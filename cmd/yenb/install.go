/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import (
	"errors"
	"fmt"
	"strings"
)

type installCmdBase struct {
	RootCmd
}

var (
	ErrInstallCmdParseArgs = errors.New("ErrInstallCmdParseArgs")
)

func (c *installCmdBase) usage() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "usage: yenb install <args>")
	fmt.Fprintln(&sb)
	c.writeVersion(&sb)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "args:")
	return sb.String()
}
