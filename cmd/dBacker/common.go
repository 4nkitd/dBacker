package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/4nkitd/dBacker/services"
	"github.com/4nkitd/dBacker/stores"
	"github.com/robfig/cron/v3"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	log "unknwon.dev/clog/v2"
)

type task func(c *cli.Context) error

func getService(c *cli.Context, service string) services.Service {
	var config services.Service
	switch service {
	case "gitea":
		config = newGogsConfig(c)
	case "mysql":
		config = newMysqlConfig(c)
	case "postgres":
		config = newPostgresConfig(c)
	case "tarball":
		config = newTarballConfig(c)
	case "consul":
		config = newConsulConfig(c)
	default:
		log.Fatal("Unsupported service: %s", service)
	}

	return config
}

func getStore(c *cli.Context, store string) stores.Storer {
	var config stores.Storer
	switch store {
	case "s3":
		config = newS3Config(c)
	case "filesystem":
		config = newFilesystemConfig(c)
	default:
		log.Fatal("Unsupported store: %s", store)
	}

	return config
}

func runTask(c *cli.Context, command string, serviceName string, storeName string) error {
	service := getService(c, serviceName)
	store := getStore(c, storeName)

	switch command {
	case "backup":
		return runScheduler(c, func(c *cli.Context) error {
			return backupTask(c, service, store)
		})
	case "restore":
		return runScheduler(c, func(c *cli.Context) error {
			return restoreTask(c, service, store)
		})
	default:
		log.Fatal("Unsupported command: %s", command)
	}
	return nil
}

func backupTask(c *cli.Context, service services.Service, store stores.Storer) error {
	filepath, err := service.Backup()
	if err != nil {
		return fmt.Errorf("service backup failed: %v", err)
	}

	log.Trace("Backup saved to %s", filepath)

	filename := path.Base(filepath)

	if err = store.Store(filepath, filename); err != nil {
		return fmt.Errorf("couldn't upload file to store: %v", err)
	}

	err = store.RemoveOlderBackups(c.GlobalInt("max-backups"))
	if err != nil {
		return fmt.Errorf("couldn't remove old backups from store: %v", err)
	}

	return nil
}

func restoreTask(c *cli.Context, service services.Service, store stores.Storer) error {
	var err error
	var filename string

	if key := c.GlobalString("restore-file"); key != "" {
		// restore directly from this file
		filename = key
	} else {
		// find the latest file in the store
		filename, err = store.FindLatestBackup()
		if err != nil {
			return fmt.Errorf("cannot find the latest backup: %v", err)
		}
	}

	filepath, err := store.Retrieve(filename)
	if err != nil {
		return fmt.Errorf("cannot download file %s: %v", filename, err)
	}

	defer store.Close()

	if err = service.Restore(filepath); err != nil {
		return fmt.Errorf("service restore failed: %v", err)
	}

	return nil
}

func runScheduler(c *cli.Context, task task) error {
	cr := cron.New()
	schedule := c.GlobalString("schedule")

	if schedule == "" || schedule == "none" {
		log.Trace("Running task directly")
		return task(c)
	}

	log.Trace("Starting scheduled backup task")
	timeoutchan := make(chan bool, 1)

	_, err := cr.AddFunc(schedule, func() {
		delay := c.GlobalInt("random-delay")
		if delay <= 0 {
			log.Warn("Schedule random delay was set to a number <= 0, using 1 as default")
			delay = 1
		}

		seconds := rand.Intn(delay)

		// run immediately is no delay is configured
		if seconds == 0 {
			if err := task(c); err != nil {
				log.Error("Failed to run scheduled task: %v", err)
			}
			return
		}

		log.Trace("Waiting for %d seconds before starting scheduled job", seconds)

		select {
		case <-timeoutchan:
			log.Trace("Random timeout cancelled")
			break
		case <-time.After(time.Duration(seconds) * time.Second):
			log.Trace("Running scheduled task")

			if err := task(c); err != nil {
				log.Error("Failed to run scheduled task: %v", err)
			}
			break
		}
	})

	if err != nil {
		return err
	}

	cr.Start()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	timeoutchan <- true
	close(timeoutchan)

	log.Trace("Stopping scheduled task")
	ctx := cr.Stop()
	<-ctx.Done()

	return nil
}

func fileOrString(c *cli.Context, name string) string {
	if filepath := c.String(name + "-file"); filepath != "" {
		f, err := os.Open(filepath)
		if err != nil {
			log.Error("Cannot open password file: %v", err)
			return ""
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			return scanner.Text()
		}

		log.Warn("Using empty password file")
		return ""
	}

	return c.String(name)
}

func applyConfigValues(flags []cli.Flag) cli.BeforeFunc {
	return func(c *cli.Context) error {
		config := c.App.Metadata["config"]
		if config != nil {
			cfg, ok := config.(altsrc.InputSourceContext)
			if ok {
				return altsrc.ApplyInputSourceValues(c, cfg, flags)
			}

			return fmt.Errorf("invalid config type for metadata")
		}

		return nil
	}
}
