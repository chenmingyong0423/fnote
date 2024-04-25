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

package apiwrap

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBody[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func (er ResponseBody[T]) Error() string {
	return fmt.Sprintf("%d:%s", er.Code, er.Message)
}

func SuccessResponse() *ResponseBody[any] {
	return &ResponseBody[any]{
		Code:    0,
		Message: "success",
	}
}

func SuccessResponseWithData[T any](data T) *ResponseBody[T] {
	return &ResponseBody[T]{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func NewResponseBody[T any](code int, message string, data T) *ResponseBody[T] {
	return &ResponseBody[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type HttpCodeError int

type ErrorResponseBody struct {
	HttpCode int
	Message  string
}

func NewErrorResponseBody(httpCode int, message string) ErrorResponseBody {
	return ErrorResponseBody{
		HttpCode: httpCode,
		Message:  message,
	}
}

func (er ErrorResponseBody) Error() string {
	return er.Message
}

type ListVO[T any] struct {
	List []T `json:"list"`
}

func NewListVO[T any](t []T) ListVO[T] {
	return ListVO[T]{
		List: t,
	}
}

type Page struct {
	// 当前页
	PageNo int64 `form:"pageNo" binding:"required"`
	// 每页数量
	PageSize int64 `form:"pageSize" binding:"required"`
}

type PageVO[T any] struct {
	Page
	// 总页数
	TotalPages int64 `json:"totalPages"`
	// 总数量
	TotalCount int64 `json:"totalCount"`
	List       []T   `json:"list"`
}

func (p *PageVO[T]) SetTotalCountAndCalculateTotalPages(totalCount int64) {
	if p.PageSize == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = (totalCount + p.PageSize - 1) / p.PageSize
	}
	p.TotalCount = totalCount
}

func NewPageVO[T any](pageNo, pageSize, totalCount int64, list []T) *PageVO[T] {
	pageVO := &PageVO[T]{Page: Page{PageNo: pageNo, PageSize: pageSize}, List: list}
	pageVO.SetTotalCountAndCalculateTotalPages(totalCount)
	return pageVO
}

func Wrap[T any](fn func(ctx *gin.Context) (T, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := fn(ctx)
		if err != nil {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func ErrorHandler(ctx *gin.Context, err error) {
	l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
	var e ErrorResponseBody
	var r ResponseBody[any]
	switch {
	case errors.As(err, &e):
		l.ErrorContext(ctx, e.Error())
		ctx.JSON(e.HttpCode, nil)
	case errors.As(err, &r):
		ctx.JSON(http.StatusOK, r)
	default:
		l.ErrorContext(ctx, err.Error())
		fmt.Printf("stack trace: \n%+v\n", err)
		ctx.JSON(http.StatusInternalServerError, nil)
	}
}

func WrapWithBody[T any, R any](fn func(ctx *gin.Context, req R) (T, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req R
		bodyErr := ctx.Bind(&req)
		if bodyErr != nil {
			ErrorHandler(ctx, bodyErr)
			return
		}
		result, err := fn(ctx, req)
		if err != nil {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
