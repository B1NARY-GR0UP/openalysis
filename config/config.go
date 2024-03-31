// Copyright 2024 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/spf13/viper"
)

var GlobalConfig = &Config{}

type Config struct {
	Groups     []Group    `yaml:"groups"`
	DataSource DataSource `yaml:"datasource"`
	Backend    Backend    `yaml:"backend"`
}

type Group struct {
	Name  string   `yaml:"name"`
	Orgs  []string `yaml:"orgs"`
	Repos []string `yaml:"repos"`
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
	Cron  string `yaml:"cron"`
	Token string `yaml:"token"`
	Retry int    `yaml:"retry"`
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

func Init(path string) error {
	if path == "" {
		path = "./default.yaml"
	}
	if err := GlobalConfig.ReadInConfig(path); err != nil {
		return err
	}
	return nil
}
