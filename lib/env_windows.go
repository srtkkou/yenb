/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goark/errs"
	reg "golang.org/x/sys/windows/registry"
)

type EnvImpl struct {
	EnvBase
}

const (
	USER_ENV_PATH = "Environment"
)

var (
	ErrEnvWinOpenKey     = errors.New("ErrEnvWinOpenKey")
	ErrEnvWinSetExpandSZ = errors.New("ErrEnvWinSetExpandSZ")
)

func (e *EnvImpl) IsImplemented() bool {
	return true
}

func (e *EnvImpl) SetUserEnvs() error {
	// Open registry path.
	regKey, err := reg.OpenKey(reg.CURRENT_USER, USER_ENV_PATH, reg.SET_VALUE)
	if err != nil {
		return errs.Wrap(ErrEnvWinOpenKey, errs.WithCause(err),
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
			// Sanitize path.
			value := v.Value
			if config.OS() == "windows" {
				value = strings.ReplaceAll(value, "/", `\`)
				value = strings.ReplaceAll(value, "${", `%`)
				value = strings.ReplaceAll(value, "}", `%`)
			}
			// Set value.
			err = regKey.SetExpandStringValue(k, value)
			if err != nil {
				return errs.Wrap(ErrEnvWinSetExpandSZ, errs.WithCause(err),
					errs.WithContext("name", k),
					errs.WithContext("value", value))
			}
			Log().InfoContext(config.Context(), "SetExpandStringValue",
				slog.String("name", k),
				slog.String("type", v.Type),
				slog.String("value", value))
		}
	}
	return nil
}
