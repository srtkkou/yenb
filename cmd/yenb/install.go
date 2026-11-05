/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package yenb

import "github.com/spf13/pflag"

type installCmd struct {
	RootCmd
	fs *pflag.FlagSet
}

func newInstallCmd(cmd *RootCmd) *installCmd {
	c := &installCmd{
		RootCmd: *cmd,
		fs:      pflag.NewFlagSet("install", pflag.ExitOnError),
	}
	return c
}

func (c *installCmd) usage() string {
	return ""
}

func (c *installCmd) exec(args []string) error {
	return nil
}
