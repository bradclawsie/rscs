// Package db provides an embeddable interface directly to the Rscs db.
// This is useful for client programs who want to access Rscs in isolation
// without the indirection of a network abstraction.
package db

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestWithTempDB(t *testing.T) {
	tmpDBFile, tmpFileErr := ioutil.TempFile("", "db_test_tmp")
	if tmpFileErr != nil {
		t.Fatalf(tmpFileErr.Error())
	}
	defer os.Remove(tmpDBFile.Name())
	rscsDB, newErr := NewRscsDB(tmpDBFile.Name())
	if newErr != nil {
		t.Fatalf("fail on tmpfile new:%s", newErr.Error())
	}
	createErr := rscsDB.CreateTable()
	if createErr != nil {
		t.Fatalf("fail on create table:%s", createErr.Error())
	}
	dupcreateErr := rscsDB.CreateTable()
	if dupcreateErr == nil {
		t.Fatalf("fail on duplicate create table:%s", dupcreateErr.Error())
	}
	dropErr := rscsDB.DropTable()
	if dropErr != nil {
		t.Fatalf("fail on drop table:%s", dropErr.Error())
	}
	recreateErr := rscsDB.CreateTable()
	if recreateErr != nil {
		t.Fatalf("fail on recreate table:%s", recreateErr.Error())
	}

	const testKey = "testkey"
	const testValue = "testValue"

	t.Run("insert", func(t *testing.T) {
		insertErr := rscsDB.Insert(testKey, testValue)
		if insertErr != nil {
			t.Errorf("insert fail:%s", insertErr.Error())
		}
	})

	t.Run("get-valid", func(t *testing.T) {
		value, found, getErr := rscsDB.Get(testKey)
		if getErr != nil {
			t.Errorf("get fail:%s", getErr.Error())
		}
		if !found {
			t.Errorf("%s not found", testKey)
		}
		if value != testValue {
			t.Errorf("%s and %s not equal", value, testValue)
		}
	})

	t.Run("get-invalid", func(t *testing.T) {
		badKey := "some-bad-key"
		_, found, getErr := rscsDB.Get(badKey)
		if getErr != nil {
			t.Errorf("get fail:%s", getErr.Error())
		}
		if found {
			t.Errorf("%s found?", badKey)
		}
	})
}

func TestWithReadonlyTempDB(t *testing.T) {
	tmpDBFile, tmpFileErr := ioutil.TempFile("", "db_test_tmp_rdonly")
	if tmpFileErr != nil {
		t.Fatalf(tmpFileErr.Error())
	}
	defer os.Remove(tmpDBFile.Name())
	rscsDB, newErr := NewRscsDB(tmpDBFile.Name())
	if newErr != nil {
		t.Fatalf("fail on tmpfile new:%s", newErr.Error())
	}
	createErr := rscsDB.CreateTable()
	if createErr != nil {
		t.Fatalf("fail on create table:%s", createErr.Error())
	}

	const testKey = "testkey"
	const testValue = "testValue"

	insertErr := rscsDB.Insert(testKey, testValue)
	if insertErr != nil {
		t.Fatalf("insert fail:%s", insertErr.Error())
	}

	// Now make it readonly.
	chmodErr := os.Chmod(tmpDBFile.Name(),0444)
	if chmodErr != nil {
		t.Fatalf("chmod:%s", chmodErr.Error())
	}

	rscsDBReadOnly, newErrReadOnly := NewRscsDB(tmpDBFile.Name())
	if newErrReadOnly != nil {
		t.Fatalf("fail on readonly tmpfile new:%s", newErrReadOnly.Error())
	}

	// Should be able to read from read-only file.
	value, found, getErr := rscsDBReadOnly.Get(testKey)
	if getErr != nil {
		t.Errorf("get fail:%s", getErr.Error())
	}
	if !found {
		t.Errorf("%s not found", testKey)
	}
	if value != testValue {
		t.Errorf("%s and %s not equal", value, testValue)
	}

	// We should not be able to write to the read-only file.
	insertReadOnlyErr := rscsDBReadOnly.Insert("testkey2", testValue)
	if insertReadOnlyErr == nil {
		t.Errorf("should not be able to insert into readonly file")
	}
	
	chmodErr = os.Chmod(tmpDBFile.Name(),0755)
	if chmodErr != nil {
		t.Fatalf("chmod:%s", chmodErr.Error())
	}
}

func TestNewRscsDBEmptyStr(t *testing.T) {
	_, err := NewRscsDB("")
	if err == nil {
		t.Errorf("fail on empty file str")
	}
}

func TestNewRscsDBNotThere(t *testing.T) {
	_, err := NewRscsDB("thisisnotthere.sqlite3")
	if err == nil {
		t.Errorf("fail on missing file")
	}
}
