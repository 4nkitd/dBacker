package services

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
)

// GiteaConfig has the config options for the GiteaConfig service
type GiteaConfig struct {
	ConfigPath string
	DataPath   string
	SaveDir    string
}

// GiteaAppPath points to the gitea binary location
var GiteaAppPath = "/app/gitea/gitea"

func (g *GiteaConfig) newGiteaCmd() *CmdConfig {
	uid := uint32(getEnvInt("PUID", 1000))
	gid := uint32(getEnvInt("PGID", 1000))
	creds := &syscall.Credential{Uid: uid, Gid: gid}

	env := os.Environ()
	home := fmt.Sprintf("HOME=%s", path.Join(g.DataPath, "git"))
	env = append(env, "USER=git", home)

	return &CmdConfig{
		OutputFile: os.Stdout,
		Env:        env,
		Credential: creds,
		WorkDir:    g.SaveDir,
	}
}

// Backup generates a tarball of the GiteaConfig repositories and returns the path where is stored
func (g *GiteaConfig) Backup() (string, error) {
	filename := generateFilename("", "gitea-dump") + ".zip"
	args := []string{"dump", "--skip-log", "--tempdir", g.SaveDir, "--file", filename}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	app := g.newGiteaCmd()

	if err := app.CmdRun(GiteaAppPath, args...); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", GiteaAppPath, err)
	}

	return path.Join(g.SaveDir, filename), nil
}

// Restore takes a GiteaConfig backup and restores it to the service
func (g *GiteaConfig) Restore(_ string) error {
	return errors.New("not implemented")
}
