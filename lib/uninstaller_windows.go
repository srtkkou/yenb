/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

func (u *UninstallerImpl) IsImplemented() bool {
	return true
}

func (u *UninstallerImpl) Uninstall() (err error) {
	Log().Info("Uninstalled yenb settings.")
	return nil
}
