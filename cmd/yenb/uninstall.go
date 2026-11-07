/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"errors"
	"fmt"
	"strings"
)

type uninstallCmdBase struct {
	RootCmd
}

var (
	ErrUninstallCmdParseArgs = errors.New("ErrUninstallCmdParseArgs")
)

func (c *uninstallCmdBase) usage() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "usage: yenb uninstall <options>")
	fmt.Fprintln(&sb)
	c.writeVersion(&sb)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "options:")
	fmt.Fprintln(&sb, c.fs.FlagUsages())
	return sb.String()
}
