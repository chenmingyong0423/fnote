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

package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	XRequestIDKey = "X-Request-ID"
)

func RequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(XRequestIDKey)
		if rid == "" {
			rid = uuid.NewV4().String()
			ctx.Request.Header.Set(XRequestIDKey, rid)
			ctx.Set(XRequestIDKey, rid)
		}
		ctx.Writer.Header().Set(XRequestIDKey, rid)
		ctx.Next()
	}
}
