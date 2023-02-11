package db

import (
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

type Config struct {
	Id     xid.ID
	UserId xid.ID
	Value  json.RawMessage
}

func ConfigForApplication(ctx context.Context, pg *pgxpool.Conn) (*Config, error) {
	return ConfigByUserId(ctx, pg, xid.NilID())
}

func ConfigByUserId(ctx context.Context, pg *pgxpool.Conn, id xid.ID) (*Config, error) {
	config := new(Config)
	err := pgxscan.Get(ctx, pg, config, "SELECT * FROM config c WHERE c.user_id = $1", id)
	if err != nil {
		return nil, err
	}

	return config, nil
}
