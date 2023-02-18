// grove.go Copyright (c) 2023 z0ne.
// All Rights Reserved.
// Licensed under the EUPL 1.2 License.
// See LICENSE the project root for license information.
//
// SPDX-License-Identifier: EUPL-1.2

package grove

import (
	"fmt"
	"os"
	"runtime"

	"github.com/z0ne-dev/grove/internal/application"
	"github.com/z0ne-dev/grove/internal/config"
	"github.com/z0ne-dev/grove/internal/service"

	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var rootCmd = &cobra.Command{
	Use:               "grove",
	Short:             "grove activity pub server",
	Long:              `Interconnected activity pub server for the fediverse`,
	Version:           "",
	PersistentPreRunE: preRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info(
			fmt.Sprintf("Starting %s", config.ApplicationName),
			slog.String("version", config.Version),
			slog.String("go_version", runtime.Version()),
			slog.String("os", runtime.GOOS),
			slog.String("arch", runtime.GOARCH),
		)

		c, err := service.New(slog.Default(), configuration)
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

func preRun(_ *cobra.Command, _ []string) error {
	// Set default values
	if err := defaults.Set(configuration); err != nil {
		return fmt.Errorf("failed to prepare configuration: %w", err)
	}

	slog.Debug("Loading Configuration")
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

	handlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     configuration.Logging.Level,
	}
	slog.SetDefault(slog.New(handlerOptions.NewTextHandler(os.Stdout)))
	return nil
}

var (
	cfgFile       = ""
	configuration = new(config.Config)
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
		slog.Default().Error(
			"Failed to bind flag to field",
			err,
			slog.String("flag", flag),
			slog.String("field", field),
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
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		slog.Default().Error("Failed to run", err)

		return 1
	}

	return 0
}
