package pkg

import (
	"encoding/json"
	"os"
)

type Config struct {
	Version string `json:"version"`
}

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

func SaveConfig(config *Config) error {
	j, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile("tilores.json", j, 0600)
}

func DefaultConfig() *Config {
	return &Config{
		Version: "v0.0.0",
	}
}
