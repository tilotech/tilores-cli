package pkg

import (
	"encoding/json"
	"os"
)

// Config represents the tilores.json config that holds the version number.
type Config struct {
	Version string `json:"version"`
}

// LoadConfig reads the current tilores.json.
func LoadConfig() (*Config, error) {
	f, err := os.Open("tilores.json")
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	decoder := json.NewDecoder(f)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// SaveConfig updates or creates the tilores.json.
func SaveConfig(config *Config) error {
	j, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile("tilores.json", j, 0600)
}

// DefaultConfig returns the default config if tilores.json does not exist.
//
// This should only be the case for very old projects or for cases where a user
// removed the tilores.json.
func DefaultConfig() *Config {
	return &Config{
		Version: "v0.0.0",
	}
}
