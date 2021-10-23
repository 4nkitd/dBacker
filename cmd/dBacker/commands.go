package main

import (
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var defaultFlags = []cli.Flag{
	altsrc.NewIntFlag(cli.IntFlag{
		Name:   "random-delay",
		Usage:  "schedule random delay",
		Value:  1,
		EnvVar: "SCHEDULE_RANDOM_DELAY",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "savedir",
		Usage:  "directory to save/read backups",
		Value:  "/tmp",
		EnvVar: "SAVE_DIR",
	}),
}

var backupFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "schedule",
		Usage:  "cron schedule",
		Value:  "@daily",
		EnvVar: "SCHEDULE",
	}),
	altsrc.NewIntFlag(cli.IntFlag{
		Name:   "max-backups",
		Usage:  "max backups to keep (0 to disable the feature)",
		Value:  5,
		EnvVar: "MAX_BACKUPS",
	}),
}

var restoreFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "schedule",
		Usage:  "cron schedule",
		Value:  "none",
		EnvVar: "SCHEDULE",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "restore-file",
		Usage:  "restore from this file instead of searching for the most recent",
		EnvVar: "RESTORE_FILE",
	}),
}

func backupCmd() cli.Command {
	name := "backup"
	flags := append(defaultFlags, backupFlags...)
	return cli.Command{
		Name:   name,
		Usage:  "run a backup task",
		Flags:  flags,
		Before: applyConfigValues(flags),
		Subcommands: []cli.Command{
			giteaCmd(name),
			postgresCmd(name),
			mysqlCmd(name),
			tarballCmd(name),
			consulCmd(name),
		},
	}
}

func restoreCmd() cli.Command {
	name := "restore"
	flags := append(defaultFlags, restoreFlags...)
	return cli.Command{
		Name:   "restore",
		Usage:  "run a restore task",
		Flags:  flags,
		Before: applyConfigValues(flags),
		Subcommands: []cli.Command{
			giteaCmd(name),
			postgresCmd(name),
			mysqlCmd(name),
			tarballCmd(name),
			consulCmd(name),
		},
	}
}
