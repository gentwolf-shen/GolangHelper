package driver

import (
	_ "github.com/lib/pq"
)

type PostgreSqlHelper struct {
	Base
}

func (this *PostgreSqlHelper) Open(dsn string, maxOpenConnections int, maxIdleConnections int) error {
	return this.OpenDb("postgres", dsn, maxOpenConnections, maxIdleConnections)
}
