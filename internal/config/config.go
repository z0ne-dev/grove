// config.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package config

import "golang.org/x/exp/slog"

// Config is the configuration for the application.
type Config struct {
	Logging  Logging `json:"logging"`
	HTTP     HTTP    `json:"http"`
	Postgres string  `json:"postgres"`
}

// HTTP is the configuration for the http server.
type HTTP struct {
	Listen        string `default:":8421" json:"listen"`
	PublicAddress string `json:"public_address"`
}

// Logging is the configuration for the logging.
type Logging struct {
	Level slog.Level `default:"2" json:"level"`
}
