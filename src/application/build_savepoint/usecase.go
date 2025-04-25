package build_savepoint

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zeflq/dockpoint/src/core/build"
	"github.com/zeflq/dockpoint/src/core/config"
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
	Conf    config.ConfigReader
}

func NewBuildSavepointUseCase(
	parser parse.DockerfileParser,
	slicer build.DockerfileSlicer,
	writer build.TempDockerfileWriter,
	builder docker.DockerBuilder,
	checker registry.ImageChecker,
	pusher registry.ImagePusher,
	conf config.ConfigReader,
) *BuildSavepointUseCase {
	return &BuildSavepointUseCase{
		Parser:  parser,
		Slicer:  slicer,
		Writer:  writer,
		Builder: builder,
		Checker: checker,
		Pusher:  pusher,
		Conf:    conf,
	}
}

func (uc *BuildSavepointUseCase) Execute(ctx context.Context, req BuildSavepointRequest) (*BuildSavepointResult, error) {
	if req.FilePath == "" {
		req.FilePath = "Dockerfile"
	}

	repo, err := uc.Conf.GetRepo()
	if err != nil {
		return nil, coreErr.ErrRepoMissing
	}

	savepoints, err := uc.Parser.Parse(req.FilePath)
	if err != nil {
		return nil, err
	}

	// Fallback if no savepoints
	if len(savepoints) == 0 {
		savepointName := "latest"
		if req.TagOverride != "" {
			savepointName = req.TagOverride
		}
		fullSavepoint := domain.Savepoint{
			Name:      savepointName,
			StartLine: 0,
			EndLine:   9999,
		}
		return uc.buildOne(ctx, fullSavepoint, req, repo)
	}

	// Safe now: savepoints exist
	if req.TagOverride != "" && req.Savepoint == "" {
		savepoints[len(savepoints)-1].Name = req.TagOverride
	}

	if req.Savepoint == "" {
		return uc.handleBuildAll(ctx, savepoints, req, repo)
	}

	// Build specific savepoint
	var target *domain.Savepoint
	for _, sp := range savepoints {
		if sp.Name == req.Savepoint {
			target = &sp
			break
		}
	}
	if target == nil {
		return nil, coreErr.ErrSavepointNotFound
	}

	return uc.buildOne(ctx, *target, req, repo)
}

func (uc *BuildSavepointUseCase) buildOne(ctx context.Context, sp domain.Savepoint, req BuildSavepointRequest, repo string) (*BuildSavepointResult, error) {
	lines, err := os.ReadFile(req.FilePath)
	if err != nil {
		return nil, err
	}
	parsedLines := strings.Split(string(lines), "\n")

	sliced, err := uc.Slicer.Slice(parsedLines, sp)
	if err != nil {
		return nil, err
	}

	tag := repo + ":" + sp.Name

	// Inject FROM <base> if defined (i.e., not the first savepoint)
	finalLines := sliced
	if req.BaseImage != "" {
		finalLines = append([]string{fmt.Sprintf("FROM %s", req.BaseImage)}, sliced...)
	}

	if req.DryRun {
		return &BuildSavepointResult{
			Tag:           tag,
			DockerfileOut: strings.Join(finalLines, "\n"),
			Skipped:       true,
		}, nil
	}

	if !req.Force {
		exists, err := uc.Checker.TagExists(tag)
		if err != nil {
			fmt.Printf("❌ Failed to check tag existence: %v\n", err)
			return nil, err
		}
		if exists {
			fmt.Printf("⏭️  Skipped (exists): %s\n", tag)
			return &BuildSavepointResult{Tag: tag, Skipped: true}, nil
		}
	}

	dockerfilePath, err := uc.Writer.Write(finalLines, sp.Name)
	if err != nil {
		fmt.Printf("❌ Failed to write temporary Dockerfile: %v\n", err)
		return nil, err
	}
	if req.Cleanup {
		defer os.Remove(dockerfilePath)
	}

	fmt.Printf("🔨 Building: %s\n", tag)
	err = uc.Builder.Build(ctx, dockerfilePath, ".", tag)
	if err != nil {
		fmt.Printf("❌ Build failed for %s: %v\n", tag, err)
		return nil, err
	}

	if req.Push {
		fmt.Printf("📤 Pushing: %s\n", tag)
		err = uc.Pusher.Push(tag)
		if err != nil {
			fmt.Printf("❌ Push failed for %s: %v\n", tag, err)
			return nil, err
		}
	}

	return &BuildSavepointResult{Tag: tag, Skipped: false}, nil
}

func (uc *BuildSavepointUseCase) handleBuildAll(
	ctx context.Context,
	savepoints []domain.Savepoint,
	req BuildSavepointRequest,
	repo string,
) (*BuildSavepointResult, error) {
	var previews []string
	var lastResult *BuildSavepointResult

	for i, sp := range savepoints {
		currentReq := req
		if i > 0 {
			currentReq.BaseImage = fmt.Sprintf("%s:%s", repo, savepoints[i-1].Name)
		}
		res, err := uc.buildOne(ctx, sp, currentReq, repo)
		if err != nil {
			fmt.Printf("⚠️  Error building savepoint %s: %v\n", sp.Name, err)
			continue
		}
		lastResult = res
		if req.DryRun {
			previews = append(previews, fmt.Sprintf("▶ %s:\n%s\n", res.Tag, res.DockerfileOut))
		}
	}

	if req.DryRun && len(previews) > 0 {
		fmt.Println("📝 Dry-run preview of all savepoints:")
		fmt.Println("====================================")
		for _, p := range previews {
			fmt.Println(p)
			fmt.Println("------------------------------------")
		}
		fmt.Println("====================================")
	}

	return lastResult, nil
}
