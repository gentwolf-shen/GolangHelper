package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Out *out

type out struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

const (
	colorBlack = uint8(iota + 90)
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func init() {
	Out = &out{}

	flag := log.Ldate | log.Ltime | log.Llongfile

	Out.Info = log.New(os.Stdout, Info(""), flag)
	Out.Warn = log.New(os.Stdout, Warn(""), flag)

	dir := filepath.Dir(os.Args[0]) + "/log/"
	if _, err := os.Stat(dir); err != nil {
		os.Mkdir(dir, os.ModePerm)
	}

	filename := dir + "error-" + time.Now().Format("20060102") + ".txt"
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	Out.Error = log.New(io.MultiWriter(file, os.Stderr), Error(""), flag)
}

func Trace(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorCyan, "TRACE", msg)
}

func Error(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorRed, "ERROR", msg)
}

func Warn(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorYellow, "WARN", msg)
}

func Info(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorGreen, "INFO", msg)
}

func Debug(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorBlue, "DEBUG", msg)
}

func Assert(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorMagenta, "ASSERT", msg)
}
