package log

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func removeFile(testLogFile string) error {
	return os.Remove(testLogFile)
}

func writeLogToFile(testLogFile, content string) error {

	cfg := &LoggerConfig{false, "info", testLogFile}

	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}

	logger.Info(content)
	return nil
}

func checkLogContent(testLogFile, content string) error {

	fp, err := os.Open(testLogFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	bio := bufio.NewReader(fp)
	line, _, err := bio.ReadLine()
	if err != nil {
		return fmt.Errorf("fail to read log file, error: %s", err.Error())
	}

	if !strings.Contains(string(line), content) {
		return fmt.Errorf("test log content not found in file, origin: %s", string(line))
	}

	return nil
}

func TestNewLogger(t *testing.T) {

	cfg := &LoggerConfig{true, "info", ""}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if logger.cfg.Level.String() != cfg.Level {
		t.Fatalf("logger level not match, excepted: %s, result: %s ", cfg.Level, logger.cfg.Level.String())
	}
}

func TestNewLogger2(t *testing.T) {

	testLogFile := "unittest-logtest.log"
	defer removeFile(testLogFile)

	content := "test log"

	err := writeLogToFile(testLogFile, content)
	if err != nil {
		t.Fatal(err)
	}

	err = checkLogContent(testLogFile, content)
	if err != nil {
		t.Fatal(err)
	}

}

func TestLogger_SetLevel(t *testing.T) {

	cfg := &LoggerConfig{false, "info", ""}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Error(err.Error())
	}

	logger.SetLevel("debug")
	if logger.cfg.Level.String() != "debug" {
		t.Errorf("set log level fail, excepted: debug, result: %s", logger.cfg.Level.String())
	}
}

func TestLogger_Write(t *testing.T) {

	testLogFile := "unittest-logtest.log"
	defer removeFile(testLogFile)

	content := "test log"

	cfg := &LoggerConfig{false, "info", testLogFile}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logger.Write([]byte("test log string"))

	err = checkLogContent(testLogFile, content)
	if err != nil {
		t.Fatal(err)
	}
}
