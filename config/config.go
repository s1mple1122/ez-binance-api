package config

import "github.com/BurntSushi/toml"

type Config struct {
	Private Private
}
type Private struct {
	ApiKey    string
	SecretKey string
	BaseURL   string
}

func (config *Config) ReadConfigToml() error {
	if _, err := toml.DecodeFile("config.toml", config); err != nil {
		return err
	}
	return nil
}
