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

type SlogPgxTracer struct {
	logger slog.Logger
}

func (s *SlogPgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, "sql", data.SQL)
	ctx = context.WithValue(ctx, "args", data.Args)

	return ctx
}

func (s *SlogPgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	logger := s.logger.With(slog.F("sql", ctx.Value("sql")), slog.F("args", ctx.Value("args")))
	if data.Err != nil {
		logger.Error(ctx, "pgx query error", slog.F("command_tag", data.CommandTag), slog.F("err", data.Err))
		return
	}
	logger.Debug(ctx, "pgx query", slog.F("command_tag", data.CommandTag))
}

func NewSlogPgxTracer(logger slog.Logger) *SlogPgxTracer {
	return &SlogPgxTracer{
		logger: logger,
	}
}
