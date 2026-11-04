/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goark/errs"
)

// Config for app.
type AppConfig struct {
	mtx        sync.RWMutex // Mutex
	exePath    string       // Path to executable.
	commitHash string       // Commit hash.
	commitAt   time.Time    // Commit date.
	isDebug    bool         // Debug mode flag.
}

const (
	// Commit date format.
	commitAtLayout = "20060102-150405"
	// Log file ext.
	logExt = "log"
)

var (
	// Execute only once.
	configOnce sync.Once
	// Pointer to AppConfig instance.
	configPtr *AppConfig
)

var (
	ErrConfigExePath = errors.New("ErrConfigExePath")
	ErrConfigAbsPath = errors.New("ErrConfigAbsPath")
)

// Get or initialize config.
func Config() *AppConfig {
	configOnce.Do(func() {
		configPtr = &AppConfig{
			commitHash: "",
			isDebug:    false,
		}
	})
	return configPtr
}

// Stringify AppConfig.
func (c *AppConfig) String() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	var sb strings.Builder
	sb.WriteString(`{"type":"AppConfig"`)
	sb.WriteString(`,"commitHash":`)
	sb.WriteString(strconv.Quote(c.commitHash))
	sb.WriteString(`,"commitAt":`)
	sb.WriteString(strconv.Quote(
		c.commitAt.Format(time.RFC3339)))
	sb.WriteString(`,"isDebug":`)
	sb.WriteString(strconv.FormatBool(c.isDebug))
	sb.WriteString(`}`)
	return sb.String()
}

// Update commit hash.
func (c *AppConfig) SetCommitHash(str string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.commitHash = str
}

// Get commit hash.
func (c *AppConfig) CommitHash() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.commitHash
}

// Update commit date.
func (c *AppConfig) SetCommitAt(at string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	tm, err := time.Parse(commitAtLayout, at)
	if err == nil {
		c.commitAt = tm
	}
}

// Get commit date.
func (c *AppConfig) CommitAt() time.Time {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.commitAt
}

// Update debug flag.
func (c *AppConfig) SetDebug(isDebug bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.isDebug = isDebug
}

// Get debug flag.
func (c *AppConfig) IsDebug() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.isDebug
}

// Get path to log file.
func (c *AppConfig) LogPath() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	ext := filepath.Ext(c.exePath)
	logPath := c.exePath[0:(len(c.exePath)-len(ext))] + "." + logExt
	return logPath
}

// Load settings.
func (c *AppConfig) Load() (err error) {
	if err = c.loadExePath(); err != nil {
		return err
	}
	return nil
}

// Load absolute path to executable.
func (c *AppConfig) loadExePath() (err error) {
	// 実行ファイルのパスの取得
	var exePath string
	exePath, err = os.Executable()
	if err != nil {
		return errs.Wrap(ErrConfigExePath, errs.WithCause(err))
	}
	// 絶対パスの取得
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return errs.Wrap(ErrConfigAbsPath, errs.WithCause(err))
	}
	// 実行パスの更新
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.exePath = absPath
	return nil
}
