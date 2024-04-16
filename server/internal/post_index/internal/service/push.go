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

package service

import (
	"bytes"
	"context"
	"io"

	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/domain"
	httpchain "github.com/chenmingyong0423/go-http-chain"
	"github.com/spf13/viper"
)

type BaiduService struct {
	client *httpchain.Client
}

func NewBaiduService() *BaiduService {
	return &BaiduService{client: httpchain.NewDefault()}
}

func (s *BaiduService) Push(ctx context.Context, site, token, urls string) (*domain.BaiduResponse, error) {
	resp, err := s.client.Post(viper.GetString("push.baidu.endpoint")+"?site="+site).
		AddQuery("token", token).
		SetHeader(httpchain.HeaderContentType, httpchain.ContentTypeTextPlain).
		SetHeader("host", "data.zz.baidu.com").
		SetBody(urls).SetBodyEncodeFunc(func(body any) (io.Reader, error) {
		var buf bytes.Buffer
		buf.WriteString(body.(string))
		return &buf, nil
	}).Call(ctx).Result()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		var baiduErrorResponse domain.BaiduErrorResponse
		err = httpchain.DecodeRespBody(resp, &baiduErrorResponse)
		if err != nil {
			return nil, err
		}
		return &domain.BaiduResponse{BaiduErrorResponse: baiduErrorResponse}, nil
	}
	var baiduSuccessResponse domain.BaiduSuccessResponse
	err = httpchain.DecodeRespBody(resp, &baiduSuccessResponse)
	if err != nil {
		return nil, err
	}
	return &domain.BaiduResponse{BaiduSuccessResponse: baiduSuccessResponse}, nil
}
