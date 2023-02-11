package application

import (
	"cdr.dev/slog"
	"context"
	"fmt"
	"grove/internal/resource"
	"grove/lib/arrayz"
	"grove/lib/mgx"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func (a *app) MigrateDatabase() error {
	migrations, err := migrationsFrom(resource.Migrations, "/")
	if err != nil {
		return err
	}

	mig, err := mgx.New(mgx.Log(mgx.LoggerFunc(func(msg string, data map[string]any) {
		fields := make([]slog.Field, 0, len(data))
		for k, v := range data {
			fields = append(fields, slog.F(k, v))
		}

		a.container.Logger().Named("migrations").Info(context.Background(), msg, fields...)
	})), mgx.Migrations(
		migrations...,
	))
	if err != nil {
		return err
	}

	db, err := a.container.PostgresPool()
	if err != nil {
		return err
	}
	defer db.Close()

	c := context.Background()
	c, cancel := context.WithTimeout(c, 5*time.Minute)
	defer cancel()
	return mig.Migrate(c, db)
}

func migrationsFrom(httpfs http.FileSystem, p string) ([]mgx.Migration, error) {
	dir, err := httpfs.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open migrations directory: %w", err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	files = arrayz.Filter(files, func(file os.FileInfo) bool {
		return !file.IsDir() && strings.HasSuffix(file.Name(), ".sql")
	})
	files = arrayz.Sort(files, func(left fs.FileInfo, right fs.FileInfo) bool {
		return left.Name() < right.Name()

	})

	migrations := make([]mgx.Migration, 0, len(files))

	for _, file := range files {
		name := path.Base(file.Name())[0 : len(file.Name())-4]
		filePath := path.Join(p, file.Name())
		buffer, err := readFile(httpfs, filePath)
		if err != nil {
			return nil, err
		}
		sql := string(buffer)

		if strings.HasPrefix(sql, "--") {
			line := strings.SplitN(sql, "\n", 2)[0]
			name = strings.TrimSpace(strings.TrimPrefix(line, "--"))
		}

		mig := mgx.NewRawMigration(name, sql)
		migrations = append(migrations, mig)
	}

	return migrations, nil
}

func readFile(fs http.FileSystem, filePath string) ([]byte, error) {
	fileEntry, err := fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open migration file: %w", err)
	}

	stat, err := fileEntry.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat migration file: %w", err)
	}

	buffer := make([]byte, stat.Size())
	_, err = fileEntry.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration file: %w", err)
	}
	return buffer, nil
}
