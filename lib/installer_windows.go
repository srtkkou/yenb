/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	tmpl "text/template"

	"github.com/goark/errs"

	"github.com/srtkkou/yenb/template"
)

type InstallerImpl struct {
	InstallerBase
	psProfilePath string
}

const (
	psScriptName = "yenb.ps1"
	psTmplPath   = "yenb.ps1.tmpl"
)

var (
	ErrInstallerWindowsGetPowerShellProfilePath = errors.New("ErrInstallerWindowsGetPowerShellProfilePath")
)

func (i *InstallerImpl) IsImplemented() bool {
	return true
}

func (i *InstallerImpl) Install() (err error) {
	// Create script directory.
	if err = i.createScriptDir(); err != nil {
		return err
	}
	// Create PowerShell script.
	if err = i.createPSScript(); err != nil {
		return err
	}
	// Get PowerShell profile path.
	out, err := exec.Command("PowerShell.exe",
		"-Command", `$profile`).Output()
	if err != nil {
		return errs.Wrap(ErrInstallerWindowsGetPowerShellProfilePath,
			errs.WithCause(err))
	}
	// Sanitize output.
	i.psProfilePath = strings.TrimSpace(string(out))
	Log().Info("PowerShell profile path",
		slog.String("path", i.psProfilePath))
	// TODO: Append line to profile.
	// TODO: Write batch autorun path in registry.
	return nil
}

func (i *InstallerImpl) createPSScript() (err error) {
	// Initialize template.
	psTmpl, err := tmpl.ParseFS(template.Files, psTmplPath)
	if err != nil {
		return errs.Wrap(ErrInstallerTmplParse,
			errs.WithCause(err))
	}
	// Open file.
	psPath := filepath.Join(i.scriptDir, psScriptName)
	psFile, err := os.OpenFile(psPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return errs.Wrap(ErrInstallerOpenFile, errs.WithCause(err),
			errs.WithContext("path", psPath))
	}
	defer psFile.Close()
	// Write to file.
	if err = psTmpl.Execute(psFile, nil); err != nil {
		return errs.Wrap(ErrInstallerTmplExecute,
			errs.WithCause(err))
	}
	return nil
}
