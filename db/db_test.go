package db

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNewRscsDBEmptyStr(t *testing.T) {
	_, err := NewRscsDB("", true)
	if err == nil {
		t.Errorf("empty file str")
	}
}

func TestNewRscsDBNotThere(t *testing.T) {
	_, err := NewRscsDB("thisisnotthere.sqlite3", true)
	if err == nil {
		t.Errorf("missing file")
	}
}

func TestNewRscsDBValid(t *testing.T) {
	_, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("valid file")
	}
}
