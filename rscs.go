package main

import (
	"flag"
	"log"
	"github.com/bradclawsie/rscs/server"
)

func main() {

	const use = `use: rscs --cert={tsl cert pem file} --key={tls key pem file} --db={sqlite db file} [--readonly]`

	// Command line options.
	var sqliteDBFile, tlsCertFile, tlsKeyFile string
	flag.StringVar(&sqliteDBFile, "db", "", "full path to sqlite db file")
	flag.StringVar(&tlsCertFile, "cert", "", "full path to tls cert file")
	flag.StringVar(&tlsKeyFile, "key", "", "full path to tls key file")

	var readOnly bool
	flag.BoolVar(&readOnly, "readonly", true, "use rscs in readonly mode")
	flag.Parse()

	if sqliteDBFile == "" || tlsCertFile == "" || tlsKeyFile == "" {
		log.Fatal(use)
	}
	log.Printf("[rscs --cert=%s --key=%s --db=%s --readonly=%v]", tlsCertFile, tlsKeyFile, sqliteDBFile, readOnly)

	_,srvErr := server.NewRscs(sqliteDBFile)
	if srvErr != nil {
		log.Fatal(srvErr)
	}
}
