package config

import (
	"encoding/json"
	"os"

	core "github.com/zeflq/dockpoint/src/core/config"
	coreErr "github.com/zeflq/dockpoint/src/core/errors"
)

type ConfigReaderImpl struct{}

func NewConfigReader() core.ConfigReader {
	return &ConfigReaderImpl{}
}

type dotfileConfig struct {
	Repo string `json:"repo"`
}

func (r *ConfigReaderImpl) GetRepo() (string, error) {
	data, err := os.ReadFile(".dockpointrc.json")
	if err != nil {
		return "", coreErr.ErrRepoMissing
	}

	var cfg dotfileConfig
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.Repo == "" {
		return "", coreErr.ErrRepoMissing
	}

	return cfg.Repo, nil
}
