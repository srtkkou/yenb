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

type installCmd struct {
	installCmdBase
}

const (
	shimDirName     = "shim"
	regUserEnvPath  = "Environment"
	regYenbHomeName = "YENB_HOME"
)

func newInstallCmd(cmd *RootCmd) subCmd {
	c := &installCmd{}
	c.RootCmd = *cmd
	return c
}

func (c *installCmd) exec(args []string) (err error) {
	// Parse args.
	if err = c.fs.Parse(args); err != nil {
		return errs.Wrap(ErrInstallCmdParseArgs)
	}
	// Create shim directory.
	config := lib.Config()
	yenbDir := filepath.Dir(config.ExePath())
	shimDir := filepath.Join(yenbDir, shimDirName)
	if err = lib.MkdirAll(shimDir); err != nil {
		return err
	}
	// Add env var YENB_HOME.
	err = lib.SetRegExpandSZ(reg.CURRENT_USER,
		regUserEnvPath, regYenbHomeName, yenbDir)
	if err != nil {
		return err
	}
	lib.Log().InfoLog("Installed yenb.")
	return nil
}
