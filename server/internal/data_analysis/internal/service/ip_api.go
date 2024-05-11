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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/data_analysis/internal/domain"
	"github.com/chenmingyong0423/gkit/slice"
	httpchain "github.com/chenmingyong0423/go-http-chain"
	"golang.org/x/sync/errgroup"
)

func NewIpApiService() *IpApiService {
	return &IpApiService{
		host: "http://ip-api.com",
		client: httpchain.NewWithClient(&http.Client{
			Timeout: time.Second * 10,
		}),
	}
}

type IIpApiService interface {
	BatchGetLocation(ctx context.Context, ips []string) ([]domain.IpApi, error)
}

var _ IIpApiService = (*IpApiService)(nil)

type IpApiService struct {
	host   string
	client *httpchain.Client
}

func (s *IpApiService) BatchGetLocation(ctx context.Context, ips []string) ([]domain.IpApi, error) {
	const batchSize = 100
	numberOfBatches := (len(ips) + batchSize - 1) / batchSize // 计算批次数量
	results := make(chan []domain.IpApi, numberOfBatches)     // 创建带有足够缓冲的通道
	defer func() {
		close(results)
	}()
	var eg errgroup.Group
	// 将 IPs 切分为每个批次大小为 100
	for start := 0; start < len(ips); start += batchSize {
		end := start + batchSize
		if end > len(ips) {
			end = len(ips)
		}
		// 准备当前批次的 IPs
		batchIps := ips[start:end]
		eg.Go(func() error {
			batchBody := slice.Map(batchIps, func(idx int, ip string) domain.IpApiRequestBody {
				return domain.IpApiRequestBody{
					Query:  ip,
					Fields: "city,country,countryCode,query",
					Lang:   "zh-CN",
				}
			})

			// 发起 API 调用
			batchResult := make([]domain.IpApi, 0, len(batchIps))
			err := s.client.Post(s.host + "/batch").SetBody(batchBody).SetBodyEncodeFunc(func(body any) (io.Reader, error) {
				marshal, err := json.Marshal(body)
				if err != nil {
					return nil, err
				}
				return io.NopCloser(bytes.NewReader(marshal)), nil
			}).Call(ctx).DecodeRespBody(&batchResult)
			if err != nil {
				return err
			}

			results <- batchResult
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	var finalResult = make([]domain.IpApi, 0, len(ips))
	for res := range results {
		finalResult = append(finalResult, res...)
		numberOfBatches--
		if numberOfBatches == 0 {
			break
		}
	}
	return finalResult, nil
}
