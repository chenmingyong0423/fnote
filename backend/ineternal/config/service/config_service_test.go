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

package service

import (
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfigService(t *testing.T) {

	testCases := []struct {
		name string
		repo repository.IConfigRepository
		want *ConfigService
	}{
		{
			name: "test",
			repo: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			service := NewConfigService(tt.repo)
			assert.NotNil(t, service)
		})
	}
}
