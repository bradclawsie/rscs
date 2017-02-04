package server

import (
	"github.com/bradclawsie/rscs/db"
)

type RscsServer struct {
	rscsDB *db.RscsDB
}

func NewRscsServer(sqliteDBFile string, readOnly bool) (*RscsServer, error) {
	rscsDB, dbErr := db.NewRscsDB(sqliteDBFile, readOnly)
	if dbErr != nil {
		return dbErr, nil
	}
	return &RscsServer{rscsDB:rscsDB}, nil
}
