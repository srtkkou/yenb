/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package main

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/srtkkou/yenb/lib"
)

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
	}
	config.SetCommitHash(CommitHash)
	config.SetCommitAt(CommitAt)
	// Initialize log.
	lib.Log().Info("yenb:start",
		slog.Any("config", json.RawMessage([]byte(config.String()))))
}
