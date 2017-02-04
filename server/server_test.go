package server

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNewRscsServerEmptyStr(t *testing.T) {
	_, err := NewRscsServer("", true)
	if err == nil {
		t.Errorf("empty file str")
	}
}

func TestNewRscsServerNotThere(t *testing.T) {
	_, err := NewRscsServer("thisisnotthere.sqlite3", true)
	if err == nil {
		t.Errorf("missing file")
	}
}

func TestNewRscsServerValid(t *testing.T) {
	_, err := NewRscsServer("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("valid file")
	}
}
