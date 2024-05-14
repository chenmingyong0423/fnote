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

package domain

import (
	"errors"
	"slices"
)

var countStatsTypeSlice = []CountStatsType{CountStatsTypePostCountInCategory, CountStatsTypePostCountInTag}

type CountStatsType string

const (
	// CountStatsTypePostCountInCategory 分类下的文章数量
	CountStatsTypePostCountInCategory CountStatsType = "PostCountInCategory"
	// CountStatsTypePostCountInTag 标签下的文章数量
	CountStatsTypePostCountInTag CountStatsType = "PostCountInTag"
	// CountStatsTypePostCount 网站的文章数量
	CountStatsTypePostCount CountStatsType = "PostCount"
	// CountStatsTypeCategoryCount 分类数量
	CountStatsTypeCategoryCount CountStatsType = "CategoryCount"
	// CountStatsTypeTagCount 标签数量
	CountStatsTypeTagCount CountStatsType = "TagCount"
	// CountStatsTypeCommentCount 评论数量
	CountStatsTypeCommentCount CountStatsType = "CommentCount"
	// CountStatsTypeLikeCount 点赞数量
	CountStatsTypeLikeCount CountStatsType = "LikeCount"
	// CountStatsTypeWebsiteViewCount 网站总访问量
	CountStatsTypeWebsiteViewCount CountStatsType = "WebsiteViewCount"
)

func (s CountStatsType) ToString() string {
	return string(s)
}

func (s CountStatsType) Valid() error {
	if !slices.Contains(countStatsTypeSlice, s) {
		return errors.New("invalid count stats type")
	}
	return nil
}

type CountStats struct {
	Id          string
	Type        CountStatsType
	ReferenceId string
	Count       int64
}

type WebsiteCountStats struct {
	PostCount        int64
	CategoryCount    int64
	TagCount         int64
	CommentCount     int64
	LikeCount        int64
	WebsiteViewCount int64
}

func (wcs *WebsiteCountStats) SetCountByType(typ CountStatsType, count int64) {
	switch typ {
	case CountStatsTypePostCount:
		wcs.PostCount = count
	case CountStatsTypeCategoryCount:
		wcs.CategoryCount = count
	case CountStatsTypeTagCount:
		wcs.TagCount = count
	case CountStatsTypeCommentCount:
		wcs.CommentCount = count
	case CountStatsTypeLikeCount:
		wcs.LikeCount = count
	case CountStatsTypeWebsiteViewCount:
		wcs.WebsiteViewCount = count
	}
}
