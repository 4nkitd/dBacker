package main

import (
	"github.com/4nkitd/dBacker/services"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var giteaFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "gitea-config",
		Usage:  "gitea config path",
		EnvVar: "GOGS_CONFIG",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "gitea-data",
		Usage:  "gitea data path",
		Value:  "/data",
		EnvVar: "GOGS_DATA",
	}),
}

var databaseFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-host",
		Usage:  "database host",
		EnvVar: "DATABASE_HOST",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-port",
		Usage:  "database port",
		EnvVar: "DATABASE_PORT",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-name",
		Usage:  "database name",
		EnvVar: "DATABASE_NAME",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-user",
		Usage:  "database user",
		EnvVar: "DATABASE_USER",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-password",
		Usage:  "database password",
		EnvVar: "DATABASE_PASSWORD",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-password-file",
		Usage:  "database password file",
		EnvVar: "DATABASE_PASSWORD_FILE",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "database-options",
		Usage:  "extra options to pass to database service",
		EnvVar: "DATABASE_OPTIONS",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "database-compress",
		Usage:  "compress sql with gzip",
		EnvVar: "DATABASE_COMPRESS",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "database-ignore-exit-code",
		Usage:  "ignore restore process exit code",
		EnvVar: "DATABASE_IGNORE_EXIT_CODE",
	}),
}

var postgresFlags = []cli.Flag{
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "postgres-custom",
		Usage:  "use custom format (always compressed), ignored when database name is not set",
		EnvVar: "POSTGRES_CUSTOM_FORMAT",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "postgres-drop",
		Usage:  "drop database before restoring it",
		EnvVar: "POSTGRES_DROP",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "postgres-owner",
		Usage:  "change owner on database restore",
		EnvVar: "POSTGRES_OWNER",
	}),
}

var tarballFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "tarball-path",
		Usage:  "path to backup/restore",
		EnvVar: "TARBALL_PATH_SOURCE",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "tarball-name",
		Usage:  "backup file prefix",
		EnvVar: "TARBALL_NAME_PREFIX",
	}),
	altsrc.NewBoolFlag(cli.BoolFlag{
		Name:   "tarball-compress",
		Usage:  "compress tarball with gzip",
		EnvVar: "TARBALL_COMPRESS",
	}),
}

func newGogsConfig(c *cli.Context) *services.GiteaConfig {
	c = c.Parent()

	return &services.GiteaConfig{
		ConfigPath: c.String("gitea-config"),
		DataPath:   c.String("gitea-data"),
		SaveDir:    c.GlobalString("savedir"),
	}
}

func newMysqlConfig(c *cli.Context) *services.MySQLConfig {
	c = c.Parent()

	return &services.MySQLConfig{
		Host:           c.String("database-host"),
		Port:           c.String("database-port"),
		User:           c.String("database-user"),
		Password:       fileOrString(c, "database-password"),
		Database:       c.String("database-name"),
		Options:        c.String("database-options"),
		Compress:       c.Bool("database-compress"),
		SaveDir:        c.GlobalString("savedir"),
		IgnoreExitCode: c.Bool("database-ignore-exit-code"),
	}
}

func newPostgresConfig(c *cli.Context) *services.PostgresConfig {
	c = c.Parent()

	return &services.PostgresConfig{
		Host:           c.String("database-host"),
		Port:           c.String("database-port"),
		User:           c.String("database-user"),
		Password:       fileOrString(c, "database-password"),
		Database:       c.String("database-name"),
		Options:        c.String("database-options"),
		Compress:       c.Bool("database-compress"),
		Custom:         c.Bool("postgres-custom"),
		SaveDir:        c.GlobalString("savedir"),
		IgnoreExitCode: c.Bool("database-ignore-exit-code"),
		Drop:           c.Bool("postgres-drop"),
		Owner:          c.String("postgres-owner"),
	}
}

func newTarballConfig(c *cli.Context) *services.TarballConfig {
	c = c.Parent()

	return &services.TarballConfig{
		Path:     c.String("tarball-path"),
		Name:     c.String("tarball-name"),
		Compress: c.Bool("tarball-compress"),
		SaveDir:  c.GlobalString("savedir"),
	}
}

func newConsulConfig(c *cli.Context) *services.ConsulConfig {
	c = c.Parent()

	return &services.ConsulConfig{
		SaveDir: c.GlobalString("savedir"),
	}
}

func giteaCmd(parent string) cli.Command {
	name := "gitea"
	return cli.Command{
		Name:   name,
		Usage:  "connect to gitea service",
		Flags:  giteaFlags,
		Before: applyConfigValues(giteaFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func postgresCmd(parent string) cli.Command {
	name := "postgres"
	flags := append(databaseFlags, postgresFlags...)
	return cli.Command{
		Name:   name,
		Usage:  "connect to postgres service",
		Flags:  flags,
		Before: applyConfigValues(flags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func mysqlCmd(parent string) cli.Command {
	name := "mysql"
	return cli.Command{
		Name:   name,
		Usage:  "connect to mysql service",
		Flags:  databaseFlags,
		Before: applyConfigValues(databaseFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func tarballCmd(parent string) cli.Command {
	name := "tarball"
	return cli.Command{
		Name:   name,
		Usage:  "connect to tarball service",
		Flags:  tarballFlags,
		Before: applyConfigValues(tarballFlags),
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}

func consulCmd(parent string) cli.Command {
	name := "consul"
	return cli.Command{
		Name:  name,
		Usage: "connect to consul service",
		Subcommands: []cli.Command{
			s3Cmd(parent, name),
			filesystemCmd(parent, name),
		},
	}
}
