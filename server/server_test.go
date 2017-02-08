package server

import (
	"github.com/bradclawsie/rscs/db"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestVaidRscsDB(t *testing.T) {
	tmpDBFile, tmpFileErr := ioutil.TempFile("", "db_test_tmp")
	if tmpFileErr != nil {
		t.Fatalf(tmpFileErr.Error())
	}
	defer os.Remove(tmpDBFile.Name())
	rscsDB, newDBErr := db.NewRscsDB(tmpDBFile.Name())
	if newDBErr != nil {
		t.Fatalf("fail on tmpfile new:%s", newDBErr.Error())
	}
	createErr := rscsDB.CreateTable()
	if createErr != nil {
		t.Fatalf("fail on create table:%s", createErr.Error())
	}
	_, newServerErr := NewRscsServer(rscsDB)
	if newServerErr != nil {
		t.Errorf("fail on valid RscsDB:%s", newServerErr.Error())
	}
}

func TestNilRscsDB(t *testing.T) {
	_, newErr := NewRscsServer(nil)
	if newErr == nil {
		t.Errorf("fail on nil RscsDB")
	}
}
