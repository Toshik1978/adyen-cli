//nolint:dupl
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Toshik1978/csv2adyen/pkg/commands"
	"github.com/Toshik1978/csv2adyen/pkg/commands/cellular"
	"github.com/Toshik1978/csv2adyen/pkg/commands/close"
	"github.com/Toshik1978/csv2adyen/pkg/commands/install"
	"github.com/Toshik1978/csv2adyen/pkg/commands/link"
	"github.com/Toshik1978/csv2adyen/pkg/commands/method"
	"github.com/Toshik1978/csv2adyen/pkg/commands/offline"
	"github.com/Toshik1978/csv2adyen/pkg/commands/reassign"
)

const (
	version = "%s, built %s"
)

var (
	Buildstamp = "undefined" //nolint:revive
	Commit     = "undefined" //nolint:revive
)

func main() {
	logger, err := newLogger()
	if err != nil {
		log.Fatal("Failed to create the logger", err)
	}
	client := newHTTPClient()
	config := newConfig(logger)
	app := newApp(logger, client, config)

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to run the application: ", err)
	}
}

// newLogger initializes logger for console.
func newLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Encoding = "console"
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}

// newHTTPClient initializes HTTP client.
func newHTTPClient() *http.Client {
	return http.DefaultClient
}

// newConfig initializes new configuration.
func newConfig(logger *zap.Logger) *commands.Config {
	var config commands.Config
	if err := env.Parse(&config); err != nil {
		logger.
			With(zap.Error(err)).
			Fatal("Failed to initialize configuration from the environment")
	}
	return &config
}

// newApp initializes new application.
func newApp(logger *zap.Logger, client *http.Client, config *commands.Config) *cli.App { //nolint:funlen
	return &cli.App{
		Name:     "adyen-cli",
		Version:  fmt.Sprintf(version, Commit, Buildstamp),
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Anton Krivenko",
				Email: "anton@krivenko.dev",
			},
		},
		Copyright: "(c) 2023 Anton Krivenko",
		Usage:     "Operate with your Adyen account via CLI",
		Action:    cli.ShowAppHelp,
		Commands: []*cli.Command{
			{
				Name:    "link",
				Aliases: []string{"l"},
				Usage:   "Link split configurations to stores",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the required data to link",
					},
					&cli.BoolFlag{
						Name:  "balance",
						Usage: "use this parameter if you want to link on Balance Platform",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := link.New(
						logger, client, config,
						c.String("csv"), c.Bool("balance"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
			{
				Name:    "close",
				Aliases: []string{"c"},
				Usage:   "Close accounts",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the required data to close",
					},
					&cli.BoolFlag{
						Name:  "store",
						Usage: "use this parameter if you want to close store too",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := close.New(
						logger, client, config,
						c.String("csv"), c.Bool("store"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
			{
				Name:    "methods",
				Aliases: []string{"m"},
				Usage:   "Add payment methods",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the required data to add",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := method.New(
						logger, client, config,
						c.String("csv"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
			{
				Name:    "reassign",
				Aliases: []string{"r"},
				Usage:   "Reassign terminals",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the required data to reassign",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := reassign.New(
						logger, client, config,
						c.String("csv"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
			{
				Name:    "cellular",
				Aliases: []string{"e"},
				Usage:   "Enable/disable cellular on terminals",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the terminal IDs",
					},
					&cli.BoolFlag{
						Name:  "disable",
						Usage: "use this parameter if you want to disable cellular",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := cellular.New(
						logger, client, config,
						c.String("csv"), c.Bool("disable"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
			{
				Name:    "offline",
				Aliases: []string{"o"},
				Usage:   "Disable offline payments on terminals",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the terminal IDs",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := offline.New(
						logger, client, config,
						c.String("csv"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			}, {
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "Install apps on terminals",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "csv",
						Required:  true,
						TakesFile: true,
						Usage:     "the full path to CSV file, containing the terminal IDs",
					},
					&cli.BoolFlag{
						Name:  "prod",
						Usage: "use this parameter if you want to run on production environment",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "use this parameter if you want to do dry run (no changes will apply)",
					},
				},
				Action: func(c *cli.Context) error {
					p := install.New(
						logger, client, config,
						c.String("csv"), c.Bool("prod"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
		},
	}
}
