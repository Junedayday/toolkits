package l4g

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestErrorLogPrint(t *testing.T) {
	defer os.Remove("testing.log")
	LoadConfiguration("../../configs/testing/testing.xml")
	Error("Error Level is ok")
	Info("Info Level is ok")
	Debug("Debug Level is ok")
	Fine("Fine Level is ok")
	time.Sleep(100 * time.Millisecond)

	b, err := ioutil.ReadFile("testing.log")
	if err != nil {
		t.Error("read log file failed")
		return
	}
	if strings.Index(string(b), "Error Level is ok") == -1 {
		t.Error("Print error log level failed!")
	}
}

func TestInfoLogPrint(t *testing.T) {
	defer os.Remove("testing.log")
	LoadConfiguration("../../configs/testing/testing.xml")
	Info("Info Level is ok")
	time.Sleep(100 * time.Millisecond)

	b, err := ioutil.ReadFile("testing.log")
	if err != nil {
		t.Error("read log file failed")
		return
	}
	if strings.Index(string(b), "Info Level is ok") == -1 {
		t.Error("Print info log level failed!")
	}
}

func TestDebugLogPrint(t *testing.T) {
	defer os.Remove("testing.log")
	LoadConfiguration("../../configs/testing/testing.xml")
	Debug("Debug Level is ok")
	time.Sleep(100 * time.Millisecond)

	b, err := ioutil.ReadFile("testing.log")
	if err != nil {
		t.Error("read log file failed")
		return
	}
	if strings.Index(string(b), "Debug Level is ok") == -1 {
		t.Error("Print debug log level failed!")
	}
}

func TestFineLogPrint(t *testing.T) {
	defer os.Remove("testing.log")
	LoadConfiguration("../../configs/testing/testing.xml")
	Fine("Fine Level is ok")
	time.Sleep(100 * time.Millisecond)

	b, err := ioutil.ReadFile("testing.log")
	if err != nil {
		t.Error("read log file failed")
		return
	}
	if strings.Index(string(b), "Fine Level is ok") == -1 {
		t.Error("Print fine log level failed!")
	}
}
