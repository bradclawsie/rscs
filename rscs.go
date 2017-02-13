package main

import (
	"context"
	"flag"
	"github.com/bradclawsie/rscs/db"
	"github.com/bradclawsie/rscs/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	const use = `use: rscs --db={sqlite db file}`

	// Command line options.
	var sqliteDBFile string
	flag.StringVar(&sqliteDBFile, "db", "", "full path to sqlite db file")
	flag.Parse()

	if sqliteDBFile == "" {
		log.Fatal(use)
	}
	log.Printf("[rscs --db=%s]", sqliteDBFile)

	rscsDB, rscsDBErr := db.NewRscsDB(sqliteDBFile)
	if rscsDBErr != nil {
		log.Fatal(rscsDBErr)
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

	srv := &http.Server{Addr: ":8081", Handler: rtr}

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
