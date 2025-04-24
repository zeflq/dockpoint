package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeflq/dockpoint/src/application/build_savepoint"
	"github.com/zeflq/dockpoint/src/infrastructure/build"
	"github.com/zeflq/dockpoint/src/infrastructure/config"
	"github.com/zeflq/dockpoint/src/infrastructure/docker"
	"github.com/zeflq/dockpoint/src/infrastructure/parse"
	"github.com/zeflq/dockpoint/src/infrastructure/registry"
)

var buildSavepointCmd = &cobra.Command{
	Use:   "build-savepoint [savepoint]",
	Short: "Build a Dockerfile up to a given savepoint",
	Long: `Build a Dockerfile incrementally up to a named savepoint comment.
Example:
  dockpoint build-savepoint deps --push`,
	Run: func(cmd *cobra.Command, args []string) {
		savepoint := ""
		if len(args) > 0 {
			savepoint = args[0]
		}
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			filePath = "Dockerfile"
		}
		force, _ := cmd.Flags().GetBool("force")
		push, _ := cmd.Flags().GetBool("push")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		cleanup, _ := cmd.Flags().GetBool("cleanup")

		// Build DTO
		req := build_savepoint.BuildSavepointRequest{
			Savepoint: savepoint,
			Force:     force,
			Push:      push,
			DryRun:    dryRun,
			FilePath:  filePath,
			Cleanup:   cleanup,
		}
	
		// Inject all dependencies (manual wiring for now)
		usecase := build_savepoint.NewBuildSavepointUseCase(
			parse.NewDockerfileParser(),
			build.NewDockerfileSlicer(),
			build.NewTempDockerfileWriter(),
			docker.NewDockerBuilder(),
			registry.NewImageChecker(),
			registry.NewImagePusher(),
			config.NewConfigReader(),
		)
	
		// Call the use case
		ctx := context.Background()
		result, err := usecase.Execute(ctx, req)
		if err != nil {
			fmt.Println("❌ Error:", err)
			os.Exit(1)
		}
	
		// Output result
		if result.Skipped {
			fmt.Println("⏭️ Skipped: image already exists")
		} else {
			fmt.Println("✅ Build complete:", result.Tag)
		}
	
		if dryRun {
			fmt.Println("📝 Generated Dockerfile preview:")
			fmt.Println("================================")
			fmt.Println(result.DockerfileOut)
			fmt.Println("================================")
		}		
	},
}

func init() {
	buildSavepointCmd.Flags().Bool("force", false, "Force rebuild even if image exists")
	buildSavepointCmd.Flags().Bool("push", false, "Push image after build")
	buildSavepointCmd.Flags().Bool("dry-run", false, "Skip build, just output the generated Dockerfile")
	buildSavepointCmd.Flags().StringP("file", "f", "", "Path to the Dockerfile (default: Dockerfile)")
	buildSavepointCmd.Flags().Bool("cleanup", false, "Remove temporary Dockerfile after building")

	// ✅ Register the command to rootCmd
	rootCmd.AddCommand(buildSavepointCmd)
}
