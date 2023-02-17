// slog_request_logger.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package util

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cdr.dev/slog"

	"github.com/go-chi/chi/middleware"
)

var (
	_ middleware.LogFormatter = (*SlogChiFormatter)(nil)
	_ middleware.LogEntry     = (*SlogChiLogEntry)(nil)
)

type SlogChiFormatter struct {
	logger slog.Logger
}

func NewSlogChiFormatter(logger slog.Logger) *SlogChiFormatter {
	return &SlogChiFormatter{logger: logger}
}

func (formatter *SlogChiFormatter) NewLogEntry(request *http.Request) middleware.LogEntry {
	fields := []slog.Field{
		slog.F("request_host", request.Host),
		slog.F("request_uri", request.RequestURI),
		slog.F("request_method", request.Method),
		slog.F("remote_addr", request.RemoteAddr),
	}

	if reqID := middleware.GetReqID(request.Context()); reqID != "" {
		fields = append(fields, slog.F("request_id", reqID))
	}

	return &SlogChiLogEntry{
		logger:  formatter.logger.With(fields...),
		request: request,
		panic:   false,
	}
}

type SlogChiLogEntry struct {
	logger  slog.Logger
	request *http.Request
	panic   bool
	fields  []slog.Field
}

func (logEntry *SlogChiLogEntry) Panic(v interface{}, stack []byte) {
	err, ok := v.(error)
	if !ok {
		err = fmt.Errorf("%+v", v)
	}

	logEntry.panic = true
	logEntry.logger = logEntry.logger.With(slog.Error(err), slog.F("stack", string(stack)))
}

func (logEntry *SlogChiLogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ any) {
	fields := []slog.Field{
		slog.F("response_status", status),
		slog.F("response_text", http.StatusText(status)),
		slog.F("response_size", uint64(bytes)),
		slog.F("response_duration", elapsed),
	}

	if !logEntry.panic {
		if !logEntry.panic {
			logEntry.logger.Info(context.Background(), "request complete", fields...)
		} else {
			logEntry.logger.Error(context.Background(), "request complete", fields...)
		}
	}
}
