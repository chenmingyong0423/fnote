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

package ioc

import (
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg *Config) io.Writer {
	writers := make([]io.Writer, 0, 2)
	writers = append(writers, os.Stdout)
	if cfg.Logger.Filename != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   cfg.Logger.Filename,
			MaxSize:    cfg.Logger.MaxSize,
			MaxAge:     cfg.Logger.MaxAge,
			MaxBackups: cfg.Logger.MaxBackups,
			LocalTime:  cfg.Logger.LocalTime,
			Compress:   cfg.Logger.Compress,
		})
	}
	multiWriter := io.MultiWriter(writers...)
	slog.SetDefault(slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Any().(time.Time)
				a.Value = slog.StringValue(t.Format(cfg.Logger.TimeFormat))
			}
			return a
		},
	})))
	return multiWriter
}
