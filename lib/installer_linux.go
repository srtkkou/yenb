/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

type InstallerImpl struct {
	InstallerBase
}

func (i *InstallerImpl) IsImplemented() bool {
	return true
}

func (i *InstallerImpl) Install() (err error) {
	// Create script directory.
	if err = i.createScriptDir(); err != nil {
		return err
	}
	return nil
}
