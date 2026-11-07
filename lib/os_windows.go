/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/goark/errs"
	reg "golang.org/x/sys/windows/registry"
)

type PSExecutionPolicy string

const (
	PolicyAllSigned    = PSExecutionPolicy("AllSigned")
	PolicyBypass       = PSExecutionPolicy("Bypass")
	PolicyDefault      = PSExecutionPolicy("Default")
	PolicyRemoteSigned = PSExecutionPolicy("RemoteSigned")
	PolicyRestricted   = PSExecutionPolicy("Restricted")
	PolicyUndefined    = PSExecutionPolicy("Undefined")
	PolicyUnrestricted = PSExecutionPolicy("Unrestricted")
)

type PSExecutionPolicyScope string

const (
	ScopeMachinePolicy = PSExecutionPolicyScope("MachinePolicy")
	ScopeUserPolicy    = PSExecutionPolicyScope("UserPolicy")
	ScopeProcess       = PSExecutionPolicyScope("Process")
	ScopeCurrentUser   = PSExecutionPolicyScope("CurrentUser")
	ScopeLocalMachine  = PSExecutionPolicyScope("LocalMachine")
)

const (
	powerShell = "PowerShell.exe"
)

var (
	rePSExecutionPolicy = regexp.MustCompile(`CurrentUser\s+([a-zA-Z]+)`)
)

var (
	ErrOSWindowsOpenRegistryKey     = errors.New("ErrOSWindowsOpenRegistryKey")
	ErrOSWindowsSetExpandSZ = errors.New("ErrOSWindowsSetExpandSZ")
)

func SetRegExpandSZ(key reg.Key, path, name, value string) error {
	// Open registry path.
	regKey, err := reg.OpenKey(key, path, reg.SET_VALUE)
	if err != nil {
		return errs.Wrap(ErrOSWindowsOpenRegistryKey, errs.WithCause(err),
			errs.WithContext("key", key),
			errs.WithContext("path", path),
			errs.WithContext("name", name),
			errs.WithContext("value", value))
	}
	defer regKey.Close()
	// Set REG_EXPAND_SZ value.
	err = regKey.SetExpandStringValue(name, value)
	if err != nil {
		return errs.Wrap(ErrOSWindowsSetExpandSZ, errs.WithCause(err),
			errs.WithContext("key", key),
			errs.WithContext("path", path),
			errs.WithContext("name", name),
			errs.WithContext("value", value))
	}
	Log().InfoLog("SetRegExpandSZ",
		slog.Any("key", key),
		slog.String("path", path),
		slog.String("name", name),
		slog.String("value", value))
	return nil
}

func GetPSExecutionPolicy() (PSExecutionPolicy, error) {
	stdout, err := Invoke(powerShell, "-Command",
		"Get-ExecutionPolicy", "-List")
	if err != nil {
		return PolicyUndefined, err
	}
	// Read STDOUT lines.
	policyStr := ""
	for line := range strings.Lines(stdout) {
		matches := rePSExecutionPolicy.FindStringSubmatch(line)
		if len(matches) == 2 {
			policyStr = matches[1]
			continue
		}
	}
	Log().InfoLog("PowerShell execution policy of current user.",
		slog.String("policy", policyStr))
	return ToPSExecutionPolicy(policyStr), nil
}

func ToPSExecutionPolicy(str string) PSExecutionPolicy {
	policy := PSExecutionPolicy(str)
	switch policy {
	case PolicyAllSigned, PolicyBypass, PolicyDefault, PolicyRemoteSigned, PolicyRestricted, PolicyUndefined, PolicyUnrestricted:
		return policy
	default:
		return PolicyUndefined
	}
}

func ToPSExecutionPolicyScope(str string) PSExecutionPolicyScope {
	scope := PSExecutionPolicyScope(str)
	switch scope {
	case ScopeMachinePolicy, ScopeUserPolicy, ScopeProcess, ScopeCurrentUser, ScopeLocalMachine:
		return scope
	default:
		return ScopeProcess
	}
}
