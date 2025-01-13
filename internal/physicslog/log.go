// Package physicslog is a simple utility package that offers logging
// for our physics engine.
package physicslog

import (
	"log"
	"os"
	"runtime/debug"
	"time"
)

// Level represents the severity level for a log entry.
type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
	LevelOff
)

// String returns a human-readable string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelOff:
		return "OFF"
	default:
		return ""
	}
}

// PhysicsLogger is a type that implements a basic logger.
type PhysicsLogger struct {
	logger   *log.Logger // Logger is guaranteed to be serial.
	minLevel Level       // The minimum severity level log entries are written for
}

// NewPhysicsLogger creates a PhysicsLogger object with a specified logging level.
// It writes to os.Stdout by default.
func NewPhysicsLogger(level Level) *PhysicsLogger {
	return &PhysicsLogger{
		logger:   log.New(os.Stdout, "", 0),
		minLevel: level,
	}
}

// LogInfo logs a message at INFO level.
func (p *PhysicsLogger) LogInfo(message string) {
	p.log(LevelInfo, message)
}

// LogError logs a message at ERROR level.
func (p *PhysicsLogger) LogError(message string) {
	p.log(LevelError, message)
}

// LogFatal logs a message at FATAL level. It also terminates the goroutine it
// was called on with os.Exit.
func (p *PhysicsLogger) LogFatal(message string) {
	p.log(LevelFatal, message)
	os.Exit(1)
}

// log formats and writes a log entry with the specified message and log entry.
func (p *PhysicsLogger) log(level Level, message string) {
	if level < p.minLevel || level == LevelOff {
		return
	}

	trace := ""

	if level >= LevelError {
		trace = string(debug.Stack())
	}

	t := time.Now().UTC().Format(time.RFC3339)
	p.logger.Printf("[%s %s] %s %s", level.String(), t, message, trace)
}
