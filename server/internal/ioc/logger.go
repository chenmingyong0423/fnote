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

	"github.com/spf13/viper"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() io.Writer {
	writers := make([]io.Writer, 0, 2)
	writers = append(writers, os.Stdout)
	if viper.GetString("logger.file_name") != "" {
		writers = append(writers, &lumberjack.Logger{
			Filename:   viper.GetString("logger.file_name"),
			MaxSize:    viper.GetInt("logger.max_size"),
			MaxAge:     viper.GetInt("logger.max_age"),
			MaxBackups: viper.GetInt("logger.max_backups"),
			LocalTime:  viper.GetBool("logger.local_time"),
			Compress:   viper.GetBool("logger.compress"),
		})
	}
	multiWriter := io.MultiWriter(writers...)
	slog.SetDefault(slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Any().(time.Time)
				a.Value = slog.StringValue(t.Format(viper.GetString("logger.time_format")))
			}
			return a
		},
	})))
	return multiWriter
}
