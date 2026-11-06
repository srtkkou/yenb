/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/goark/errs"
)

type Installer interface {
	IsImplemented() bool
	Install() error
}

type (
	InstallerBase struct {
		scriptDir string // Path to script directory.
	}

	InstallerImpl struct {
		InstallerBase
	}
)

const (
	scriptDirName = "script"
)

var (
	ErrInstallerNotImplemented = errors.New("ErrInstallerNotImplemented")
	ErrInstallerTmplParse      = errors.New("ErrInstallerTmplParse")
	ErrInstallerTmplExecute    = errors.New("ErrInstallerTmplExecute")
	ErrInstallerMkdir          = errors.New("ErrInstallerMkdir")
	ErrInstallerOpenFile       = errors.New("ErrInstallerOpenFile")
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
			return errs.Wrap(ErrInstallerMkdir, errs.WithCause(err),
				errs.WithContext("path", i.scriptDir))
		}
	}
	return nil
}
