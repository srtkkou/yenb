/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
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
