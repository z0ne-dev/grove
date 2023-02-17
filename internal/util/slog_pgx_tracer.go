// slog_pgx_tracer.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package util

import (
	"context"

	"cdr.dev/slog"
	"github.com/jackc/pgx/v5"
)

var _ pgx.QueryTracer = (*SlogPgxTracer)(nil)

var (
	contextKeySQL  = struct{}{}
	contextKeyArgs = struct{}{}
)

// SlogPgxTracer is a pgx.QueryTracer that logs queries using slog.
type SlogPgxTracer struct {
	logger *slog.Logger
}

// TraceQueryStart is called when a sql query is started.
func (_ *SlogPgxTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	ctx = context.WithValue(ctx, contextKeySQL, data.SQL)
	ctx = context.WithValue(ctx, contextKeyArgs, data.Args)

	return ctx
}

// TraceQueryEnd is called when a sql query has completed.
func (s *SlogPgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	logger := s.logger.With(slog.F("sql", ctx.Value(contextKeySQL)), slog.F("args", ctx.Value(contextKeySQL)))
	if data.Err != nil {
		logger.Error(ctx, "pgx query error", slog.F("command_tag", data.CommandTag), slog.F("err", data.Err))
		return
	}
	logger.Debug(ctx, "pgx query", slog.F("command_tag", data.CommandTag))
}

// NewSlogPgxTracer creates a new SlogPgxTracer.
func NewSlogPgxTracer(logger *slog.Logger) *SlogPgxTracer {
	return &SlogPgxTracer{
		logger: logger,
	}
}
