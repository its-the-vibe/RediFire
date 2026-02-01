package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Redis     RedisConfig     `yaml:"redis"`
	Firestore FirestoreConfig `yaml:"firestore"`
	Mappings  []Mapping       `yaml:"mappings"`
}

// RedisConfig contains Redis connection settings
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// FirestoreConfig contains Firestore connection settings
type FirestoreConfig struct {
	ProjectID       string `yaml:"projectID"`
	CredentialsFile string `yaml:"credentialsFile"`
}

// Mapping represents a source-to-target mapping
type Mapping struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override Redis password from environment if set
	if envPassword := os.Getenv("REDIS_PASSWORD"); envPassword != "" {
		cfg.Redis.Password = envPassword
	}

	// Override Firestore credentials from environment if set
	if envCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); envCreds != "" && cfg.Firestore.CredentialsFile == "" {
		cfg.Firestore.CredentialsFile = envCreds
	}

	return &cfg, nil
}
