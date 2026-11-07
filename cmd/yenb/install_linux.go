/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"github.com/goark/errs"

	"github.com/srtkkou/yenb/lib"
)

type installCmd struct {
	installCmdBase
}

func newInstallCmd(cmd *RootCmd) subCmd {
	c := &installCmd{}
	c.RootCmd = *cmd
	return c
}

func (c *installCmd) exec(args []string) error {
	// Parse args.
	if err := c.fs.Parse(args); err != nil {
		return errs.Wrap(ErrInstallCmdParseArgs)
	}
	lib.Log().InfoLog("Installed yenb.")
	return nil
}
