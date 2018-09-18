package sql_builder

import (
	"strings"
)

var (
	DbOracle     = "ORACLE"
	DbMySQL      = "MYSQL"
	DbPostgreSQL = "POSTGRESQL"
)

type SqlBuilder struct {
	dbType  string
	sqlType string
	table   string

	wheres       []string
	selectFields []string
	updateParams []string
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{}
}

func (this *SqlBuilder) SetDbType(name string) *SqlBuilder {
	this.dbType = name
	return this
}

func (this *SqlBuilder) Select(args ...string) *SqlBuilder {
	this.sqlType = "SELECT "
	this.selectFields = make([]string, len(args))
	for i, name := range args {
		switch this.dbType {
		case DbOracle:
			this.selectFields[i] = `"` + name + `"`
		default:
			this.selectFields[i] = name
		}
	}

	return this
}

func (this *SqlBuilder) Update(args ...string) *SqlBuilder {
	this.sqlType = "UPDATE "
	this.updateParams = this.buildWhere(args...)

	return this
}

func (this *SqlBuilder) Delete(args ...string) *SqlBuilder {
	this.sqlType = "DELETE "
	if len(args) > 0 {
		this.wheres = this.buildWhere(args...)
	}
	return this
}

func (this *SqlBuilder) Insert(args ...string) *SqlBuilder {
	this.sqlType = "INSERT INTO "
	this.updateParams = this.buildWhere(args...)
	return this
}

func (this *SqlBuilder) From(table string) *SqlBuilder {
	this.table = table
	return this
}

func (this *SqlBuilder) Where(args ...string) *SqlBuilder {
	this.wheres = this.buildWhere(args...)

	return this
}

func (this *SqlBuilder) ToString() string {
	str := this.sqlType
	if len(this.selectFields) > 0 {
		str += strings.Join(this.selectFields, ",") + " FROM "
	}
	str += this.table + " "

	if len(this.updateParams) > 0 {
		str += strings.Join(this.updateParams, ",")
	}

	if len(this.wheres) > 0 {
		str += " WHERE " + strings.Join(this.wheres, " AND ")
	}
	return str
}

func (this *SqlBuilder) buildWhere(args ...string) []string {
	params := make([]string, len(args))
	for i, name := range args {
		switch this.dbType {
		case DbOracle:
			params[i] = `"` + name + `"=:1`
		default:
			this.wheres[i] = name + `=:1`
		}
	}

	return params
}
