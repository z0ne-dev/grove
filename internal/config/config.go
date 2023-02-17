// config.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package config

import "cdr.dev/slog"

type Config struct {
	Logging  Logging `json:"logging"`
	Http     Http    `json:"http"`
	Postgres string  `json:"postgres"`
}

type Http struct {
	Listen        string `default:":8421" json:"listen"`
	PublicAddress string `json:"public_address"`
}

type Logging struct {
	EnableFile bool       `json:"enable_file"`
	File       string     `default:"./logs" json:"file"`
	Level      slog.Level `default:"2" json:"level"`
}
