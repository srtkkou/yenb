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

type Uninstaller interface {
	IsImplemented() bool
	Uninstall() error
}

type (
	UninstallerBase struct {
		scriptDir string
	}

	UninstallerImpl struct {
		UninstallerBase
	}
)

var (
	ErrUninstallerNotImplemented = errors.New("ErrUninstallerNotImplemented")
	ErrUninstallerRemovePath     = errors.New("ErrUninstallerRemovePath")
)

func NewUninstaller() (Uninstaller, error) {
	u := &UninstallerImpl{}
	if !u.IsImplemented() {
		err := errs.Wrap(ErrUninstallerNotImplemented)
		return nil, err
	}
	// Set script dir.
	config := Config()
	u.scriptDir = filepath.Join(
		filepath.Dir(config.ExePath()), scriptDirName)
	return u, nil
}

func (u *UninstallerBase) IsImplemented() bool {
	return false
}

func (u *UninstallerBase) deleteScriptDir() (err error) {
	if _, err = os.Stat(u.scriptDir); os.IsExist(err) {
		if err = os.RemoveAll(u.scriptDir); err != nil {
			return errs.Wrap(ErrUninstallerRemovePath,
				errs.WithCause(err),
				errs.WithContext("path", u.scriptDir))
		}
	}
	Log().Info("Deleted script directory.",
		slog.String("path", u.scriptDir))
	return nil
}
