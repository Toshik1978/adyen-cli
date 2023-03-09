package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Toshik1978/csv2adyen/pkg/processor"
)

const (
	version = "Commit %s, built on %s"
)

var (
	Buildstamp = "undefined" //nolint:revive
	Commit     = "undefined" //nolint:revive
)

func main() {
	app := &cli.App{
		Name:     "csv2adyen",
		Version:  fmt.Sprintf(version, Commit, Buildstamp),
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Anton Krivenko",
				Email: "anton@krivenko.dev",
			},
		},
		Copyright: "(c) 2021 Anton Krivenko",
		Usage:     "Run to process Adyen's split configuration from CSV file",
		Action:    cli.ShowAppHelp,
		Commands: []*cli.Command{
			{
				Name:    "process",
				Aliases: []string{"p"},
				Usage:   "process split configuration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "csv",
						Required: true,
					},
					&cli.StringFlag{
						Name:        "url",
						DefaultText: "cal-test.adyen.com",
					},
					&cli.StringFlag{
						Name:     "key",
						Required: true,
					},
					&cli.BoolFlag{
						Name: "dry-run",
					},
				},
				Action: func(c *cli.Context) error {
					url := c.String("url")
					if url == "" {
						url = "cal-test.adyen.com"
					}
					p := processor.New(c.String("csv"), url, c.String("key"), c.Bool("dry-run"))
					return p.Run(context.Background())
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
