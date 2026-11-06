/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"log/slog"
	"regexp"
	"strings"
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

func GetPSExecutionPolicy() (PSExecutionPolicy, error) {
	stdout, err := invoke(powerShell, "-Command",
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
	Log().Info("PowerShell execution policy of current user.",
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
