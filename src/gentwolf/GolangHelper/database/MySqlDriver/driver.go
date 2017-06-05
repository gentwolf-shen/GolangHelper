package MySqlDriver

import (
	"gentwolf/GolangHelper/database"
	_ "github.com/Go-SQL-Driver/MySQL"
)

type Driver struct {
	database.Base
}

func (this *Driver) Open(dsn string, maxOpenConnections int, maxIdleConnections int) error {
	return this.OpenDb("mysql", dsn, maxOpenConnections, maxIdleConnections)
}
