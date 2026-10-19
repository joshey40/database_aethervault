package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Develop bool `koanf:"develop"`
	DB      DB   `koanf:"db"`
	API     API  `koanf:"api"`
}

type DB struct {
	Host string `koanf:"host"`
	User string `koanf:"dbuser"`
	Name string `koanf:"dbname"`
	Port int    `koanf:"port"`
}

type API struct {
	Port int `koanf:"port"`
}

func LoadConf(path string) (*Config, error) {
	k := koanf.New(".")

	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		// zap.L().Error("Loading config file failed", zap.Error(err))
		return nil, err
	}
	var config Config

	err = k.Unmarshal("", &config)
	if err != nil {
		// zap.L().Error("Unmarshalling failed", zap.Error(err))
		return nil, err
	}
	return &config, nil
}
