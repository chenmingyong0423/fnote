// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
)

var Config = &config{}

type config struct {
	IsWebsiteInitialized bool
	Domain               string
}

func IsWebsiteInitialized() bool {
	return Config.IsWebsiteInitialized
}

func IsWebsiteInitializedFn(db *mongox.Database) (func() bool, error) {
	var ok bool
	ctx := context.Background()
	result := bson.M{}
	err := db.Database().Collection("configs").FindOne(ctx, query.Eq("typ", "website")).Decode(&result)
	if err != nil {
		return nil, err
	}
	props := result["props"].(bson.D)
	if props == nil {
		return nil, errors.New("The collections of config does not have a field named initialized.")
	}

	hasWebsiteInitField := false
	for _, prop := range props {
		if prop.Key == "website_init" {
			Config.IsWebsiteInitialized, ok = prop.Value.(bool)
			if !ok {
				return nil, errors.New("The collections of config does not have a field named initialized.")
			}
			hasWebsiteInitField = true
			break
		}
	}

	if !hasWebsiteInitField {
		return nil, errors.New("The collections of config does not have a field named initialized.")
	}

	return func() bool {
		return Config.IsWebsiteInitialized
	}, nil
}
