package driver

import (
	_ "github.com/Go-SQL-Driver/MySQL"
)

type MySqlHelper struct {
	Base
}

func (this *MySqlHelper) Open(dsn string, maxOpenConnections int, maxIdleConnections int) error {
	return this.OpenDb("mysql", dsn, maxOpenConnections, maxIdleConnections)
}
