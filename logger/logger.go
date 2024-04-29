/*
Copyright © 2023 vavuthu@redhat.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"	
)

var Logger *zap.Logger

func init()  {
	// construct filename with timestamp
	currentTime := time.Now()
	timestamp := currentTime.Format("2019-02-09_17-50-01")
	filename := "itr_" + timestamp + ".log"

	// create encoderconfig with timeformat
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.StacktraceKey = ""

	// create console and file enconder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	
	// create log file
	logFile, err := os.Create(filename)
	if err != nil {
		fmt.Errorf("failed to open log file: %v", err)
	}

	// create writers for console and file
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileWriter := zapcore.AddSync(logFile)
	
	defaultLogLevel := zapcore.DebugLevel

	// create cores for writing to console and file
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, defaultLogLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel)

	// create core with both console and file
	core := zapcore.NewTee(consoleCore, fileCore)

	// create logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).WithOptions(zap.AddCallerSkip(1))
	
	Logger.Info("log file: " + filename)
}

// concatenateMsg concatenates message and values
func concatenateMsg(msg string, values ...interface{}) string {
    if len(values) == 0 {
        return msg
    }

    var concatenatedMsg string
    for _, v := range values {
        concatenatedMsg += fmt.Sprintf(" %v", v)
    }
    return msg + concatenatedMsg
}

// Info logs a message at Info level with optional variable value
func Info(msg string, values ...interface{}) {
    Logger.Info(concatenateMsg(msg, values...))
}

// Infof logs a formatted message at Info level
func Infof(format string, values ...interface{}) {
    msg := fmt.Sprintf(format, values...)
    Logger.Info(msg)
}

// Warn logs a message at Warn level with optional variable value
func Warn(msg string, values ...interface{}) {
    Logger.Warn(concatenateMsg(msg, values...))
}

// Warnf logs a formatted message at Warn level
func Warnf(format string, values ...interface{}) {
    msg := fmt.Sprintf(format, values...)
    Logger.Warn(msg)
}

// Error logs a message at Error level with optional variable value
func Error(msg string, values ...interface{}) {
    Logger.Error(concatenateMsg(msg, values...))
}

// Errorf logs a formatted message at Error level
func Errorf(format string, values ...interface{}) {
    msg := fmt.Sprintf(format, values...)
    Logger.Error(msg)
}

// Debug logs a message at Debug level with optional variable value
func Debug(msg string, values ...interface{}) {
    Logger.Debug(concatenateMsg(msg, values...))
}

// Debugf logs a formatted message at Debug level
func Debugf(format string, values ...interface{}) {
    msg := fmt.Sprintf(format, values...)
    Logger.Debug(msg)
}
