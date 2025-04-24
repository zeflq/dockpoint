package build_savepoint

// BuildSavepointRequest represents CLI input and options passed into the use case.
type BuildSavepointRequest struct {
	Savepoint string
	Force     bool
	Push      bool
	DryRun    bool
	FilePath  string // path to Dockerfile, default is "Dockerfile"
}

// BuildSavepointResult represents the result of the build process.
type BuildSavepointResult struct {
	Tag           string
	DockerfileOut string // Used in dry-run mode
	Skipped       bool
}
