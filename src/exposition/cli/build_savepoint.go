package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeflq/dockpoint/src/application/build_savepoint"
	"github.com/zeflq/dockpoint/src/infrastructure/build"
	"github.com/zeflq/dockpoint/src/infrastructure/docker"
	"github.com/zeflq/dockpoint/src/infrastructure/parse"
	"github.com/zeflq/dockpoint/src/infrastructure/registry"
)

var buildSavepointCmd = &cobra.Command{
	Use:   "build-savepoint",
	Short: "Incrementally build a Dockerfile up to savepoints",
	Long: `dockpoint build-savepoint builds a Dockerfile step-by-step based on savepoints defined inside the file.

Each savepoint is a special comment like:
    # savepoint: mypoint

By default, all savepoints are built sequentially.
You must specify a full target image using -t in the format 'repo/image:tag'.

Examples:
  # Build and push all savepoints to Docker Hub
  dockpoint build-savepoint -t docker.io/myuser/myapp:1.1.2 --push

  # Build only one specific savepoint
  dockpoint build-savepoint base -t docker.io/myuser/myapp:1.1.2

Flags:
  -t, --target       Full image reference to build (e.g., docker.io/user/app:1.1.2) [REQUIRED]
  --savepoint        Build only up to a specific savepoint name (optional)
  -f, --file         Path to the Dockerfile (default: Dockerfile)
  --force            Force rebuild even if the image already exists
  --push             Push images to the registry after building
  --dry-run          Skip actual build, output generated Dockerfile previews
`,
	Run: func(cmd *cobra.Command, args []string) {
		savepoint := ""
		if len(args) > 0 {
			savepoint = args[0]
		}

		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			filePath = "Dockerfile"
		}
		target, _ := cmd.Flags().GetString("target")
		force, _ := cmd.Flags().GetBool("force")
		push, _ := cmd.Flags().GetBool("push")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		req := build_savepoint.BuildSavepointRequest{
			Savepoint:  savepoint,
			FilePath:   filePath,
			FullTarget: target,
			Force:      force,
			Push:       push,
			DryRun:     dryRun,
		}

		usecase := build_savepoint.NewBuildSavepointUseCase(
			parse.NewDockerfileParser(),
			build.NewDockerfileSlicer(),
			build.NewTempDockerfileWriter(),
			docker.NewDockerBuilder(),
			registry.NewImageChecker(),
			registry.NewImagePusher(),
		)

		ctx := context.Background()
		result, err := usecase.Execute(ctx, req)
		if err != nil {
			fmt.Println("❌ Error:", err)
			os.Exit(1)
		}
		_ = result
	},
}

func init() {
	buildSavepointCmd.Flags().Bool("force", false, "Force rebuild even if image exists")
	buildSavepointCmd.Flags().Bool("push", false, "Push image after build")
	buildSavepointCmd.Flags().Bool("dry-run", false, "Skip build, just output the generated Dockerfile")
	buildSavepointCmd.Flags().StringP("file", "f", "", "Path to the Dockerfile (default: Dockerfile)")
	buildSavepointCmd.Flags().StringP("target", "t", "", "Full image name (e.g., docker.io/user/app:1.1.2)")

	// ✅ Register the command to rootCmd
	rootCmd.AddCommand(buildSavepointCmd)
}