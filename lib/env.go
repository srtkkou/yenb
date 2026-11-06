/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"

	"github.com/goark/errs"
)

type Env interface {
	IsImplemented() bool
	SetUserEnvs() error
}

type EnvBase struct{}

var (
	ErrEnvNotImplemented = errors.New("ErrEnvNotImplemented")
)

func NewEnv() (Env, error) {
	e := &EnvImpl{}
	if !e.IsImplemented() {
		err := errs.Wrap(ErrEnvNotImplemented)
		return nil, err
	}
	return e, nil
}

func (e *EnvBase) IsImplemented() bool {
	return false
}
