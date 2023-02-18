// slog_request_logger.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package util

import (
	"net/http"
	"time"

	"github.com/z0ne-dev/grove/lib/slogz"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

var (
	_ middleware.LogFormatter = (*SlogChiFormatter)(nil)
	_ middleware.LogEntry     = (*SlogChiLogEntry)(nil)
)

// SlogChiFormatter is a go-chi/chi/middleware.LogFormatter implementation.
type SlogChiFormatter struct {
	logger *slog.Logger
}

// NewSlogChiFormatter creates a new SlogChiFormatter.
func NewSlogChiFormatter(logger *slog.Logger) *SlogChiFormatter {
	return &SlogChiFormatter{logger: logger}
}

// NewLogEntry is called when a request is received.
func (formatter *SlogChiFormatter) NewLogEntry(request *http.Request) middleware.LogEntry {
	logger := formatter.logger.With(
		slog.String("request_host", request.Host),
		slog.String("request_uri", request.RequestURI),
		slog.String("request_method", request.Method),
		slog.String("remote_addr", request.RemoteAddr),
	)

	if reqID := middleware.GetReqID(request.Context()); reqID != "" {
		logger.With(slog.String("request_id", reqID))
	}

	return &SlogChiLogEntry{
		logger:  logger,
		request: request,
	}
}

// SlogChiLogEntry is a go-chi/chi/middleware.LogEntry implementation.
type SlogChiLogEntry struct {
	request *http.Request
	logger  *slog.Logger
	err     error
}

// Panic is called when a panic occurs.
func (logEntry *SlogChiLogEntry) Panic(v any, stack []byte) {
	err, ok := v.(error)
	if !ok {
		panic(v)
	}

	logEntry.err = err
	logEntry.logger = logEntry.logger.With(
		slog.String("stack", string(stack)),
	)
}

//revive:disable:argument-limit Interface required by go-chi/chi/middleware
func (logEntry *SlogChiLogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ any) {
	//revive:enable:argument-limit
	logger := logEntry.logger.With(
		slog.Int("response_status", status),
		slog.String("response_text", http.StatusText(status)),
		slog.Int("response_size", bytes),
		slog.Duration("response_duration", elapsed),
	)

	if logEntry.err != nil {
		logger.Error("request error", logEntry.err, slogz.Stringer("err", logEntry.err))
	} else {
		logger.Info("request")
	}
}
