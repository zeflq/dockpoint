package main

import (
	"fmt"

	"github.com/zeflq/dockpoint/src/infrastructure/build"
)

func TestDockerfileWriter() {
	writer := build.NewTempDockerfileWriter()

	dockerLines := []string{
		"FROM alpine",
		"RUN echo 'Hello from savepoint!'",
	}

	filePath, err := writer.Write(dockerLines, "test")
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Dockerfile written to:", filePath)
}
