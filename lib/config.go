/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/goark/errs"
	"github.com/goccy/go-yaml"
)

// Config for app.
type AppConfig struct {
	mtx        sync.RWMutex      // Mutex
	ctx        context.Context   // Context
	version    string            // Version.
	commitHash string            // Commit hash.
	commitAt   time.Time         // Commit date.
	isDebug    bool              // Debug mode flag.
	os         string            // Operating system.
	exePath    string            // Path to executable.
	userEnvs   map[string]envVar // Environment variables for user.
	userPaths  []string          // Paths to add to user PATH env.
}

// Structure of config file.
type configFile struct {
	User struct {
		Envs  map[string]envVar
		Paths []string
	}
}

// Environment variable.
type envVar struct {
	Type  string
	Value string
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
	ErrConfigExePath  = errors.New("ErrConfigExePath")
	ErrConfigAbsPath  = errors.New("ErrConfigAbsPath")
	ErrConfigReadYAML = errors.New("ErrConfigReadYAML")
	ErrConfigLoadYAML = errors.New("ErrConfigLoadYAML")
)

// Get or initialize config.
func Config() *AppConfig {
	configOnce.Do(func() {
		configPtr = &AppConfig{
			ctx:        context.Background(),
			commitHash: "",
			isDebug:    false,
			os:         "unknown",
		}
	})
	return configPtr
}

// Marshal to JSON.
func (c *AppConfig) MarshalJSON() ([]byte, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	var b bytes.Buffer
	fmt.Fprintf(&b, `{"commitHash":%q`, c.commitHash)
	fmt.Fprintf(&b, `,"commitAt":%q`,
		c.commitAt.Format(time.RFC3339))
	fmt.Fprintf(&b, `,"isDebug":%t`, c.isDebug)
	fmt.Fprintf(&b, `,"os":%q`, c.os)
	fmt.Fprintf(&b, `,"exePath":%q`, c.exePath)
	fmt.Fprintf(&b, `,"userEnvs":{`)
	if len(c.userEnvs) > 0 {
		for k, v := range c.userEnvs {
			fmt.Fprintf(&b, `%q:{"type":%q,"value":%q},`,
				k, v.Type, v.Value)
		}
		b.Truncate(b.Len() - 1)
	}
	fmt.Fprintf(&b, `},"userPaths":[`)
	if len(c.userPaths) > 0 {
		for _, path := range c.userPaths {
			fmt.Fprintf(&b, `%q,`, path)
		}
		b.Truncate(b.Len() - 1)
	}
	fmt.Fprintf(&b, `]}`)
	return b.Bytes(), nil
}

// Stringify AppConfig.
func (c *AppConfig) String() string {
	blob, _ := c.MarshalJSON()
	return string(blob)
}

// Get context.
func (c *AppConfig) Context() context.Context {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.ctx
}

// Update version.
func (c *AppConfig) SetVersion(str string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.version = str
}

// Get version.
func (c *AppConfig) Version() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.version
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

// Get path to exe file.
func (c *AppConfig) ExePath() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.exePath
}

// Get OS.
func (c *AppConfig) OS() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.os
}

// Get path to log file.
func (c *AppConfig) LogPath() string {
	ext := filepath.Ext(c.ExePath())
	logPath := c.exePath[0:(len(c.exePath)-len(ext))] + "." + logExt
	return logPath
}

// Update user envs.
func (c *AppConfig) SetUserEnvs(m map[string]envVar) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.userEnvs = m
}

// Get user envs.
func (c *AppConfig) UserEnvs() map[string]envVar {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.userEnvs
}

// Update user paths.
func (c *AppConfig) SetUserPaths(s []string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.userPaths = s
}

// Get user paths.
func (c *AppConfig) UserPaths() []string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.userPaths
}

// Get path to YAML file.
func (c *AppConfig) YAMLPath() string {
	exePath := c.ExePath()
	ext := filepath.Ext(exePath)
	path := c.exePath[0:(len(exePath)-len(ext))] + ".yml"
	return path
}

// Load settings.
func (c *AppConfig) Load() (err error) {
	c.loadRuntime()
	if err = c.loadExePath(); err != nil {
		return err
	}
	if err = c.loadYAML(); err != nil {
		return err
	}
	return nil
}

// Load runtime information.
func (c *AppConfig) loadRuntime() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.os = runtime.GOOS
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
	c.exePath = filepath.ToSlash(absPath)
	return nil
}

// Load config YAML file.
func (c *AppConfig) loadYAML() error {
	// 設定ファイルパスの読み込み
	ymlPath := c.YAMLPath()
	blob, err := os.ReadFile(ymlPath)
	if err != nil {
		return errs.Wrap(ErrConfigReadYAML, errs.WithCause(err),
			errs.WithContext("path", ymlPath))
	}
	// YAMLの読み込み
	values := new(configFile)
	if err := yaml.Unmarshal(blob, values); err != nil {
		return errs.Wrap(ErrConfigLoadYAML, errs.WithCause(err),
			errs.WithContext("path", ymlPath))
	}
	// 値の設定
	c.SetUserEnvs(values.User.Envs)
	c.SetUserPaths(values.User.Paths)
	return nil
}
