package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zeflq/dockpoint/src/infrastructure/docker"
)

func TestDockerBuild() {
	builder := docker.NewDockerBuilder()
	ctx := context.Background()

	// Check if Dockerfile exists
	dockerfilePath := filepath.Join(".", "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		log.Fatalf("Dockerfile not found at %s", dockerfilePath)
	}

	fmt.Println("🔨 Starting Docker build...")
	err := builder.Build(ctx, dockerfilePath, ".", "dockpoint:test")
	if err != nil {
		log.Printf("❌ Build failed: %v", err)
		return
	}

	fmt.Println("✅ Build completed")
}
