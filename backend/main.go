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
	"context"
	"errors"
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/http"
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) < 3 {
		panic(errors.New("missing parameters"))
	}
	username := os.Args[1]
	password := os.Args[2]
	db := initDb(username, password)

	configColl := db.Collection("config")

	r := gin.Default()
	_ = http.NewConfigHandler(r, service.NewConfigService(repository.NewConfigRepository(dao.NewConfigDao(configColl))))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initDb(username, password string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username:   username,
		Password:   password,
		AuthSource: "fnote",
	}).SetDirect(true))
	if err != nil {
		panic(err)
	}
	return client.Database("fnote")
}
