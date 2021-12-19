package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// HilcoLogger with zeroLog
type HilcoLogger struct {
	logger *zerolog.Logger
}

//NewMultiWithFile with file
func NewMultiWithFile(isDebug bool) *HilcoLogger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	//file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error: HilcoLogger file error")
	}
	//multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout, file)
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	return &HilcoLogger{logger: &logger}
}

// New for zerolog
func New(isDebug bool) *HilcoLogger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &HilcoLogger{logger: &logger}
}

//NewConsole with Concole
func NewConsole(isDebug bool) *HilcoLogger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &HilcoLogger{logger: &logger}
}

// Output duplicates the global logger and sets w as its output.
func (l *HilcoLogger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

// Info is a standard error message
func (l *HilcoLogger) Info(message string) {
	l.logger.Info().Msg(message)
}

//Debug is a standard error message
func (l *HilcoLogger) Debug(message string) {
	l.logger.Debug().Msg(message)
}

//Fatal is a standard error message
func (l *HilcoLogger) Fatal(message string) {
	l.logger.Fatal().Msg(message)
}

//Erroris a standard error message
func (l *HilcoLogger) Error(message string) {
	l.logger.Error().Msg(message)
}
