// Package db provides an embeddable interface directly to the Rscs db.
// This is useful for client programs who want to access Rscs in isolation
// without the indirection of a network abstraction.
package db

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	memoryDBName = "file::memory:?mode=memory&cache=shared"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestWithTempDB(t *testing.T) {
	rscsDB, newErr := NewRscsDB(memoryDBName)
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

	t.Run("name", func(t *testing.T) {
		if rscsDB.DBFileName() != memoryDBName {
			t.Errorf("cannot name db")
		}
	})

	t.Run("insert", func(t *testing.T) {
		rowCount, insertErr := rscsDB.Insert("", testValue)
		if insertErr == nil {
			t.Errorf("insert on empty key")
		}
		if rowCount != 0 {
			t.Errorf("rowcount nonzero")
		}

		var badKey bytes.Buffer
		for i := 0; i <= 255; i++ {
			badKey.WriteString("a")
		}
		rowCount, insertErr = rscsDB.Insert(badKey.String(), testValue)
		if insertErr == nil {
			t.Errorf("insert on oversized key")
		}
		if rowCount != 0 {
			t.Errorf("rowcount nonzero")
		}

		rowCount, insertErr = rscsDB.Insert(testKey, testValue)
		if insertErr != nil {
			t.Errorf("insert fail:%s", insertErr.Error())
		}
		if rowCount != 1 {
			t.Errorf("insert rowcount:%d", rowCount)
		}

		var goodKey bytes.Buffer
		for i := 0; i < 255; i++ {
			goodKey.WriteString("a")
		}

		rowCount, insertErr = rscsDB.Insert(goodKey.String(), testValue)
		if insertErr != nil {
			t.Errorf("insert fail:%s", insertErr.Error())
		}
		if rowCount != 1 {
			t.Errorf("insert rowcount:%d", rowCount)
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

	t.Run("update-valid", func(t *testing.T) {
		rowCount, updateErr := rscsDB.Update("", "newval")
		if updateErr == nil {
			t.Errorf("update empty key")
		}
		if rowCount != 0 {
			t.Errorf("rowcount nonzero")
		}

		rowCount, updateErr = rscsDB.Update("notthere", "newval")
		if updateErr != nil {
			t.Errorf("update fail:%s", updateErr.Error())
		}
		if rowCount != 0 {
			t.Errorf("update rowcount:%d", rowCount)
		}

		rowCount, updateErr = rscsDB.Update(testKey, "newval")
		if updateErr != nil {
			t.Errorf("update fail:%s", updateErr.Error())
		}
		if rowCount != 1 {
			t.Errorf("update rowcount:%d", rowCount)
		}

		value, found, getErr := rscsDB.Get(testKey)
		if getErr != nil {
			t.Errorf("get fail:%s", getErr.Error())
		}
		if !found {
			t.Errorf("%s not found", testKey)
		}
		if value != "newval" {
			t.Errorf("not updated value")
		}
	})

	t.Run("delete", func(t *testing.T) {
		rowCount, deleteErr := rscsDB.Delete("")
		if deleteErr == nil {
			t.Errorf("delete empty key")
		}
		if rowCount != 0 {
			t.Errorf("rowcount nonzero")
		}

		rowCount, deleteErr = rscsDB.Delete("notthere")
		if deleteErr != nil {
			t.Errorf("delete missing key")
		}
		if rowCount != 0 {
			t.Errorf("rowcount nonzero")
		}

		rowCount, deleteErr = rscsDB.Delete(testKey)
		if deleteErr != nil {
			t.Errorf("delete fail:%s", deleteErr.Error())
		}
		if rowCount != 1 {
			t.Errorf("delete rowcount:%d", rowCount)
		}
	})

	t.Run("get-deleted", func(t *testing.T) {
		_, found, getErr := rscsDB.Get(testKey)
		if getErr != nil {
			t.Errorf("get fail:%s", getErr.Error())
		}
		if found {
			t.Errorf("deleted %s found", testKey)
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

	if rscsDB.DBFileName() != tmpDBFile.Name() {
		t.Errorf("cannot name db")
	}

	createErr := rscsDB.CreateTable()
	if createErr != nil {
		t.Fatalf("fail on create table:%s", createErr.Error())
	}

	const testKey = "testkey"
	const testValue = "testValue"

	rowCount, insertErr := rscsDB.Insert(testKey, testValue)
	if insertErr != nil {
		t.Errorf("insert fail:%s", insertErr.Error())
	}
	if rowCount != 1 {
		t.Errorf("insert rowcount:%d", rowCount)
	}

	// Now make it readonly.
	chmodErr := os.Chmod(tmpDBFile.Name(), 0444)
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
	_, insertReadOnlyErr := rscsDBReadOnly.Insert("testkey2", testValue)
	if insertReadOnlyErr == nil {
		t.Errorf("should not be able to insert into readonly file")
	}
	_, updateReadOnlyErr := rscsDBReadOnly.Update(testKey, "newval")
	if updateReadOnlyErr == nil {
		t.Errorf("should not be able to update into readonly file")
	}
	_, deleteReadOnlyErr := rscsDBReadOnly.Delete("testkey")
	if deleteReadOnlyErr == nil {
		t.Errorf("should not be able to delete from readonly file")
	}

	chmodErr = os.Chmod(tmpDBFile.Name(), 0755)
	if chmodErr != nil {
		t.Fatalf("chmod:%s", chmodErr.Error())
	}
}

func Example() {
	rscsDB, newErr := NewRscsDB("file::memory:?mode=memory&cache=shared")
	if newErr != nil {
		log.Fatalf("fail on tmpfile new:%s", newErr.Error())
	}
	createErr := rscsDB.CreateTable()
	if createErr != nil {
		log.Fatalf("fail on create table:%s", createErr.Error())
	}
	dupcreateErr := rscsDB.CreateTable()
	if dupcreateErr == nil {
		log.Fatalf("fail on duplicate create table:%s", dupcreateErr.Error())
	}
	dropErr := rscsDB.DropTable()
	if dropErr != nil {
		log.Fatalf("fail on drop table:%s", dropErr.Error())
	}
	recreateErr := rscsDB.CreateTable()
	if recreateErr != nil {
		log.Fatalf("fail on recreate table:%s", recreateErr.Error())
	}

	var rowCount int
	var insertErr, updateErr, deleteErr error

	testKey := "my-test-key"
	testValue := "my-test-value"

	rowCount, insertErr = rscsDB.Insert(testKey, testValue)
	if insertErr != nil {
		log.Fatalf("insert fail:%s", insertErr.Error())
	}
	if rowCount != 1 {
		log.Fatalf("insert rowcount:%d", rowCount)
	}

	value, found, getErr := rscsDB.Get(testKey)
	if getErr != nil {
		log.Fatalf("get fail:%s", getErr.Error())
	}
	if !found {
		log.Fatalf("%s not found", testKey)
	}
	if value != testValue {
		log.Fatalf("%s and %s not equal", value, testValue)
	}

	rowCount, updateErr = rscsDB.Update(testKey, "my-new-val")
	if updateErr != nil {
		log.Fatalf("update fail:%s", updateErr.Error())
	}
	if rowCount != 1 {
		log.Fatalf("update rowcount:%d", rowCount)
	}

	rowCount, deleteErr = rscsDB.Delete(testKey)
	if deleteErr != nil {
		log.Fatalf("delete fail:%s", deleteErr.Error())
	}
	if rowCount != 1 {
		log.Fatalf("delete rowcount:%d", rowCount)
	}
}
