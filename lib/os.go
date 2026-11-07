/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	"github.com/goark/errs"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	ErrOSMkdirAll         = errors.New("ErrOSMkdirAll")
	ErrOSRemoveAll        = errors.New("ErrOSRemoveAll")
	ErrOSInvokePipeStdout = errors.New("ErrOSInvokePipeStdout")
	ErrOSInvokePipeStderr = errors.New("ErrOSInvokePipeStderr")
	ErrOSInvokeCmdStart   = errors.New("ErrOSInvodeCmdStart")
	ErrOSInvokeReadStdout = errors.New("ErrOSInvokeReadStdout")
	ErrOSInvokeReadStderr = errors.New("ErrOSInvokeReadStderr")
	ErrOSInvokeCmdWait    = errors.New("ErrOSInvokeCmdWait")
	ErrOSInvokeCmdErr     = errors.New("ErrOSInvokeCmdErr")
	ErrOSUTF8Convert      = errors.New("ErrOSUTF8Convert")
)

func MkdirAll(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0o755); err != nil {
			return errs.Wrap(ErrOSMkdirAll, errs.WithCause(err),
				errs.WithContext("path", path))
		}
		Log().InfoLog("Created directory.", slog.String("path", path))
	}
	return nil
}

func RemoveAll(path string) (err error) {
	if _, err = os.Stat(path); os.IsExist(err) {
		if err = os.RemoveAll(path); err != nil {
			return errs.Wrap(ErrOSRemoveAll, errs.WithCause(err),
				errs.WithContext("path", path))
		}
		Log().InfoLog("Deleted directory.", slog.String("path", path))
	}
	return nil
}

func Invoke(command string, args ...string) (string, error) {
	// Initialize command.
	config := Config()
	cmd := exec.CommandContext(config.Context(), command, args...)
	commandStr := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	// Connect STDOUT pipe.
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		err = errs.Wrap(ErrOSInvokePipeStdout, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	defer stdoutPipe.Close()
	// Connect STDERR pipe.
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		err = errs.Wrap(ErrOSInvokePipeStderr, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	defer stderrPipe.Close()
	// Start command.
	if err = cmd.Start(); err != nil {
		err = errs.Wrap(ErrOSInvokeCmdStart, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	// Read STDOUT.
	stdoutBlob, err := io.ReadAll(stdoutPipe)
	if err != nil {
		err = errs.Wrap(ErrOSInvokeReadStdout, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	// Read STDERR.
	stderrBlob, err := io.ReadAll(stderrPipe)
	if err != nil {
		err = errs.Wrap(ErrOSInvokeReadStderr, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	// Wait for command to finish.
	if err = cmd.Wait(); err != nil {
		err = errs.Wrap(ErrOSInvokeCmdWait, errs.WithCause(err),
			errs.WithContext("command", commandStr))
		return "", err
	}
	// Format STDOUT.
	stdoutUTF8Blob, err := toUTF8(stdoutBlob)
	if err != nil {
		return "", err
	}
	stdout := string(bytes.TrimSpace(stdoutUTF8Blob))
	// Format STDERR.
	stderrUTF8Blob, err := toUTF8(stderrBlob)
	if err != nil {
		return "", err
	}
	stderr := string(bytes.TrimSpace(stderrUTF8Blob))
	// Check if STDERR is empty.
	if len(stderr) > 0 {
		err = errs.Wrap(ErrOSInvokeCmdErr,
			errs.WithContext("command", commandStr),
			errs.WithContext("stderr", stderr))
		return "", err
	}
	Log().InfoLog("Invoke", slog.String("command", commandStr),
		slog.String("stdout", stdout), slog.String("stderr", stderr))
	return stdout, nil
}

func toUTF8(blob []byte) ([]byte, error) {
	// Do nothing if bytes are valid UTF-8.
	if !utf8.Valid(blob) {
		return blob, nil
	}
	// Try unicode conversion.
	fallback := transform.Nop
	r := transform.NewReader(bytes.NewReader(blob),
		unicode.BOMOverride(fallback))
	u8Blob, err := io.ReadAll(r)
	if err != nil {
		err := errs.Wrap(ErrOSUTF8Convert,
			errs.WithCause(err),
			errs.WithContext("blob", blob))
		return []byte{}, err
	}
	return u8Blob, nil
}
