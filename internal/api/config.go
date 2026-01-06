package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-api" envconfig:"SERVICE_NAME"`
	InstanceID  string `default:"" envconfig:"INSTANCE_ID"`
	AppHostName string `default:"localhost:8080" envconfig:"APP_HOSTNAME"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	if cfg.InstanceID == "" {
		// The most commonly used UUID library for Go is github.com/google/uuid
		cfg.InstanceID = uuid.New().String()
	}
	return cfg, nil
}
