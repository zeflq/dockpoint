package build

// TempDockerfileWriter writes a Dockerfile subset to a temporary file named after the savepoint.
type TempDockerfileWriter interface {
	Write(lines []string, savepoint string) (string, error)
}
