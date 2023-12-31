// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ioc

import (
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	System  SystemConfig  `json:"system" yaml:"system"`
	Gin     GinConfig     `json:"gin" yaml:"gin"`
	Logger  LoggerConfig  `json:"logger" yaml:"logger"`
	MongoDb MongodbConfig `json:"mongo_db" yaml:"mongo_db"`
}

type SystemConfig struct {
	StaticPath string `json:"static_path" yaml:"static_path"`
}

type GinConfig struct {
	AllowedOrigins []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers" yaml:"allowed_headers"`
}

type LoggerConfig struct {
	Filename   string `json:"file_name" yaml:"file_name"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	LocalTime  bool   `json:"local_time" yaml:"local_time"`
	Compress   bool   `json:"compress" yaml:"compress"`
	Level      string `json:"level" yaml:"level"`
	TimeFormat string `json:"time_format" yaml:"time_format"`
}

type MongodbConfig struct {
	Uri        string `json:"uri" yaml:"uri"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	AuthSource string `json:"auth_source" yaml:"auth_source"`
	Database   string `json:"database" yaml:"database"`
}

func InitConfig(cfgPath string) *Config {
	cfg := defaultConfig()
	if file, err := os.ReadFile(cfgPath); err != nil {
		// ignore error
		slog.Warn("InitConfig", err)
	} else if err = yaml.Unmarshal(file, cfg); err != nil {
		// ignore error
		slog.Warn("InitConfig", err)
	}
	return cfg
}

func defaultConfig() *Config {
	return &Config{
		System: SystemConfig{StaticPath: "/fnote/static/"},
		Logger: LoggerConfig{
			MaxSize:    100,
			Compress:   true,
			Level:      "INFO",
			TimeFormat: time.DateTime,
		},
	}
}
