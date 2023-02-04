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
