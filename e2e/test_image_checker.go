package main

import (
	"fmt"
	"log"

	"github.com/zeflq/dockpoint/src/infrastructure/registry"
)

func TestImageChecker() {
	checker := registry.NewImageChecker()
	
	// Test with a public image
	tag := "redis:latest"
	exists, err := checker.TagExists(tag)
	if err != nil {
		log.Fatalf("Error checking tag: %v", err)
	}
	
	fmt.Printf("Tag %s exists: %v\n", tag, exists)
}