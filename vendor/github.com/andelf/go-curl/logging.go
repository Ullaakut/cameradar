package curl

import (
    "log"
)

const (
    _DEBUG = 10 * (iota + 1)
    _INFO
    _WARN
    _ERROR
)

const _DEFAULT_LOG_LEVEL = _WARN

var log_level = _DEFAULT_LOG_LEVEL

// SetLogLevel changes the log level which determines the granularity of the
// messages that are logged.  Available log levels are: "DEBUG", "INFO",
// "WARN", "ERROR" and "DEFAULT_LOG_LEVEL".
func SetLogLevel(levelName string) {
    switch levelName {
    case "DEBUG":
        log_level = _DEBUG
    case "INFO":
        log_level = _INFO
    case "WARN":
        log_level = _WARN
    case "ERROR":
        log_level = _ERROR
    case "DEFAULT_LOG_LEVEL":
        log_level = _DEFAULT_LOG_LEVEL
    }
}

func logf(limitLevel int, format string, args ...interface{}) {
    if log_level <= limitLevel {
        log.Printf(format, args...)
    }
}

func debugf(format string, args ...interface{}) {
    logf(_DEBUG, format, args...)
}

func infof(format string, args ...interface{}) {
    logf(_INFO, format, args...)
}

func warnf(format string, args ...interface{}) {
    logf(_WARN, format, args...)
}

func errorf(format string, args ...interface{}) {
    logf(_ERROR, format, args...)
}
