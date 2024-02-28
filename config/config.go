package config

import "github.com/spf13/viper"

var DefaultConfig = &Config{}

type Config struct {
	Openalysis Openalysis `yaml:"openalysis"`
	DataSource DataSource `yaml:"datasource"`
	Backend    Backend    `yaml:"backend"`
}

type Openalysis struct {
}

type DataSource struct {
	MySQL MySQL `yaml:"mysql"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Backend struct {
	GraphQL GraphQL `yaml:"graphql"`
}

type GraphQL struct {
	Token string `yaml:"token"`
}

func (cfg *Config) ReadInConfig(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}
	return nil
}
