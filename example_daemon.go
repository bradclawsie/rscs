// This example is very similar to rscs.go with some code removed for simplicity.
// This requires Go 1.8* to compile.
package main

import (
	"flag"
	"github.com/bradclawsie/rscs/db"
	"github.com/bradclawsie/rscs/server"
	"log"
	"net/http"
)

func main() {

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
	rscsServer, rscsSrvErr := server.NewRscsServer(rscsDB)
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
