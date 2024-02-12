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

package ioc

import (
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/jwtutil"
	"github.com/gin-gonic/gin"
)

// JwtParseMiddleware jwt 解析中间件
func JwtParseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Request.RequestURI
		// 非 admin 接口不需要 jwt
		if !strings.HasPrefix(uri, "/admin") {
			ctx.Next()
			return
		}
		// 登录和初始化接口不需要 jwt
		if uri == "/admin/login" || uri == "/admin/init" {
			ctx.Next()
			return
		}

		jwtStr := ctx.GetHeader("Authorization")
		if jwtStr == "" {
			ctx.AbortWithStatusJSON(401, nil)
			return
		}
		// 解析 jwt
		claims, err := jwtutil.ParseJwt(jwtStr)
		if err != nil {
			ctx.AbortWithStatusJSON(401, nil)
			return
		}
		ctx.Set("jwtClaims", claims)
		ctx.Next()
	}
}
