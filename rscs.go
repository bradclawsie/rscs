package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bradclawsie/rscs/db"
	"github.com/bradclawsie/rscs/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	const use = `use: rscs --db={sqlite db file} [--create-only] [--port={portnum}]`

	// Command line options.
	var sqliteDBFile string
	var createOnly bool
	var portNum int
	flag.StringVar(&sqliteDBFile, "db", "", "full path to sqlite db file")
	flag.BoolVar(&createOnly, "create-only", false, "only create table in file and exit")
	flag.IntVar(&portNum, "port", 8081, "port to listen on")
	flag.Parse()

	if sqliteDBFile == "" {
		log.Fatal(use)
	}

	rscsDB, rscsDBErr := db.NewRscsDB(sqliteDBFile)
	if rscsDBErr != nil {
		log.Fatal(rscsDBErr)
	}

	if createOnly {
		createErr := rscsDB.CreateTable()
		if createErr != nil {
			log.Fatal(createErr.Error())
		}
		log.Printf("created")
		os.Exit(0)
	}

	rscsServer, rscsSrvErr := server.NewRscsServer(rscsDB)
	if rscsSrvErr != nil {
		log.Fatal(rscsSrvErr)
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	rtr, rtrErr := rscsServer.NewRouter()
	if rtrErr != nil {
		log.Fatal(rtrErr.Error())
	}

	addrStr := fmt.Sprintf(":%d", portNum)
	srv := &http.Server{Addr: addrStr, Handler: rtr}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
