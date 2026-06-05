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
	"context"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func NewMongoDB() *mongox.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(viper.GetString("mongodb.uri")).SetAuth(options.Credential{
		Username:   viper.GetString("mongodb.username"),
		Password:   viper.GetString("mongodb.password"),
		AuthSource: viper.GetString("mongodb.auth_source"),
	}).SetDirect(true))
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	client := mongox.NewClient(mongoClient, &mongox.Config{})

	return client.NewDatabase(viper.GetString("mongodb.database"))
}
