// Package db provides an embeddable interface directly to the Rscs db.
// This is useful for client programs who want to access Rscs in isolation
// without the indirection of a network abstraction.
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
		t.Errorf("fail on empty file str")
	}
}

func TestNewRscsDBNotThere(t *testing.T) {
	_, err := NewRscsDB("thisisnotthere.sqlite3", true)
	if err == nil {
		t.Errorf("fail on missing file")
	}
}

func TestNewRscsDBValid(t *testing.T) {
	_, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("fail on valid file")
	}
}

func TestNewRscsDBValidReadOnly(t *testing.T) {
	rscsDB, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("fail on valid file")
	}
	if !rscsDB.ReadOnly() {
		t.Errorf("fail on readOnly")
	}
}

func TestNewRscsDBValidSHA(t *testing.T) {
	rscsDB, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("fail on valid file")
	}
	if rscsDB.SHA256() == "" {
		t.Errorf("fail on valid sha256")
	}
}

func TestNewRscsDBValidGet(t *testing.T) {
	rscsDB, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("fail on valid file")
	}
	value, found, getErr := rscsDB.Get("testkey")
	if getErr != nil {
		t.Errorf(getErr.Error())
	}
	if !found {
		t.Errorf("should have found expected key")
	}
	if value != "testvalue" {
		t.Errorf("return should be 'testvalue'")
	}
}

func TestNewRscsDBInvalidGet(t *testing.T) {
	rscsDB, err := NewRscsDB("../test/test.sqlite3", true)
	if err != nil {
		t.Errorf("fail on valid file")
	}
	_, found, getErr := rscsDB.Get("testkeyunexpected")
	if getErr != nil {
		t.Errorf(getErr.Error())
	}
	if found {
		t.Errorf("should have not found unexpected key")
	}
}
