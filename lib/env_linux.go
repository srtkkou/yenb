/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

type EnvImpl struct {
	EnvBase
}

func (e *EnvImpl) IsImplemented() bool {
	return true
}

func (e *EnvImpl) SetUserEnvs() error {
	return nil
}
