/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/goark/errs"
)

type Installer interface {
	IsImplemented() bool
	Install() error
}

type InstallerBase struct {
	scriptDir string // Path to script directory.
}

const (
	scriptDirName = "script"
)

var (
	ErrInstallerNotImplemented  = errors.New("ErrInstallerNotImplemented")
	ErrInstallerCreateScriptDir = errors.New("ErrInstallerCreateScriptDir")
)

func NewInstaller() (Installer, error) {
	i := &InstallerImpl{}
	if !i.IsImplemented() {
		err := errs.Wrap(ErrUninstallerNotImplemented)
		return nil, err
	}
	// Set script dir.
	config := Config()
	i.scriptDir = filepath.Join(
		filepath.Dir(config.ExePath()), scriptDirName)
	return i, nil
}

func (i *InstallerBase) IsImplemented() bool {
	return false
}

// Create script directory.
func (i *InstallerBase) createScriptDir() (err error) {
	if _, err = os.Stat(i.scriptDir); os.IsNotExist(err) {
		if err = os.Mkdir(i.scriptDir, 0o755); err != nil {
			return errs.Wrap(ErrInstallerCreateScriptDir, errs.WithCause(err),
				errs.WithContext("path", i.scriptDir))
		}
		Log().Info("Created script dir.",
			slog.String("path", i.scriptDir))
	}
	return nil
}
