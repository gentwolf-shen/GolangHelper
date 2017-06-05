package Sqlite3Helper

import (
	"gentwolf/GolangHelper/database"
	_ "github.com/mattn/go-sqlite3"
)

type Driver struct {
	database.Base
}

func (this *Driver) Open(dsn string, maxOpenConnections int, maxIdleConnections int) error {
	return this.OpenDb("sqlite3", dsn, maxOpenConnections, maxIdleConnections)
}
