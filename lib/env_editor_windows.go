/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"log/slog"

	"github.com/goark/errs"
	reg "golang.org/x/sys/windows/registry"
)

const (
	USER_ENV_PATH = "Environment"
)

var (
	ErrEnvEditorWinOpenKey     = errors.New("ErrEnvEditorWinOpenKey")
	ErrEnvEditorWinSetExpandSZ = errors.New("ErrEnvEditorWinSetExpandSZ")
)

func (e *EnvEditorImpl) IsImplemented() bool {
	return true
}

func (e *EnvEditorImpl) SetUserEnvs() error {
	// Open registry path.
	regKey, err := reg.OpenKey(reg.CURRENT_USER, USER_ENV_PATH, reg.SET_VALUE)
	if err != nil {
		return errs.Wrap(ErrEnvEditorWinOpenKey, errs.WithCause(err),
			errs.WithContext("key", reg.CURRENT_USER),
			errs.WithContext("path", USER_ENV_PATH),
			errs.WithContext("access", reg.SET_VALUE))
	}
	defer regKey.Close()
	// Loop for each user env vars.
	config := Config()
	for k, v := range config.UserEnvs() {
		switch v.Type {
		case "REG_EXPAND_SZ":
			// Set value.
			err = regKey.SetExpandStringValue(k, v.Value)
			if err != nil {
				return errs.Wrap(ErrEnvEditorWinSetExpandSZ, errs.WithCause(err),
					errs.WithContext("name", k),
					errs.WithContext("value", v.Value))
			}
			Log().InfoContext(config.Context(), "SetExpandStringValue",
				slog.String("name", k),
				slog.String("type", v.Type),
				slog.String("value", v.Value))
		}
	}
	return nil
}
