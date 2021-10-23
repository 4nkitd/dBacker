package services

import (
	"fmt"
	"path"

	"github.com/mholt/archiver/v3"
)

// TarballConfig has the config options for the TarballConfig service
type TarballConfig struct {
	Name     string
	Path     string
	Compress bool
	SaveDir  string
}

// Backup creates a tarball of the specified directory
func (f *TarballConfig) Backup() (string, error) {
	var name string
	if f.Name != "" {
		name = f.Name + "-backup"
	} else {
		name = path.Base(f.Path) + "-backup"
	}

	filepath := generateFilename(f.SaveDir, name) + ".tar"

	if f.Compress {
		filepath += ".gz"
	}

	err := archiver.Archive([]string{f.Path}, filepath)
	if err != nil {
		return "", fmt.Errorf("cannot create tarball on %s, %v", filepath, err)
	}

	return filepath, nil
}

// Restore extracts a tarball to the specified directory
func (f *TarballConfig) Restore(filepath string) error {
	err := removeDirectoryContents(f.Path)
	if err != nil {
		return fmt.Errorf("failed to empty directory contents before restoring: %v", err)
	}

	err = archiver.Unarchive(filepath, path.Dir(f.Path))
	if err != nil {
		return fmt.Errorf("cannot unpack backup: %v", err)
	}

	return nil
}
