package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zeflq/dockpoint/src/infrastructure/docker"
	"github.com/zeflq/dockpoint/src/infrastructure/registry"
)

func ImagePusherTest() {
	// First build a test image
	imageName := "zefla/test:latest" // Replace with your Docker Hub username
	
	fmt.Println("🔨 Building test image...")
	builder := docker.NewDockerBuilder()
	ctx := context.Background()
	
	err := builder.Build(ctx, "./Dockerfile", ".", imageName)
	if err != nil {
		log.Fatalf("Failed to build image: %v", err)
	}
	
	fmt.Println("✅ Image built successfully")
	fmt.Println("📤 Pushing image to registry...")
	
	// Now push the image
	pusher := registry.NewImagePusher()
	err = pusher.Push(imageName)
	if err != nil {
		log.Fatalf("Failed to push image: %v", err)
	}
	
	fmt.Println("✅ Image pushed successfully")
}
