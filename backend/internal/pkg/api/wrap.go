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

package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Wrap[T any](fn func(ctx *gin.Context) (T, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := fn(ctx)
		if err != nil {
			ErrorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, SuccessResponseWithData(result))
	}
}

func ErrorHandler(ctx *gin.Context, err error) {
	l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
	var e ErrorResponseBody
	switch {
	case errors.As(err, &e):
		l.ErrorContext(ctx, e.Error())
		ctx.JSON(e.HttpCode, nil)
	default:
		l.ErrorContext(ctx, fmt.Sprintf("%+v", err))
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
		ctx.JSON(http.StatusOK, SuccessResponseWithData(result))
	}
}
