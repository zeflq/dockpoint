package build_savepoint

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zeflq/dockpoint/src/core/build"
	"github.com/zeflq/dockpoint/src/core/docker"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
	"github.com/zeflq/dockpoint/src/core/parse"
	"github.com/zeflq/dockpoint/src/core/registry"
	"github.com/zeflq/dockpoint/src/domain"
)

type BuildSavepointUseCase struct {
	Parser  parse.DockerfileParser
	Slicer  build.DockerfileSlicer
	Writer  build.TempDockerfileWriter
	Builder docker.DockerBuilder
	Checker registry.ImageChecker
	Pusher  registry.ImagePusher
	Validator domain.SavepointValidator
}
type BuildSavepointResultList struct {
	Results []*BuildSavepointResult
}

func NewBuildSavepointUseCase(
	parser parse.DockerfileParser,
	slicer build.DockerfileSlicer,
	writer build.TempDockerfileWriter,
	builder docker.DockerBuilder,
	checker registry.ImageChecker,
	pusher registry.ImagePusher,
	validator domain.SavepointValidator,
) *BuildSavepointUseCase {
	return &BuildSavepointUseCase{
		Parser:  parser,
		Slicer:  slicer,
		Writer:  writer,
		Builder: builder,
		Checker: checker,
		Pusher:  pusher,
		Validator: validator,
	}
}

func (uc *BuildSavepointUseCase) Execute(ctx context.Context, req BuildSavepointRequest) (*BuildSavepointResultList, error) {
	if req.FilePath == "" {
			req.FilePath = "Dockerfile"
	}

	if req.FullTarget == "" {
			return nil, fmt.Errorf("❌ Missing required -t flag. You must specify a full image name like 'docker.io/user/app:1.1.2'")
	}

	ref := req.FullTarget
	if !strings.Contains(ref, ":") {
			return nil, fmt.Errorf("❌ Invalid image reference '%s'. You must specify repo/image:tag format (missing ':tag')", ref)
	}

	parts := strings.SplitN(ref, ":", 2)
	candidateRepo := parts[0]
	candidateTag := parts[1]

	if candidateRepo == "" || candidateTag == "" {
			return nil, fmt.Errorf("❌ Invalid image reference '%s'. Must be repo/image:tag format", ref)
	}

	if !strings.Contains(candidateRepo, "/") {
			return nil, fmt.Errorf("❌ Invalid image reference '%s'. The repository must contain at least one '/' like 'docker.io/user/app:tag'", ref)
	}

	imageRepo := candidateRepo
	baseTag := candidateTag

	fmt.Printf("🔍 Building base image repo: %s, base tag: %s\n", imageRepo, baseTag)

	savepoints, err := uc.Parser.Parse(req.FilePath)
	if err != nil {
			return nil, err
	}
	
	if err := uc.Validator.Validate(savepoints); err != nil {
		return nil, err
	}

	if len(savepoints) == 0 {
			// No savepoints: build full Dockerfile
			fullSavepoint := domain.Savepoint{
					Name:      baseTag,
					StartLine: 0,
					EndLine:   9999,
			}
			finalTag := fmt.Sprintf("%s:%s", imageRepo, baseTag)
			res, err := uc.buildOneWithTag(ctx, fullSavepoint, req, finalTag)
			if err != nil {
				return nil, err
			}
			return &BuildSavepointResultList{Results: []*BuildSavepointResult{res}}, nil
	}

	if req.Savepoint == "" {
			// Build all savepoints
			results, err := uc.handleBuildAll(ctx, savepoints, req, imageRepo, baseTag)
			if err != nil {
				return nil, err
			}
			return &BuildSavepointResultList{Results: results}, nil
	}

	// Build specific savepoint
	var target *domain.Savepoint
	var targetIndex int = -1
	for i, sp := range savepoints {
			if sp.Name == req.Savepoint {
					target = &sp
					targetIndex = i
					break
			}
	}
	if target == nil {
			return nil, coreErr.ErrSavepointNotFound
	}

	finalTag := fmt.Sprintf("%s:%s", imageRepo, target.Name)

	// ✅ Fix: set BaseImage manually if needed
	if targetIndex > 0 {
			req.BaseImage = fmt.Sprintf("%s:%s", imageRepo, savepoints[targetIndex-1].Name)
	}

	res, err := uc.buildOneWithTag(ctx, *target, req, finalTag)
	if err != nil {
		return nil, err
	}

	return &BuildSavepointResultList{Results: []*BuildSavepointResult{res}}, nil
	
}



func (uc *BuildSavepointUseCase) buildOneWithTag(
	ctx context.Context,
	sp domain.Savepoint,
	req BuildSavepointRequest,
	finalTag string,
) (*BuildSavepointResult, error) {
	lines, err := os.ReadFile(req.FilePath)
	if err != nil {
			return nil, err
	}
	parsedLines := strings.Split(string(lines), "\n")

	sliced, err := uc.Slicer.Slice(parsedLines, sp)
	if err != nil {
			return nil, err
	}

	finalLines := sliced
	if req.BaseImage != "" {
			finalLines = append([]string{fmt.Sprintf("FROM %s", req.BaseImage)}, sliced...)
	}

	dockerfileContent := strings.Join(finalLines, "\n")

	if req.DryRun {
			// fmt.Println("📝 Dry-run preview:")
			// fmt.Println("====================================")
			// fmt.Printf("▶ %s:\n", finalTag)
			// fmt.Println(dockerfileContent)
			// fmt.Println("------------------------------------")
			// fmt.Println("====================================")
			return &BuildSavepointResult{
					Tag:           finalTag,
					DockerfileOut: dockerfileContent,
					Skipped:       true,
			}, nil
	}

	if !req.Force {
			exists, err := uc.Checker.TagExists(finalTag)
			if err != nil {
					fmt.Printf("❌ Failed to check tag existence: %v\n", err)
					return nil, err
			}
			if exists {
					fmt.Printf("⏭️  Skipped (already exists): %s\n", finalTag)
					return &BuildSavepointResult{Tag: finalTag, Skipped: true}, nil
			}
	}

	dockerfilePath, err := uc.Writer.Write(finalLines, sp.Name)
	if err != nil {
			fmt.Printf("❌ Failed to write temporary Dockerfile: %v\n", err)
			return nil, err
	}
	defer os.Remove(dockerfilePath)

	fmt.Printf("🔨 Building: %s\n", finalTag)
	err = uc.Builder.Build(ctx, dockerfilePath, ".", finalTag)
	if err != nil {
			fmt.Printf("❌ Build failed for %s: %v\n", finalTag, err)
			return nil, err
	}

	if req.Push {
			fmt.Printf("📤 Pushing: %s\n", finalTag)
			err = uc.Pusher.Push(finalTag)
			if err != nil {
					fmt.Printf("❌ Push failed for %s: %v\n", finalTag, err)
					return nil, err
			}
	}

	return &BuildSavepointResult{Tag: finalTag, Skipped: false}, nil
}


func (uc *BuildSavepointUseCase) handleBuildAll(
	ctx context.Context,
	savepoints []domain.Savepoint,
	req BuildSavepointRequest,
	imageRepo string,
	baseTag string,
) ([]*BuildSavepointResult, error) {
	var results []*BuildSavepointResult

	for i, sp := range savepoints {
			currentReq := req
			if i > 0 {
					currentReq.BaseImage = fmt.Sprintf("%s:%s", imageRepo, savepoints[i-1].Name)
			}

			var finalTag string
			if i == len(savepoints)-1 {
					// Last savepoint → baseTag (provided or default latest)
					finalTag = fmt.Sprintf("%s:%s", imageRepo, baseTag)
			} else {
					// Inner savepoints → savepoint name
					finalTag = fmt.Sprintf("%s:%s", imageRepo, sp.Name)
			}

			fmt.Printf("🔨 Building Savepoint: %s -> %s\n", sp.Name, finalTag)

			res, err := uc.buildOneWithTag(ctx, sp, currentReq, finalTag)
			if err != nil {
					fmt.Printf("⚠️  Error building savepoint %s: %v\n", sp.Name, err)
					continue
			}
			results = append(results, res)
	}

	return results, nil
}
