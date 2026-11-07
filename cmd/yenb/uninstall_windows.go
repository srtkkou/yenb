/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"github.com/goark/errs"

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
	lib.Log().InfoLog("Uninstalled yenb.")
	return nil
}
