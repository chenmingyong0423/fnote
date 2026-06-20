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

package repository

import (
	"reflect"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository/dao"
)

func TestContentLocations(t *testing.T) {
	locations := contentLocations("/static/example.png", &dao.UsageContent{
		Content:  `<p><img src="http://localhost:8080/static/example.png"></p>`,
		CoverImg: "/static/example.png",
	}, "文章正文", "文章封面")

	expected := []string{"文章正文", "文章封面"}
	if !reflect.DeepEqual(locations, expected) {
		t.Fatalf("expected %v, got %v", expected, locations)
	}
}

func TestConfigLocations(t *testing.T) {
	locations := configLocations("/static/example.png", dao.UsageConfigProps{
		WebsiteIcon:    "/static/example.png",
		WebsiteRecords: []string{`<img src="/static/example.png">`, `/static/example.png`},
		List: []dao.UsageConfigItem{
			{Image: "/static/example.png"},
			{CoverImg: "/static/example.png"},
		},
	})

	expected := []string{"站点 Logo", "备案信息", "收款二维码", "轮播图"}
	if !reflect.DeepEqual(locations, expected) {
		t.Fatalf("expected %v, got %v", expected, locations)
	}
}

func TestConfigName(t *testing.T) {
	if got := configName("pay"); got != "收款配置" {
		t.Fatalf("expected 收款配置, got %s", got)
	}
	if got := configName("custom"); got != "custom" {
		t.Fatalf("expected custom, got %s", got)
	}
}
