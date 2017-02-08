package main

import (
	"flag"
	"github.com/bradclawsie/rscs/db"
	"github.com/bradclawsie/rscs/server"
	"log"
)

func main() {

	const use = `use: rscs --db={sqlite db file}`

	// Command line options.
	var sqliteDBFile, tlsCertFile, tlsKeyFile string
	flag.StringVar(&sqliteDBFile, "db", "", "full path to sqlite db file")

	if sqliteDBFile == "" {
		log.Fatal(use)
	}
	log.Printf("[rscs --db=%s]", sqliteDBFile)

	rscsDB, rscsDBErr := db.NewRscsDB(sqliteDBFile)
	if rscsDBErr != nil {
		log.Fatal(rscsDBErr)
	}
	_, rscsSrvErr := server.NewRscsServer(rscsDB)
	if rscsSrvErr != nil {
		log.Fatal(rscsSrvErr)
	}
}
