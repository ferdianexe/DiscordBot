package config

import (
	"context"
	"log"
)

// Service used to lists resources that used by configuration service
type Service struct {
	appConfigFile *Config // appConfigFile is a configuration for app mode binary
}

// fileReader is the interface that wraps the basic Read method to read a configuration file
type fileReader interface {
	// Read reads the configuration file at given filePath location
	// and store the config at given dest.
	// A successful call return err == nil
	Read(ctx context.Context, dest interface{}, filePath string) error
}

// NewService instantiate file reader, resource, and
// returns Service object by given config root path
func NewService(fileReader fileReader) *Service {
	var svc Service
	if err := fileReader.Read(context.Background(), &svc.appConfigFile, "config.yaml"); err != nil {
		log.Panic(context.Background(), nil, err, "fileReader.Read() app config file got error - NewResource")
	}
	return &svc
}

// GetConfig returns the configuration used in the infra as configuration.Config.
// It also satisfy infrastructure.InfraConfigProvider and app.AccessProvider interface
func (s *Service) GetConfig() *Config {
	return s.appConfigFile
}
