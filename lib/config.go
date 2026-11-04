/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/goark/errs"
	"github.com/pelletier/go-toml/v2"
)

// Config for app.
type AppConfig struct {
	mtx        sync.RWMutex      // Mutex
	ctx        context.Context   // Context
	commitHash string            // Commit hash.
	commitAt   time.Time         // Commit date.
	isDebug    bool              // Debug mode flag.
	exePath    string            // Path to executable.
	userEnvs   map[string]string // Environment variables for user.
	userPaths  []string          // Paths to add to user PATH env.
}

type tomlStructure struct {
	User struct {
		Envs  map[string]string `toml:"env"`  // Environment variables.
		Paths []string          `toml:"path"` // Paths
	} `toml:"user"`
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
	ErrConfigReadTOML = errors.New("ErrConfigReadTOML")
	ErrConfigLoadTOML = errors.New("ErrConfigLoadTOML")
)

// Get or initialize config.
func Config() *AppConfig {
	configOnce.Do(func() {
		configPtr = &AppConfig{
			ctx:        context.Background(),
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
	var b bytes.Buffer
	b.WriteString(`{"type":"AppConfig"`)
	b.WriteString(`,"commitHash":`)
	b.WriteString(strconv.Quote(c.commitHash))
	b.WriteString(`,"commitAt":`)
	b.WriteString(strconv.Quote(
		c.commitAt.Format(time.RFC3339)))
	b.WriteString(`,"isDebug":`)
	b.WriteString(strconv.FormatBool(c.isDebug))
	b.WriteString(`,"exePath":`)
	b.WriteString(strconv.Quote(c.exePath))
	b.WriteString(`,"userEnvs":{`)
	if len(c.userEnvs) > 0 {
		for key, value := range c.userEnvs {
			b.WriteString(strconv.Quote(key))
			b.WriteString(`:`)
			b.WriteString(strconv.Quote(value))
			b.WriteString(`,`)
		}
		b.Truncate(b.Len() - 1)
	}
	b.WriteString(`},"userPaths":[`)
	if len(c.userPaths) > 0 {
		for _, path := range c.userPaths {
			b.WriteString(strconv.Quote(path))
			b.WriteString(`,`)
		}
		b.Truncate(b.Len() - 1)
	}
	b.WriteString(`]}`)
	return b.String()
}

// Get context.
func (c *AppConfig) Context() context.Context {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.ctx
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

// Get path to log file.
func (c *AppConfig) LogPath() string {
	ext := filepath.Ext(c.ExePath())
	logPath := c.exePath[0:(len(c.exePath)-len(ext))] + "." + logExt
	return logPath
}

// Update user envs.
func (c *AppConfig) SetUserEnvs(m map[string]string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.userEnvs = m
}

// Get user envs.
func (c *AppConfig) UserEnvs() map[string]string {
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

// Get path to TOML file.
func (c *AppConfig) TOMLPath() string {
	ext := filepath.Ext(c.ExePath())
	tomlPath := c.exePath[0:(len(c.exePath)-len(ext))] + ".toml"
	return tomlPath
}

// Load settings.
func (c *AppConfig) Load() (err error) {
	if err = c.loadExePath(); err != nil {
		return err
	}
	if err = c.loadTOML(); err != nil {
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
	c.exePath = filepath.ToSlash(absPath)
	return nil
}

// 設定tomlファイルの読み込み
func (c *AppConfig) loadTOML() error {
	// 設定ファイルパスの読み込み
	tomlPath := c.TOMLPath()
	blob, err := os.ReadFile(tomlPath)
	if err != nil {
		return errs.Wrap(ErrConfigReadTOML, errs.WithCause(err),
			errs.WithContext("path", tomlPath))
	}
	// TOMLの読み込み
	values := new(tomlStructure)
	if err := toml.Unmarshal(blob, values); err != nil {
		return errs.Wrap(ErrConfigLoadTOML, errs.WithCause(err),
			errs.WithContext("path", tomlPath))
	}
	// 値の設定
	c.SetUserEnvs((values.User.Envs))
	c.SetUserPaths((values.User.Paths))
	return nil
}
