package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"stellarisWikiCrawler/internal/config"
)

// ---------- level ----------

type Level int

const (
	LevelError Level = iota + 1
	LevelSuc
	LevelWarn
	LevelInfo
	LevelDebug
)

// ---------- color ----------

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

// ---------- prefix ----------

const (
	prefixError = "[error] "
	prefixSuc   = "[suc] "
	prefixWarn  = "[warn] "
	prefixInfo  = "[info] "
	prefixDebug = "[debug] "
)

// ---------- logger ----------

type Logger struct {
	l        *log.Logger
	level    Level
	useColor bool

	task string

	mu sync.Mutex
}

// parseLogLevel converts a string log level to the Level type.
func parseLogLevel(levelStr string) Level {
	switch strings.ToLower(levelStr) {
	case "error":
		return LevelError
	case "suc":
		return LevelSuc
	case "warn":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		return LevelInfo // Default to LevelInfo if the input is invalid
	}
}

// ---------- constructor ----------

func NewLogger() *Logger {
	cfg := config.Get()
	return &Logger{
		l:        log.New(os.Stdout, "", log.LstdFlags),
		level:    parseLogLevel(cfg.LogLevel),
		useColor: isTTY(),
	}
}

// 派生一个带 task 的 logger
func (lg *Logger) WithTask(task string) *Logger {
	return &Logger{
		l:        lg.l,
		level:    lg.level,
		useColor: lg.useColor,
		task:     task,
	}
}

// ---------- public api ----------

func (lg *Logger) Error(v ...any) {
	if lg.level >= LevelError {
		lg.print(colorRed, prefixError, v...)
	}
}

func (lg *Logger) Suc(v ...any) {
	if lg.level >= LevelSuc {
		lg.print(colorGreen, prefixSuc, v...)
	}
}

func (lg *Logger) Warn(v ...any) {
	if lg.level >= LevelWarn {
		lg.print(colorYellow, prefixWarn, v...)
	}
}

func (lg *Logger) Info(v ...any) {
	if lg.level >= LevelInfo {
		lg.print(colorBlue, prefixInfo, v...)
	}
}

func (lg *Logger) Debug(v ...any) {
	if lg.level >= LevelDebug {
		lg.print(colorGray, prefixDebug, v...)
	}
}

// ---------- internal ----------

func (lg *Logger) print(color, prefix string, v ...any) {
	msg, fields := splitMessageAndFields(v...)

	lg.mu.Lock()
	defer lg.mu.Unlock()

	prefixTask := ""
	if lg.task != "" {
		prefixTask = fmt.Sprintf("[%s] ", lg.task)
	}

	if lg.useColor {
		lg.l.SetPrefix(color + prefix)
	} else {
		lg.l.SetPrefix(prefix)
	}

	if fields == "" {
		lg.l.Println(prefixTask + msg)
	} else {
		lg.l.Println(prefixTask + msg + " | " + fields)
	}

	if lg.useColor {
		lg.l.SetPrefix(colorReset)
	} else {
		lg.l.SetPrefix("")
	}
}

func splitMessageAndFields(v ...any) (string, string) {
	if len(v) == 0 {
		return "", ""
	}

	msg := fmt.Sprint(v[0])

	if len(v) < 3 {
		return msg, ""
	}

	var b strings.Builder
	for i := 1; i+1 < len(v); i += 2 {
		key := fmt.Sprint(v[i])
		val := fmt.Sprint(v[i+1])
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(val)
		b.WriteString(" ")
	}

	return msg, strings.TrimSpace(b.String())
}

// ---------- tty ----------

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
