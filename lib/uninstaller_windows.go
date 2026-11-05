/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

func (u *UninstallerImpl) IsImplemented() bool {
	return true
}

func (u *UninstallerImpl) Uninstall() (err error) {
	// Delete script dir.
	if err = u.deleteScriptDir(); err != nil {
		return err
	}
	Log().Info("Uninstalled yenb settings.")
	return nil
}
