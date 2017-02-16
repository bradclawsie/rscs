package server

import (
	"bytes"
	"encoding/json"
	"flag"
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
		t.Errorf("status:not 200")
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

func TestInsert(t *testing.T) {
	key := "key1"
	val := "val1"
	vIn := Value{Value: val}
	vJSON, jsonErr := json.Marshal(vIn)
	if jsonErr != nil {
		t.Errorf(jsonErr.Error())
	}

	emptyKeyRoute := KVRoutePrefix
	emptyKeyResp, _ := testRequest(t, testServer, http.MethodPost, emptyKeyRoute, bytes.NewReader(vJSON))
	if emptyKeyResp.StatusCode != http.StatusNotFound {
		t.Errorf("inserted empty route")
	}

	noKeyRoute := KVRoutePrefix + "/"
	noKeyResp, _ := testRequest(t, testServer, http.MethodPost, noKeyRoute, bytes.NewReader(vJSON))
	if noKeyResp.StatusCode == http.StatusCreated {
		t.Errorf("inserted empty key")
	}

	route := KVRoutePrefix + "/" + key
	badBytes := []byte(`{"V":"val"}`)
	badJSONResp, _ := testRequest(t, testServer, http.MethodPost, route, bytes.NewReader(badBytes))
	if badJSONResp.StatusCode == http.StatusCreated {
		t.Errorf("inserted bad JSON")
	}

	insertResp, _ := testRequest(t, testServer, http.MethodPost, route, bytes.NewReader(vJSON))
	if insertResp.StatusCode != http.StatusCreated {
		t.Errorf("insert:not 201")
	}

	getResp, getBody := testRequest(t, testServer, http.MethodGet, route, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Errorf("get:not 200")
	}
	var vOut Value
	umErr := json.Unmarshal([]byte(getBody), &vOut)
	if umErr != nil {
		t.Errorf(umErr.Error())
	}
	if vIn.Value != vOut.Value {
		t.Errorf("round trip values not equal")
	}
}

func TestUpdate(t *testing.T) {
	key := "key2"
	val := "val2"
	vIn := Value{Value: val}
	vJSON, jsonErr := json.Marshal(vIn)
	if jsonErr != nil {
		t.Errorf(jsonErr.Error())
	}

	route := KVRoutePrefix + "/" + key
	insertResp, _ := testRequest(t, testServer, http.MethodPost, route, bytes.NewReader(vJSON))
	if insertResp.StatusCode != http.StatusCreated {
		t.Errorf("insert:not 201")
	}

	getResp, getBody := testRequest(t, testServer, http.MethodGet, route, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Errorf("get:not 200")
	}
	var vOut Value
	umErr := json.Unmarshal([]byte(getBody), &vOut)
	if umErr != nil {
		t.Errorf(umErr.Error())
	}
	if vIn.Value != vOut.Value {
		t.Errorf("round trip values not equal")
	}

	emptyKeyRoute := KVRoutePrefix
	emptyKeyResp, _ := testRequest(t, testServer, http.MethodPut, emptyKeyRoute, bytes.NewReader(vJSON))
	if emptyKeyResp.StatusCode != http.StatusNotFound {
		t.Errorf("updated empty route")
	}

	noKeyRoute := KVRoutePrefix + "/"
	noKeyResp, _ := testRequest(t, testServer, http.MethodPut, noKeyRoute, bytes.NewReader(vJSON))
	if noKeyResp.StatusCode == http.StatusOK {
		t.Errorf("updated empty key")
	}

	badBytes := []byte(`{"V":"val"}`)
	badJSONResp, _ := testRequest(t, testServer, http.MethodPut, route, bytes.NewReader(badBytes))
	if badJSONResp.StatusCode == http.StatusOK {
		t.Errorf("updated bad JSON")
	}

	newVal := "val2new"
	vUpdate := Value{Value: newVal}
	vJSON, jsonErr = json.Marshal(vUpdate)
	if jsonErr != nil {
		t.Errorf(jsonErr.Error())
	}

	updateResp, _ := testRequest(t, testServer, http.MethodPut, route, bytes.NewReader(vJSON))
	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("update:not 200")
	}

	getResp, getBody = testRequest(t, testServer, http.MethodGet, route, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Errorf("get:not 200")
	}
	var vUpdateGet Value
	umErr = json.Unmarshal([]byte(getBody), &vUpdateGet)
	if umErr != nil {
		t.Errorf(umErr.Error())
	}
	if vUpdate.Value != vUpdateGet.Value {
		t.Errorf("round trip values not equal")
	}
}

func TestDelete(t *testing.T) {
	key := "key3"
	val := "val3"
	vIn := Value{Value: val}
	vJSON, jsonErr := json.Marshal(vIn)
	if jsonErr != nil {
		t.Errorf(jsonErr.Error())
	}

	route := KVRoutePrefix + "/" + key
	insertResp, _ := testRequest(t, testServer, http.MethodPost, route, bytes.NewReader(vJSON))
	if insertResp.StatusCode != http.StatusCreated {
		t.Errorf("insert:not 201")
	}

	getResp, getBody := testRequest(t, testServer, http.MethodGet, route, nil)
	if getResp.StatusCode != http.StatusOK {
		t.Errorf("get:not 200")
	}
	var vOut Value
	umErr := json.Unmarshal([]byte(getBody), &vOut)
	if umErr != nil {
		t.Errorf(umErr.Error())
	}
	if vIn.Value != vOut.Value {
		t.Errorf("round trip values not equal")
	}

	emptyKeyRoute := KVRoutePrefix
	emptyKeyResp, _ := testRequest(t, testServer, http.MethodDelete, emptyKeyRoute, nil)
	if emptyKeyResp.StatusCode != http.StatusNotFound {
		t.Errorf("deleted on empty route")
	}

	noKeyRoute := KVRoutePrefix + "/"
	noKeyResp, _ := testRequest(t, testServer, http.MethodDelete, noKeyRoute, nil)
	if noKeyResp.StatusCode == http.StatusOK {
		t.Errorf("deleted empty key")
	}

	deleteResp, _ := testRequest(t, testServer, http.MethodDelete, route, nil)
	if deleteResp.StatusCode != http.StatusOK {
		t.Errorf("delete:not 200")
	}

	getResp, _ = testRequest(t, testServer, http.MethodGet, route, nil)
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("get:not 404")
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

func Example() {

	// You need a sqlite file. We assume here that you have one with
	// a pre-existing `kv` table, if not, look at rscs.go in this directory
	// for an example of how to create an empty one.
	var sqliteDBFile string
	flag.StringVar(&sqliteDBFile, "db", "", "full path to sqlite db file")
	flag.Parse()

	if sqliteDBFile == "" {
		log.Fatal(`use: example_daemon --db={sqlite db file}`)
	}

	// Create a DB instance for sqlite.
	rscsDB, rscsDBErr := db.NewRscsDB(sqliteDBFile)
	if rscsDBErr != nil {
		log.Fatal(rscsDBErr)
	}

	// Create an http abstraction layer.
	rscsServer, rscsSrvErr := NewRscsServer(rscsDB)
	if rscsSrvErr != nil {
		log.Fatal(rscsSrvErr)
	}

	// Instantiate a router.
	rtr, rtrErr := rscsServer.NewRouter()
	if rtrErr != nil {
		log.Fatal(rtrErr.Error())
	}

	// Listen.
	srv := &http.Server{Addr: ":8081", Handler: rtr}
	log.Fatal(srv.ListenAndServe())
}
