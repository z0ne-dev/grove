// routes_middleware.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package application

import (
	"encoding/json"
	"net/http"

	"grove/internal/db"

	"github.com/CloudyKit/jet/v6"
	"github.com/creasty/defaults"
	"github.com/jackc/pgx/v5/pgxpool"
)

func jetGlobalsMiddleware(set *jet.Set, pool *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fetch a connection from the pool
			conn, err := pool.Acquire(r.Context())
			if err != nil {
				panic(err)
			}

			defer conn.Release()

			rawConfig, err := db.ConfigForApplication(r.Context(), conn)
			if err != nil {
				panic(err)
			}

			config := new(db.ConfigPayloadApplication)
			err = json.Unmarshal(rawConfig.Value, &config)
			if err != nil {
				panic(err)
			}

			defaults.MustSet(config)
			set.AddGlobal("config", config)
			next.ServeHTTP(w, r)
		})
	}
}
