package l4g

import (
	"github.com/alecthomas/log4go"
)

// LoadConfiguration from xml file
func LoadConfiguration(filename string) {
	log4go.LoadConfiguration(filename)
}

// Fine : print fine level log
func Fine(arg0 interface{}, args ...interface{}) {
	log4go.Fine(arg0, args...)
}

// Debug : print debug level log
func Debug(arg0 interface{}, args ...interface{}) {
	log4go.Debug(arg0, args...)
}

// Info : print info level log
func Info(arg0 interface{}, args ...interface{}) {
	log4go.Info(arg0, args...)
}

// Error : print error level log
func Error(arg0 interface{}, args ...interface{}) {
	log4go.Error(arg0, args...)
}
