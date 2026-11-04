/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"

	"github.com/goark/errs"
)

type EnvEditor interface {
	IsImplemented() bool
}

type EnvEditorBase struct{}

type EnvEditorImpl struct {
	EnvEditorBase
}

var (
	ErrEnvEditorNotImplemented = errors.New("ErrEnvEditorNotImplemented")
)

func NewEnvEditor() (EnvEditor, error) {
	e := &EnvEditorImpl{}
	if !e.IsImplemented() {
		err := errs.Wrap(ErrEnvEditorNotImplemented)
		return nil, err
	}
	return e, nil
}

func (e *EnvEditorBase) IsImplemented() bool {
	return false
}
