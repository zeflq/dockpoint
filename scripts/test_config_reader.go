package main

import (
	"fmt"

	"github.com/zeflq/dockpoint/src/infrastructure/config"
)

func TestConfigReader() {
	reader := config.NewConfigReader()
	repo, err := reader.GetRepo()
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ Repo found:", repo)
}
