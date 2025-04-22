package config

// ConfigReader loads .dockpointrc.json and returns the base image repo.
type ConfigReader interface {
	GetRepo() (string, error)
}
