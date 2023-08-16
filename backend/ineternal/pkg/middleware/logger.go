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
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"os"
	"time"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Any().(time.Time)
		a.Value = slog.StringValue(t.Format(time.DateTime))
	}
	return a
}}))

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(XRequestIDKey)
		start := time.Now()

		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		body, _ := ctx.GetRawData()
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		requestLogger := logger.With("Url", path, "Query", query, "Method", method, "Body", string(body), "XRequestId", rid)
		requestLogger.Info("REQUEST")
		mr := &MyResponse{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = mr

		ctx.Next()

		elapsedTime := time.Since(start)

		statusCode := ctx.Writer.Status()

		responseLogger := logger.With("Code", statusCode, "Body", mr.body.String(), "ElapseTime", elapsedTime, "XRequestId", rid)
		responseLogger.Info("RESPONSE")
	}
}

type MyResponse struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w MyResponse) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w MyResponse) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
