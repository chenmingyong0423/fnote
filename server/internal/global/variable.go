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

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chenmingyong0423/go-mongox/builder/query"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	isWebsiteInitialized bool
)

func IsWebsiteInitialized() bool {
	return isWebsiteInitialized
}

func WebsiteInitialized() {
	isWebsiteInitialized = true
}

func IsWebsiteInitializedFn(db *mongo.Database) (func() bool, error) {
	var ok bool
	ctx := context.Background()
	result := bson.M{}
	err := db.Collection("configs").FindOne(ctx, query.Eq("typ", "website")).Decode(&result)
	if err != nil {
		return nil, err
	}
	props := result["props"].(bson.M)
	if props == nil {
		return nil, errors.New("The collections of config does not have a field named initialized.")
	}

	if isWebsiteInitialized, ok = props["website_init"].(bool); !ok {
		return nil, errors.New("The collections of config does not have a field named initialized.")
	} else {
		return func() bool {
			return isWebsiteInitialized
		}, nil
	}
}
