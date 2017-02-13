package server

import (
	"encoding/json"
	"github.com/bradclawsie/rscs/db"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	memoryDBName = "file::memory:?mode=memory&cache=shared"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {

	rscsDB, rscsDBErr := db.NewRscsDB(memoryDBName)
	if rscsDBErr != nil {
		log.Fatal(rscsDBErr)
	}
	rscsServer, rscsSrvErr := NewRscsServer(rscsDB)
	if rscsSrvErr != nil {
		log.Fatal(rscsSrvErr)
	}

	var rtrErr error
	rtr, rtrErr := rscsServer.NewRouter()
	if rtrErr != nil {
		log.Fatal(rtrErr.Error())
	}

	testServer = httptest.NewServer(rtr)
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestStatus(t *testing.T) {
	resp, statusJSON := testRequest(t, testServer, http.MethodGet, StatusRoute, nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("not 200")
	}
	var status StatusResult
	umErr := json.Unmarshal([]byte(statusJSON), &status)
	if umErr != nil {
		t.Errorf(umErr.Error())
	}
}

func TestVaidRscsDB(t *testing.T) {
	rscsDB, newDBErr := db.NewRscsDB(memoryDBName)
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

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
