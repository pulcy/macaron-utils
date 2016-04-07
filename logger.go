// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"time"

	"github.com/op/go-logging"
	"gopkg.in/macaron.v1"
)

// LoggerOption is used to control the logging process
type LoggerOption func(ctx *macaron.Context) bool

// Logger creates a handler that logs the current request.
func Logger(log *logging.Logger, options ...LoggerOption) macaron.Handler {
	return func(ctx *macaron.Context) {
		start := time.Now()

		rw := ctx.Resp.(macaron.ResponseWriter)
		ctx.Next()

		for _, opt := range options {
			if !opt(ctx) {
				return
			}
		}

		ms := int(time.Since(start) / time.Millisecond)
		log.Infof("%s %s %d %d %d", ctx.Req.Method, ctx.Req.RequestURI, rw.Status(), rw.Size(), ms)
	}
}

func DontLogHead() LoggerOption {
	return func(ctx *macaron.Context) bool {
		if ctx.Req.Method == "HEAD" {
			return false
		}
		return true
	}
}
