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
	psScriptPath  string
}

const (
	psScript = "yenb.ps1"
	psTmpl   = "yenb.ps1.tmpl"
)

var (
	ErrInstallerWindowsGetPSPolicyStdout   = errors.New("ErrInstallerWindowsGetPSPolicy")
	ErrInstallerWindowsGetPSPolicyCmdStart = errors.New("ErrInstallerWindowsGetPSPolicyCmdStart")
	ErrInstallerWindowsGetPSPolicyScan     = errors.New("ErrInstallerWindowsGetPSPolicyScan")
	ErrInstallerWindowsPSScriptTmplParse   = errors.New("ErrInstallerWindowsPSScriptTmplParse")
	ErrInstallerWindowsPSScriptOpen        = errors.New("ErrInstallerWindowsPSScriptOpen")
	ErrInstallerWindowsPSScirptTmplExecute = errors.New("ErrInstallerWindowsPSScirptTmplExecute")
	ErrInstallerWindowsGetPSProfilePath    = errors.New("ErrInstallerWindowsGetPSProfilePath")
	ErrInstallerWindowsCreatePSProfileDir  = errors.New("ErrInstallerWindowsCreatePSProfileDir")
	ErrInstallerWindowsCreatePSProfile     = errors.New("ErrInstallerWindowsCreatePSProfile")
	ErrInstallerWindowsUpdatePSProfile     = errors.New("ErrInstallerWindowsUpdatePSProfile")
	ErrInstallerWindowsScanPSProfile       = errors.New("ErrInstallerWindowsScanPSProfile")
	ErrInstallerWindowsSeekPSProfile       = errors.New("ErrInstallerWindowsSeekPSProfile")
)

func (i *InstallerImpl) IsImplemented() bool {
	return true
}

func (i *InstallerImpl) Install() (err error) {
	// Get PS script policy.
	/*
		policy, err := i.getPSPolicy()
		if err != nil {
			return err
		}
		Log().Info("PS Script execution policy", slog.String("policy", policy))
	*/
	// Create script directory.
	if err = i.createScriptDir(); err != nil {
		return err
	}
	// Get PowerShell script path.
	i.psScriptPath = filepath.Join(i.scriptDir, psScript)
	// Create PowerShell script.
	if err = i.createPSScript(); err != nil {
		return err
	}
	// Get PowerShell profile path.
	i.psProfilePath, err = i.getPSProfilePath()
	if err != nil {
		return err
	}
	/*
		// Create or update PS profile.
		if _, err = os.Stat(i.psProfilePath); os.IsNotExist(err) {
			if err = i.createPSProfile(); err != nil {
				return err
			}
		} else {
			if err = i.updatePSProfile(); err != nil {
				return err
			}
		}
	*/
	// TODO: Create autorun batch.
	// TODO: Write batch autorun path in registry.
	return nil
}

func (i *InstallerImpl) createPSScript() (err error) {
	// Initialize template.
	psTmpl, err := tmpl.ParseFS(template.Files, psTmpl)
	if err != nil {
		return errs.Wrap(ErrInstallerWindowsPSScriptTmplParse,
			errs.WithCause(err))
	}
	// Open file.
	psFile, err := os.OpenFile(i.psScriptPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return errs.Wrap(ErrInstallerWindowsPSScriptOpen, errs.WithCause(err),
			errs.WithContext("path", i.psScriptPath))
	}
	defer psFile.Close()
	// Write to file.
	if err = psTmpl.Execute(psFile, nil); err != nil {
		return errs.Wrap(ErrInstallerWindowsPSScirptTmplExecute,
			errs.WithCause(err))
	}
	return nil
}

func (i *InstallerImpl) getPSProfilePath() (string, error) {
	out, err := exec.Command("PowerShell.exe",
		"-Command", `$profile`).Output()
	if err != nil {
		err = errs.Wrap(ErrInstallerWindowsGetPSProfilePath,
			errs.WithCause(err))
		return "", err
	}
	// Sanitize output.
	path := strings.TrimSpace(string(out))
	Log().Info("PowerShell profile path",
		slog.String("path", path))
	return path, nil
}

/*
func (i *InstallerImpl) createPSProfile() (err error) {
	// Create directory.
	dirPath := filepath.Dir(i.psProfilePath)
	if err = os.MkdirAll(dirPath, 0o755); err != nil {
		return errs.Wrap(ErrInstallerWindowsCreatePSProfileDir,
			errs.WithCause(err), errs.WithContext("path", dirPath))
	}
	// Create profile.
	f, err := os.OpenFile(i.psProfilePath, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return errs.Wrap(ErrInstallerWindowsCreatePSProfile,
			errs.WithCause(err), errs.WithContext("path", i.psProfilePath))
	}
	defer f.Close()
	// Append lines to profile.
	fmt.Fprintln(f, `# Below lines are added by yenb.`)
	fmt.Fprintln(f, "&", i.psScriptPath)
	fmt.Fprintln(f, `# Above lines are added by yenb.`)
	Log().Info("Created new PS profile.",
		slog.String("path", i.psProfilePath))
	return nil
}

func (i *InstallerImpl) updatePSProfile() (err error) {
	f, err := os.OpenFile(i.psProfilePath, os.O_RDWR, 0o644)
	if err != nil {
		return errs.Wrap(ErrInstallerWindowsUpdatePSProfile,
			errs.WithCause(err), errs.WithContext("path", i.psProfilePath))
	}
	defer f.Close()
	// Scan lines and check if yenb lines exist.
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "yenb") {
			Log().Info("yenb is already installed in profile.")
			return nil
		}
	}
	if err = s.Err(); err != nil {
		return errs.Wrap(ErrInstallerWindowsScanPSProfile,
			errs.WithCause(err), errs.WithContext("path", i.psProfilePath))
	}
	// Go to end of file.
	if _, err = f.Seek(0, io.SeekEnd); err != nil {
		return errs.Wrap(ErrInstallerWindowsSeekPSProfile,
			errs.WithCause(err), errs.WithContext("path", i.psProfilePath))
	}
	// Append lines to PS profile.
	fmt.Fprintln(f, `# Below lines are added by yenb.`)
	fmt.Fprintln(f, "&", i.psScriptPath)
	fmt.Fprintln(f, `# Above lines are added by yenb.`)
	Log().Info("Updated PS profile.",
		slog.String("path", i.psProfilePath))
	return nil
}
*/
