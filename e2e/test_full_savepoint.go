package main

import (
	"context"
	"fmt"

	"github.com/zeflq/dockpoint/src/application/build_savepoint"
	"github.com/zeflq/dockpoint/src/infrastructure/build"
	"github.com/zeflq/dockpoint/src/infrastructure/docker"
	"github.com/zeflq/dockpoint/src/infrastructure/parse"
	"github.com/zeflq/dockpoint/src/infrastructure/registry"
	"github.com/zeflq/dockpoint/src/infrastructure/validator"
)
func TestFullBuildSavepointFlow() {
	fmt.Println("🔁 Testing build-savepoint end-to-end")

	ctx := context.Background()

	usecase := build_savepoint.NewBuildSavepointUseCase(
		parse.NewDockerfileParser(),
		build.NewDockerfileSlicer(),
		build.NewTempDockerfileWriter(),
		docker.NewDockerBuilder(),
		registry.NewImageChecker(),
		registry.NewImagePusher(),
		validator.NewSavepointValidator(),
		build.NewDockerfileHasher(),
		build.NewTagBuilder(),
	)

	req := build_savepoint.BuildSavepointRequest{
		Savepoint:  "base",
		Force:      true,
		Push:       true,
		DryRun:     true,
		FilePath:   "Dockerfile",
		FullTarget: "docker.io/user/app:latest",
	}

	results, err := usecase.Execute(ctx, req)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	result := results.Results[0]
	if result.Skipped {
		fmt.Println("⏭️ Image already exists:", result.Tag)
	} else {
		fmt.Println("✅ Build completed:", result.Tag)
	}

	if req.DryRun {
		fmt.Println("📝 Generated Dockerfile preview:")
		fmt.Println("================================")
		fmt.Println(result.DockerfileOut)
		fmt.Println("================================")
	}
}
