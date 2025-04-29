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
	Hasher domain.DockerfileHasher
	TagBuilder domain.TagBuilder
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
	hasher domain.DockerfileHasher,
	tagBuilder domain.TagBuilder,
) *BuildSavepointUseCase {
	return &BuildSavepointUseCase{
		Parser:  parser,
		Slicer:  slicer,
		Writer:  writer,
		Builder: builder,
		Checker: checker,
		Pusher:  pusher,
		Validator: validator,
		Hasher: hasher,
		TagBuilder: tagBuilder,
	}
}

func ParseAndValidateFullTarget(ref string) (string, string, error) {
	if ref == "" {
		return "", "", fmt.Errorf("❌ Full target is empty")
	}

	if !strings.Contains(ref, ":") {
		return "", "", fmt.Errorf("❌ Invalid image reference '%s'. You must specify repo/image:tag format (missing ':tag')", ref)
	}

	parts := strings.SplitN(ref, ":", 2)
	repo := parts[0]
	tag := parts[1]

	if repo == "" || tag == "" {
		return "", "", fmt.Errorf("❌ Invalid image reference '%s'. Must be repo/image:tag format", ref)
	}

	if !strings.Contains(repo, "/") {
		return "", "", fmt.Errorf("❌ Invalid image reference '%s'. The repository must contain at least one '/' like 'docker.io/user/app:tag'", ref)
	}

	return repo, tag, nil
}

func FindSavepoint(savepoints []domain.Savepoint, name string) (*domain.Savepoint, int, error) {
	for i, sp := range savepoints {
		if sp.Name == name {
			return &sp, i, nil
		}
	}
	return nil, -1, coreErr.ErrSavepointNotFound
}
func (uc *BuildSavepointUseCase) readDockerfile(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to read Dockerfile: %v", err)
	}
	return string(fileBytes), nil
}

func (uc *BuildSavepointUseCase) prepareSavepoints(filePath string) ([]domain.Savepoint, error) {
	savepoints, err := uc.Parser.Parse(filePath)
	if err != nil {
		return nil, err
	}

	if err := uc.Validator.Validate(savepoints); err != nil {
		return nil, err
	}

	return savepoints, nil
}

func (uc *BuildSavepointUseCase) validateRequest(req *BuildSavepointRequest) error {
	if req.FilePath == "" {
		req.FilePath = "Dockerfile"
	}

	if req.FullTarget == "" {
		return fmt.Errorf("❌ Missing required -t flag. You must specify a full image name like 'docker.io/user/app:1.1.2'")
	}
	return nil
}

func (uc *BuildSavepointUseCase) Execute(ctx context.Context, req BuildSavepointRequest) (*BuildSavepointResultList, error) {
	if err := uc.validateRequest(&req); err != nil {
		return nil, err
	}

	imageRepo, baseTag, err := ParseAndValidateFullTarget(req.FullTarget)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔍 Building base image repo: %s, base tag: %s\n", imageRepo, baseTag)

	fileContent, err := uc.readDockerfile(req.FilePath)
	if err != nil {
		return nil, err
	}

	savepoints, err := uc.prepareSavepoints(req.FilePath)
	if err != nil {
		return nil, err
	}

	return uc.processSavepoints(ctx, savepoints, req, fileContent, imageRepo, baseTag)
}


func (uc *BuildSavepointUseCase) processSavepoints(
	ctx context.Context,
	savepoints []domain.Savepoint,
	req BuildSavepointRequest,
	fileContent string,
	imageRepo string,
	baseTag string,
) (*BuildSavepointResultList, error) {
	var targetSavepoints []domain.Savepoint

	if len(savepoints) == 0 {
		// Handle no savepoints case
		lines := strings.Split(fileContent, "\n")
		realEndLine := len(lines) - 1
		targetSavepoints = []domain.Savepoint{{
			Name:      baseTag,
			StartLine: 0,
			EndLine:   realEndLine,
		}}
	} else if req.Savepoint != "" {
		// Handle single savepoint case
		target, _, err := FindSavepoint(savepoints, req.Savepoint)
		if err != nil {
			return nil, err
		}
		targetSavepoints = []domain.Savepoint{*target}
	} else {
		// Handle all savepoints case
		targetSavepoints = savepoints
	}

	results, err := uc.handleBuildAll(ctx, targetSavepoints, req, fileContent, imageRepo, baseTag)
	if err != nil {
		return nil, err
	}

	return &BuildSavepointResultList{Results: results}, nil
}

func PrepareDockerfileContent(fileContent string, sp domain.Savepoint, baseImage string) (string, error) {
	lines := strings.Split(fileContent, "\n")
	if sp.StartLine < 0 || sp.EndLine >= len(lines) || sp.StartLine > sp.EndLine {
		return "", fmt.Errorf("invalid savepoint line range: start=%d end=%d totalLines=%d", sp.StartLine, sp.EndLine, len(lines))
	}

	sliced := lines[sp.StartLine : sp.EndLine+1]

	finalLines := sliced
	if baseImage != "" {
		finalLines = append([]string{fmt.Sprintf("FROM %s", baseImage)}, sliced...)
	}

	dockerfileContent := strings.Join(finalLines, "\n")
	return dockerfileContent, nil
}

func (uc *BuildSavepointUseCase) buildOneWithTag(
	ctx context.Context,
	sp domain.Savepoint,
	req BuildSavepointRequest,
	dockerfileContent string,
	finalTag string,
) (*BuildSavepointResult, error) {
	if req.DryRun {
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

	dockerfilePath, err := uc.Writer.Write(strings.Split(dockerfileContent, "\n"), sp.Name)
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
	fileContent string,
	imageRepo string,
	baseTag string,
) ([]*BuildSavepointResult, error) {

	var results []*BuildSavepointResult
	var previousTag string

	for i, sp := range savepoints {
		currentReq := req
		var baseImage string

		if i > 0 {
			baseImage = previousTag
			currentReq.BaseImage = previousTag
		}

		dockerfileContent, err := PrepareDockerfileContent(fileContent, sp, baseImage)
		if err != nil {
			fmt.Printf("❌ Failed preparing Dockerfile content for savepoint %s: %v\n", sp.Name, err)
			continue
		}
		isLast := i == len(savepoints)-1
		if isLast {
			//force the last image to be rebuilt
			currentReq.Force = true
		}
		tagCtx := domain.TagContext{
			FileContent: dockerfileContent,
			Repo:        imageRepo,
			SavepointName: sp.Name,
			FinalImageTag: baseTag,
			IsLast:      isLast,
		}		

		finalTag, err := uc.TagBuilder.BuildFinalTag(tagCtx)
		if err != nil {
			fmt.Printf("❌ Failed to build tag for savepoint %s: %v\n", sp.Name, err)
			continue
		}

		fmt.Printf("🔨 Building Savepoint: %s -> %s\n", sp.Name, finalTag)

		res, err := uc.buildOneWithTag(ctx, sp, currentReq, dockerfileContent, finalTag)
		if err != nil {
			fmt.Printf("⚠️  Error building savepoint %s: %v\n", sp.Name, err)
			continue
		}
		previousTag = res.Tag
		results = append(results, res)
	}

	return results, nil
}
