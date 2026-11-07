/*
Copyright © 2026 srtkkou <srtkkou@outlook.jp>
*/
package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/goark/errs"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ロガー
type AppLogger struct {
	slog.Logger
	level *slog.LevelVar // ログレベル
}

const (
	// ログファイルのローテーションサイズ(MB)
	LOG_MAX_SIZE = 100
	// ログファイルの保存個数
	LOG_MAX_BACKUPS = 3
)

var (
	// ロガー初期化排他
	logOnce sync.Once
	// ロガーへのポインタ
	appLoggerPtr *AppLogger
)

var (
	ErrLogFileOpen  = errors.New("ErrLogFileOpen")
	ErrLogFileClose = errors.New("ErrLogFileClose")
)

func Log() *AppLogger {
	// ロガーの初期化
	logOnce.Do(func() {
		var writer io.Writer
		config := Config()
		// ログファイルの新規作成
		logPath := config.LogPath()
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			err = errs.Wrap(ErrLogFileOpen, errs.WithCause(err),
				errs.WithContext("path", logPath))
			log.Fatalf("%+v\n", err)
		}
		if err = f.Close(); err != nil {
			err = errs.Wrap(ErrLogFileClose, errs.WithCause(err),
				errs.WithContext("path", logPath))
			log.Fatalf("%+v\n", err)
		}
		// ログファイルのローテーション設定
		fileWriter := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    LOG_MAX_SIZE,
			MaxBackups: LOG_MAX_BACKUPS,
		}
		writer = io.MultiWriter(fileWriter, os.Stdout)
		// インスタンスの初期化
		appLogger := AppLogger{
			level: new(slog.LevelVar),
		}
		// レベルの設定
		if config.IsDebug() {
			appLogger.SetLevelDebug()
		} else {
			appLogger.SetLevelInfo()
		}
		// slogの初期化
		slogHandler := slog.NewJSONHandler(writer,
			&slog.HandlerOptions{Level: appLogger.level})
		slogLogger := slog.New(slogHandler)
		appLogger.Logger = *slogLogger
		appLoggerPtr = &appLogger
	})
	return appLoggerPtr
}

// DEBUGログの出力
func (l *AppLogger) DebugLog(msg string, args ...any) {
	config := Config()
	l.DebugContext(config.Context(), msg, args...)
}

// INFOログの出力
func (l *AppLogger) InfoLog(msg string, args ...any) {
	config := Config()
	l.InfoContext(config.Context(), msg, args...)
}

// ERRORログの出力
func (l *AppLogger) Errorf(msg string, err error) {
	config := Config()
	str := fmt.Sprintf("%+v", err)
	rawJSON := json.RawMessage([]byte(str))
	fmt.Fprintln(os.Stderr, str)
	l.ErrorContext(config.Context(), msg,
		slog.Any("error", rawJSON))
}

// ログレベルをDebug以上にする。
func (l *AppLogger) SetLevelDebug() {
	l.level.Set(slog.LevelDebug)
}

// ログレベルをInfo以上にする。
func (l *AppLogger) SetLevelInfo() {
	l.level.Set(slog.LevelInfo)
}

// ログレベルをWarn以上にする。
func (l *AppLogger) SetLevelWarn() {
	l.level.Set(slog.LevelWarn)
}
