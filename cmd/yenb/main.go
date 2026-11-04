/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"

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
	configJSON := json.RawMessage([]byte(config.String()))
	lib.Log().InfoContext(config.Context(), "yenb:start",
		slog.Any("config", configJSON))
	// Initialize env variable editor.
	_, err := lib.NewEnvEditor()
	if err != nil {
		lib.Log().Errorf(err.Error(), err)
		os.Exit(1)
	}
}
