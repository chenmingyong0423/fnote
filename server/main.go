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

package main

import (
	"bytes"
	"flag"
	"os"

	"github.com/spf13/viper"
)

var (
	configPath = flag.String("website_config", "./config/fnote.yaml", "the path of website_config")
	port       = flag.String("port", ":8080", "HTTP port")
)

func main() {
	flag.Parse()
	err := initViper(*configPath)
	if err != nil {
		panic(err)
	}
	app, err := initializeApp()
	if err != nil {
		panic(err)
	}
	err = app.Run(*port)
	if err != nil {
		panic(err)
	}
}

func initViper(cfgPath string) error {
	viper.SetConfigType("yaml")
	readFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	// 寻找配置文件并读取
	err = viper.ReadConfig(bytes.NewReader(readFile))
	if err != nil {
		return err
	}
	return nil
}
