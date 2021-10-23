package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	log "unknwon.dev/clog/v2"
)

const versionFormatter = `dBacker
Version:      %s
Git commit:   %s
Built:        %s
Compilation:  %s
`

func printVersion(c *cli.Context) {
	_, _ = fmt.Fprintf(c.App.Writer, versionFormatter, Version, Commit, BuildTime, BuildNumber)
}

func main() {
	app := cli.NewApp()
	app.Usage = "run backups from various services to S3-like storage"
	app.Version = Version
	cli.VersionPrinter = printVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "load config from yaml file",
			EnvVar: "CONFIG_FILE",
		},
	}

	app.Commands = []cli.Command{
		backupCmd(),
		restoreCmd(),
	}

	app.Before = func(c *cli.Context) error {
		if err := log.NewConsole(); err != nil {
			return err
		}

		log.Info("Starting dBacker, version: %s, commit: %s, built: %s, compilation: %s",
			Version,
			Commit,
			BuildTime,
			BuildNumber)

		if c.String("config") != "" {
			cfg, err := altsrc.NewYamlSourceFromFile(c.String("config"))
			if err != nil {
				app.Metadata = map[string]interface{}{
					"config": cfg,
				}
			}

			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Unrecoverable error: %v", err)
	}

	log.Info("Shutdown complete")
	log.Stop()
}
