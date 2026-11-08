/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"path/filepath"

	"github.com/goark/errs"
	reg "golang.org/x/sys/windows/registry"

	"github.com/srtkkou/yenb/lib"
)

type uninstallCmd struct {
	uninstallCmdBase
}

func newUninstallCmd(cmd *RootCmd) subCmd {
	c := &uninstallCmd{}
	c.RootCmd = *cmd
	return c
}

func (c *uninstallCmd) exec(args []string) (err error) {
	// Parse args.
	if err := c.fs.Parse(args); err != nil {
		return errs.Wrap(ErrUninstallCmdParseArgs)
	}
	// Delete env var YENB_HOME.
	err = lib.DeleteRegValue(reg.CURRENT_USER,
		regUserEnvPath, regYenbHomeName)
	if err != nil {
		return err
	}
	// Delete shim directory.
	config := lib.Config()
	shimDir := filepath.Join(
		filepath.Dir(config.ExePath()), shimDirName)
	if err = lib.MkdirAll(shimDir); err != nil {
		return err
	}
	lib.Log().InfoLog("Uninstalled yenb.")
	return nil
}
