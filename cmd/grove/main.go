package main

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogjson"
	"context"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"grove/internal/application"
	"grove/internal/config"
	"grove/internal/service"
	"os"
	"path/filepath"
	"runtime"
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
			return fmt.Errorf("failed to migrate database: %v", err)
		}
		if err := app.ConfigureRouter(); err != nil {
			return fmt.Errorf("failed to configure routes: %v", err)
		}
		app.ListenAndServe()
		return nil
	},
}

var cfgFile = ""
var configuration = new(config.Config)
var slogger = slog.Make(sloghuman.Sink(os.Stdout))

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "custom configuration to load")
	rootCmd.PersistentFlags().CountP("verbose", "v", "Increase verboseness")
	bindFlag("verbose", "logging.level")
}

func bindFlag(flag string, field string) {
	err := viper.BindPFlag(field, rootCmd.PersistentFlags().Lookup(flag))
	if err != nil {
		slogger.Critical(context.Background(), "Failed to bind flag to field", slog.F("flag", flag), slog.F("field", field), slog.Error(err))
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slogger.Fatal(context.Background(), "Failed to run", slog.Error(err))
	}
}

func main() {
	Execute()
}
