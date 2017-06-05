package PostgresSqlHelper

import (
	"gentwolf/GolangHelper/database"
	_ "github.com/lib/pq"
)

type Driver struct {
	database.Base
}

func (this *Driver) Open(dsn string, maxOpenConnections int, maxIdleConnections int) error {
	return this.OpenDb("postgres", dsn, maxOpenConnections, maxIdleConnections)
}
