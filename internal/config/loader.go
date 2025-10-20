package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrEmptyConfigPath = errors.New("config path is empty")
)

func Load() (*Config, error) {
	const op = "Config.Load"

	var (
		cfg        Config
		configPath string
		err        error
	)

	flag.StringVar(&configPath, "c", "", "path/to/config/file.yaml")
	flag.Parse()
	if configPath == "" {
		return nil, ErrEmptyConfigPath
	}

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to ReadConfig from file with %w", op, err)
	}

	return &cfg, nil
}
