package services

import (
	"fmt"
)

// ConsulConfig has the config options for the ConsulConfig service
type ConsulConfig struct {
	SaveDir string
}

// ConsulAppPath points to the consul binary location
var ConsulAppPath = "/bin/consul"

// Backup generates a tarball of the consul database and returns the path where is stored
func (c *ConsulConfig) Backup() (string, error) {
	filepath := generateFilename(c.SaveDir, "consul-backup") + ".snap"
	args := []string{"snapshot", "save", filepath}

	app := CmdConfig{}

	if err := app.CmdRun(ConsulAppPath, args...); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", ConsulAppPath, err)
	}

	return filepath, nil
}

// Restore takes a GiteaConfig backup and restores it to the service
func (c *ConsulConfig) Restore(filepath string) error {
	args := []string{"snapshot", "restore", filepath}

	app := CmdConfig{}

	if err := app.CmdRun(ConsulAppPath, args...); err != nil {
		return fmt.Errorf("couldn't execute consul restore, %v", err)
	}

	return nil
}
