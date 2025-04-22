package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
		force, _ := cmd.Flags().GetBool("force")
		push, _ := cmd.Flags().GetBool("push")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Placeholder for actual use case call
		fmt.Printf("🚧 BuildSavepoint called with:\n- Savepoint: %s\n- Force: %v\n- Push: %v\n- DryRun: %v\n",
			savepoint, force, push, dryRun)
	},
}

func init() {
	buildSavepointCmd.Flags().Bool("force", false, "Force rebuild even if image exists")
	buildSavepointCmd.Flags().Bool("push", false, "Push image after build")
	buildSavepointCmd.Flags().Bool("dry-run", false, "Skip build, just output the generated Dockerfile")

	// ✅ Register the command to rootCmd
	rootCmd.AddCommand(buildSavepointCmd)
}
