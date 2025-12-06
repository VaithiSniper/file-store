package logger

import (
	"fmt"
	"runtime"
	"time"
)

type LogLevel int

const (
	LevelEmergency LogLevel = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInfo
	LevelDebug
	LevelTrace
	LevelForce LogLevel = -1
)

var logLevelNames = map[LogLevel]string{
	LevelEmergency: "emergency",
	LevelAlert:     "alert",
	LevelCritical:  "critical",
	LevelError:     "error",
	LevelWarning:   "warning",
	LevelNotice:    "notice",
	LevelInfo:      "info",
	LevelDebug:     "debug",
	LevelTrace:     "trace",
	LevelForce:     "FORCE",
}

func (l LogLevel) String() string {
	return logLevelNames[l]
}

func GetLogLevelFromStr(levelStr string) LogLevel {
	for level, name := range logLevelNames {
		if name == levelStr {
			return level
		}
	}
	return DefaultLogLevel // Default
}

const DefaultLogLevel LogLevel = LevelNotice

var CurrentLogLevel = DefaultLogLevel // Hardcoded default

// formatLogMessage formats in the format:
// "ts"="<Timestamp>"|"l"="<Severity>"|"module":"p2p/store/api_server"|"m"="<Message>"
func formatLogMessage(
	severity LogLevel, message string, module string, args ...interface{},
) string {
	msg := fmt.Sprintf(message, args...)
	return fmt.Sprintf(
		`"ts"="%s"|"l"="%s"|"mod":"%s"|"m"="%s"`,
		time.Now().Format("2006-01-02 15:04:05.000"),
		severity, module, msg,
	)
}

func log(
	severity LogLevel, message string, module string, args ...interface{},
) {
	if severity > CurrentLogLevel {
		return
	}
	fmt.Println(formatLogMessage(severity, message, module, args...))
}

func LogEmergency(module string, message string, args ...interface{}) {
	log(LevelEmergency, message, module, args...)
	// TODO: Maybe let's shut down the program here?
}

func LogAlert(module string, message string, args ...interface{}) {
	log(LevelAlert, message, module, args...)
}

func LogCritical(module string, message string, args ...interface{}) {
	log(LevelCritical, message, module, args...)
}

func LogError(module string, message string, args ...interface{}) {
	log(LevelError, message, module, args...)
}

func LogWarning(module string, message string, args ...interface{}) {
	log(LevelWarning, message, module, args...)
}

func LogNotice(module string, message string, args ...interface{}) {
	log(LevelNotice, message, module, args...)
}

func LogInfo(module string, message string, args ...interface{}) {
	log(LevelInfo, message, module, args...)
}

func LogDebug(module string, message string, args ...interface{}) {
	log(LevelDebug, message, module, args...)
}

func LogTrace(module string, message string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1) // use 2 if you wrap LogTrace further
	callerFuncName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			callerFuncName = fn.Name()
		}
	}
	fmtMsg := fmt.Sprintf(
		`%s|"line":"%d"|"function":"%s"`,
		formatLogMessage(LevelTrace, message, module, args...), line,
		callerFuncName,
	)
	fmt.Println(fmtMsg)
}

func LogForce(module string, message string, args ...interface{}) {
	log(LevelForce, message, module, args...)
}
