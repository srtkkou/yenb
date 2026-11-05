/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goark/errs"
)

type installCmd struct {
	RootCmd
}

var (
	ErrInstallCmdParseArgs = errors.New("ErrInstallCmdParseArgs")
)

func newInstallCmd(cmd *RootCmd) *installCmd {
	c := &installCmd{
		RootCmd: *cmd,
	}
	return c
}

func (c *installCmd) usage() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "usage: yenb install <args>")
	fmt.Fprintln(&sb)
	c.writeVersion(&sb)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "args:")
	return sb.String()
}

func (c *installCmd) exec(args []string) error {
	// Parse args.
	if err := c.fs.Parse(args); err != nil {
		return errs.Wrap(ErrInstallCmdParseArgs)
	}
	return nil
}
