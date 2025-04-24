package build_savepoint

import (
	"context"
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
	repo, err := uc.Conf.GetRepo()
	if err != nil {
		return nil, coreErr.ErrRepoMissing
	}
	
	if strings.TrimSpace(req.FilePath) == "" {
		req.FilePath = "Dockerfile"
	}

	savepoints, err := uc.Parser.Parse(req.FilePath)
	if err != nil {
		return nil, err
	}

	// Handle fallback case: no savepoint provided, and none found in file
	if req.Savepoint == "" && len(savepoints) == 0 {
		lines, err := os.ReadFile(req.FilePath)
		if err != nil {
			return nil, err
		}
		content := strings.Split(string(lines), "\n")
		tag := repo + ":latest"

		if req.DryRun {
			return &BuildSavepointResult{
				Tag:           tag,
				DockerfileOut: strings.Join(content, "\n"),
				Skipped:       true,
			}, nil
		}

		if !req.Force {
			exists, err := uc.Checker.TagExists(tag)
			if err != nil {
				return nil, err
			}
			if exists {
				return &BuildSavepointResult{Tag: tag, Skipped: true}, nil
			}
		}

		path, err := uc.Writer.Write(content, "full")
		if err != nil {
			return nil, err
		}

		err = uc.Builder.Build(ctx, path, ".", tag)
		if err != nil {
			return nil, err
		}
		if req.Push {
			err = uc.Pusher.Push(tag)
			if err != nil {
				return nil, err
			}
		}

		return &BuildSavepointResult{Tag: tag, Skipped: false}, nil
	}

	// Find the requested savepoint
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

	// Read Dockerfile
	data, err := os.ReadFile(req.FilePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")


	sliced, err := uc.Slicer.Slice(lines, *target)
	if err != nil {
		return nil, err
	}

	tag := repo + ":" + req.Savepoint

	// Dry-run mode
	if req.DryRun {
		return &BuildSavepointResult{
			Tag:           tag,
			DockerfileOut: strings.Join(sliced, "\n"),
			Skipped:       true,
		}, nil
	}

	// Skip if tag exists
	if !req.Force {
		exists, err := uc.Checker.TagExists(tag)
		if err != nil {
			return nil, err
		}
		if exists {
			return &BuildSavepointResult{
				Tag:     tag,
				Skipped: true,
			}, nil
		}
	}

	// Write temp Dockerfile
	dockerfilePath, err := uc.Writer.Write(sliced, req.Savepoint)
	if err != nil {
		return nil, err
	}
	if req.Cleanup {
		defer os.Remove(dockerfilePath)
	}
	// Build
	err = uc.Builder.Build(ctx, dockerfilePath, ".", tag)
	if err != nil {
		return nil, err
	}

	// Push if requested
	if req.Push {
		err = uc.Pusher.Push(tag)
		if err != nil {
			return nil, err
		}
	}

	return &BuildSavepointResult{
		Tag:     tag,
		Skipped: false,
	}, nil
}
