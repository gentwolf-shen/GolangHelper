package database

import (
	"strings"

	"gentwolf/GolangHelper/convert"
)

type Helper struct {
	driver string
	dbType string
	db     *Base
	caches map[string]bool
}

func NewHelper(driver string) *Helper {
	helper := &Helper{}
	helper.driver = driver

	helper.db = Factory.Driver(driver)
	helper.dbType = helper.db.DbType

	return helper
}

func (this *Helper) QuerySql(name, sql string, params ...interface{}) ([]map[string]string, error) {
	return this.db.PrepareQuery(name, sql, params...)
}

func (this *Helper) Query(name, table, fields string, where map[string]interface{}, ext string) ([]map[string]string, error) {
	sql := ""
	params, _, values := this.formUpdateData(where, nil)

	if _, bl := this.db.StmtList[name]; !bl {
		sql := "SELECT " + fields + " FROM " + table

		if len(values) > 0 {
			sql += " WHERE " + strings.Join(params, " AND ")
		}
	}

	return this.db.PrepareQuery(name, sql, values...)
}

func (this *Helper) QueryRow(name, table, fields string, where map[string]interface{}, ext string) (map[string]string, error) {
	rows, err := this.Query(name, table, fields, where, ext)
	if err != nil || len(rows) == 0 {
		return nil, err
	}

	return rows[0], nil
}

func (this *Helper) QueryScalar(name, table, field string, where map[string]interface{}, ext string) (string, error) {
	row, err := this.QueryRow(name, table, field, where, ext)
	if err != nil || len(row) == 0 {
		return "", err
	}

	return row[field], nil
}

func (this *Helper) Insert(name, table string, data map[string]interface{}) (int64, error) {
	sql := ""
	keys, params, values := this.formatInsertData(data)

	if _, bl := this.db.StmtList[name]; !bl {
		sql = "INSERT INTO " + table + "(" + strings.Join(keys, ",") + ") VALUES(" + strings.Join(params, ",") + ")"

		if this.dbType == "postgres" {
			sql += "RETURNING ID"
		}
	}

	return this.db.Insert(sql, values...)
}

func (this *Helper) Update(name, table string, where, data map[string]interface{}) (int64, error) {
	sql := ""
	whereParams, dataParams, values := this.formUpdateData(where, data)

	if _, bl := this.db.StmtList[name]; !bl {
		sql = "UPDATE " + table + " SET " + strings.Join(whereParams, ",") + " WHERE " + strings.Join(dataParams, " AND ")
	}

	return this.db.Update(sql, values...)
}

func (this *Helper) Delete(name, table string, where map[string]interface{}) (int64, error) {
	sql := ""
	params, _, values := this.formUpdateData(where, nil)

	if _, bl := this.db.StmtList[name]; !bl {
		sql = "DELETE FROM " + table + " SET WHERE " + strings.Join(params, ",")
	}

	return this.db.Update(sql, values...)
}

func (this *Helper) formatInsertData(data map[string]interface{}) ([]string, []string, []interface{}) {
	length := len(data)
	keys := make([]string, length)
	params := make([]string, length)
	values := make([]interface{}, length)
	index := 0
	for k, v := range data {
		keys[index] = k

		if this.dbType == "mysql" {
			params[index] = "=?"
		} else {
			params[index] = "=$" + convert.ToStr(index)
		}

		values[index] = v

		index++
	}

	return keys, params, values
}

func (this *Helper) formUpdateData(where, data map[string]interface{}) ([]string, []string, []interface{}) {
	whereLength := len(where)
	dataLength := len(data)

	index := 0
	values := make([]interface{}, whereLength+dataLength)

	i := 0
	whereParams := make([]string, whereLength)
	for k, v := range where {
		if this.dbType == "mysql" {
			whereParams[i] = k + "=?"
		} else {
			whereParams[i] = k + "=$" + convert.ToStr(index)
		}
		values[index] = v

		i++
		index++
	}

	i = 0
	dataParams := make([]string, dataLength)
	for k, v := range data {
		if this.dbType == "mysql" {
			dataParams[i] = k + "=?"
		} else {
			dataParams[i] = k + "$" + convert.ToStr(index)
		}

		values[index] = v

		i++
		index++
	}

	return whereParams, dataParams, values
}
