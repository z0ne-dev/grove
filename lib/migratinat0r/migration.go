package migratinat0r

import (
	"context"
	"fmt"
)

type MigrationFunc func(context.Context, Commands) error

type Migration interface {
	fmt.Stringer
	Run(context.Context, Commands) error
}

type migrationFuncWrapper struct {
	name string
	fn   MigrationFunc
}

func (m *migrationFuncWrapper) Run(ctx context.Context, tx Commands) error {
	return m.fn(ctx, tx)
}

func (m *migrationFuncWrapper) String() string {
	return m.name
}
