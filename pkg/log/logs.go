package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level int

func LevelFromString(level string) Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return TRACE
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		panic(level + " is illegal log level. Legal levels: TRACE, DEBUG, INFO, WARN, ERROR, FATAL")
	}
}

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var logLevel Level

func Init(level Level, logsPath string) {
	file, err := os.OpenFile(logsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logLevel = level
}

func Trace(message ...any) {
	if logLevel <= TRACE {
		output(fmt.Sprintln("[TRACE]", message))
	}
}

func TraceF(message string, args ...any) {
	if logLevel <= TRACE {
		output(fmt.Sprintf("[TRACE] "+message+"\n", args))
	}
}

func Debug(message ...any) {
	if logLevel <= DEBUG {
		output(fmt.Sprintln("[DEBUG]", message))
	}
}

func DebugF(message string, args ...any) {
	if logLevel <= DEBUG {
		output(fmt.Sprintf("[DEBUG] "+message+"\n", args))
	}
}

func Info(message ...any) {
	if logLevel <= INFO {
		output(fmt.Sprintln("[INFO]", message))
	}
}

func InfoF(message string, args ...any) {
	if logLevel <= INFO {
		output(fmt.Sprintf("[INFO] "+message+"\n", args))
	}
}

func Warn(message ...any) {
	if logLevel <= WARN {
		output(fmt.Sprintln("[WARN]", message))
	}
}

func WarnF(message string, args ...any) {
	if logLevel <= WARN {
		output(fmt.Sprintf("[WARN] "+message+"\n", args))
	}
}

func Error(message ...any) {
	if logLevel <= ERROR {
		output(fmt.Sprintln("[ERROR]", message))
	}
}

func ErrorF(message string, args ...any) {
	if logLevel <= ERROR {
		output(fmt.Sprintf("[ERROR] "+message+"\n", args))
	}
}

func Panic(message ...any) {
	s := fmt.Sprintln(message)
	if logLevel <= ERROR {
		output("[ERROR] " + s)
	}
	panic(s)
}

func PanicF(message string, args ...any) {
	s := fmt.Sprintf(message+"\n", args)
	if logLevel <= ERROR {
		output("[ERROR] " + s)
	}
	panic(s)
}

func Fatal(message ...any) {
	s := fmt.Sprintln(message)
	output("[FATAL] " + s)
	os.Exit(1)
}

func FatalF(message string, args ...any) {
	s := fmt.Sprintf(message+"\n", args)
	output("[FATAL] " + s)
	os.Exit(1)
}

func output(output string) {
	err := log.Output(3, output)
	if err != nil {
		log.Println("!!!!!!!!!!!!!ERROR couldn't write log!!!!!!!!!!!!!\n", err)
		log.Println("Failed with", output)
	}
}
