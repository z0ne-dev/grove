// main.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/z0ne-dev/grove/internal/application"
	"github.com/z0ne-dev/grove/internal/config"
	"github.com/z0ne-dev/grove/internal/service"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogjson"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "grove",
	Short:   "grove activity pub server",
	Long:    `Interconnected activity pub server for the fediverse`,
	Version: "",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := defaults.Set(configuration); err != nil {
			return fmt.Errorf("failed to prepare configuration: %w", err)
		}

		slogger.Debug(context.Background(), "Loading Configuration")
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		if err := viper.Unmarshal(configuration, func(c *mapstructure.DecoderConfig) {
			c.WeaklyTypedInput = true
			c.Squash = true
			c.TagName = "json"
		}); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if configuration.Logging.EnableFile {
			logFile := configuration.Logging.File
			if !filepath.IsAbs(logFile) {
				wd, _ := os.Getwd()
				logFile, _ = filepath.Abs(filepath.Join(wd, logFile))
			}
			// #nosec G304
			f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}

			slogger = slogger.AppendSinks(slogjson.Sink(f))
		}

		slogger = slogger.Leveled(configuration.Logging.Level)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		slogger.Info(
			context.Background(),
			fmt.Sprintf("Starting %s", config.ApplicationName),
			slog.F("version", config.Version),
			slog.F("go_version", runtime.Version()),
			slog.F("os", runtime.GOOS),
			slog.F("arch", runtime.GOARCH),
		)

		c, err := service.New(slogger, configuration)
		if err != nil {
			return err
		}

		app := application.New(c)
		if err := app.MigrateDatabase(); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}
		if err := app.ConfigureRouter(); err != nil {
			return fmt.Errorf("failed to configure routes: %w", err)
		}
		app.ListenAndServe()
		return nil
	},
}

var (
	cfgFile       = ""
	configuration = new(config.Config)
	slogger       = slog.Make(sloghuman.Sink(os.Stdout))
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "custom configuration to load")
	rootCmd.PersistentFlags().CountP("verbose", "v", "Increase verboseness")
	bindFlag("verbose", "logging.level")
}

func bindFlag(flag, field string) {
	err := viper.BindPFlag(field, rootCmd.PersistentFlags().Lookup(flag))
	if err != nil {
		slogger.Critical(
			context.Background(),
			"Failed to bind flag to field",
			slog.F("flag", flag),
			slog.F("field", field),
			slog.Error(err),
		)
	}
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, _ := os.Getwd()
		viper.AddConfigPath(wd)           // look for config in the working directory
		viper.AddConfigPath("/etc/grove") // look for config in /etc
		viper.SetConfigName("grove")      // name of config file (without extension)
	}
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slogger.Fatal(context.Background(), "Failed to run", slog.Error(err))
	}
}

func main() {
	Execute()
}
