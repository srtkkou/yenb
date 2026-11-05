/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goark/errs"

	"github.com/srtkkou/yenb/lib"
)

type uninstallCmd struct {
	RootCmd
}

var (
	ErrUninstallCmdParseArgs = errors.New("ErrUninstallCmdParseArgs")
)

func newUninstallCmd(cmd *RootCmd) *uninstallCmd {
	c := &uninstallCmd{
		RootCmd: *cmd,
	}
	return c
}

func (c *uninstallCmd) usage() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "usage: yenb uninstall <options>")
	fmt.Fprintln(&sb)
	c.writeVersion(&sb)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "options:")
	fmt.Fprintln(&sb, c.fs.FlagUsages())
	return sb.String()
}

func (c *uninstallCmd) exec(args []string) error {
	// Parse args.
	if err := c.fs.Parse(args); err != nil {
		return errs.Wrap(ErrInstallCmdParseArgs)
	}
	// Initialize uninstaller.
	u, err := lib.NewUninstaller()
	if err != nil {
		return err
	}
	// Run uninstaller.
	if err = u.Uninstall(); err != nil {
		return err
	}
	return nil
}
