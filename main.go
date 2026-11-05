/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package main

import (
	"log"
	"os"

	"github.com/srtkkou/yenb/cmd/yenb"
	"github.com/srtkkou/yenb/lib"
)

// Version.
const Version = "0.0.1"

var (
	// Commit hash of VCS.
	CommitHash string
	// Commit date of VCS.
	CommitAt string
)

// Main
func main() {
	// Update config.
	config := lib.Config()
	if err := config.Load(); err != nil {
		log.Fatalf("%+v\n", err)
		os.Exit(1)
	}
	config.SetVersion(Version)
	config.SetCommitHash(CommitHash)
	config.SetCommitAt(CommitAt)
	// Execute main command.
	cmd := yenb.NewRootCmd()
	if err := cmd.Exec(os.Args); err != nil {
		lib.Log().Errorf("RootCmd.Exec", err)
		os.Exit(1)
	}
}
