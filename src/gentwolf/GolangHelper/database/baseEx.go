package database

import (
	"strings"

	"gentwolf/GolangHelper/convert"
	"gentwolf/GolangHelper/logger"
)

// 简单查询
func (this *Base) SimpleQuery(name, table, fields string, keyId int64) map[string]string {
	var row map[string]string
	var err error

	s := "SELECT " + fields + " FROM " + table + " WHERE id=?"

	if name == "" {
		row, err = this.QueryRow(s, keyId)
	} else {
		row, err = this.PrepareQueryRow(name, s, keyId)
	}

	if err != nil {
		logger.Out.Error.Println(s)
		logger.Out.Error.Println("\t", err)
		return nil
	}

	return row
}

// 更新/删除
func (this *Base) SimpleUpdate(name, sql string, args ...interface{}) int64 {
	var n int64 = -1
	var err error

	if name == "" {
		n, err = this.Update(sql, args...)
	} else {
		n, err = this.PrepareExec(name, sql, args...)
	}

	if err != nil {
		logger.Out.Error.Println(err)
		return -1
	}
	return n
}

// 检测表中值是否存在
func (this *Base) TableItemExists(table, key string, value interface{}, keyId int64) bool {
	strKey := convert.ToStr(keyId)
	name := table + "TableItemExists" + key + "-" + strKey

	s := "SELECT id FROM " + table + " WHERE "
	if keyId > 0 {
		s += " id!=" + strKey + " AND "
		name += "Key"
	}
	s += key + "=? LIMIT 1"

	id, err := this.PrepareQueryScalar(name, s, value)
	if err != nil {
		logger.Out.Error.Println(s)
		logger.Out.Error.Println("\t", err)
		return true
	}

	return id != ""
}

// 修改表中某值
func (this *Base) TableItemUpdate(table, key, value string, keyId int64) bool {
	name := table + "TableItemUpdate" + key

	sql := "UPDATE " + table + " SET " + key + "=? WHERE id=?"
	_, err := this.PrepareExec(name, sql, value, keyId)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
		return false
	}

	return true
}

// 删除表中记录
func (this *Base) TableItemDelete(table, key string, keyId interface{}) bool {
	name := table + "TableItemDelete"

	sql := "DELETE FROM " + table + " WHERE " + key + "=?"
	_, err := this.PrepareExec(name, sql, keyId)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
		return false
	}

	return true
}

// 添加内容到表中
func (this *Base) TableInsert(name, table string, data map[string]interface{}) int64 {
	keys, params, values := this.formatInsertData(data)

	s := "INSERT INTO " + table + "(" + strings.Join(keys, ",") + ") VALUES(" + strings.Join(params, ",") + ")"
	if this.dbType == "postgres" {
		s += "RETURNING ID"
	}

	id, err := this.PrepareExec(name, s, values...)
	if err != nil {
		logger.Out.Error.Println(s)
		logger.Out.Error.Println("\t", err)
		return 0
	}

	return id
}

// 修改内容到表中
func (this *Base) TableUpdate(name, table string, data map[string]interface{}, where map[string]interface{}) int64 {
	dataKeys, whereKeys, values := this.formUpdateData(data, where)

	s := "UPDATE " + table + " SET " + strings.Join(dataKeys, ",") + " WHERE " + strings.Join(whereKeys, " AND ")
	n, err := this.PrepareExec(name, s, values...)

	if err != nil {
		logger.Out.Error.Println(s)
		logger.Out.Error.Println("\t", err)
		return -1
	}

	return n
}

// 取表中记录 [简单查询]
func (this *Base) TableItems(name, sql string, params ...interface{}) []map[string]string {
	rows, err := this.PrepareQuery(name, sql, params...)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
	}
	return rows
}

func (this *Base) formatInsertData(data map[string]interface{}) ([]string, []string, []interface{}) {
	length := len(data)
	keys := make([]string, length)
	params := make([]string, length)
	values := make([]interface{}, length)
	index := 0
	for k, v := range data {
		keys[index] = k

		if this.dbType == "mysql" {
			params[index] = "?"
		} else {
			params[index] = "$" + convert.ToStr(index)
		}

		values[index] = v

		index++
	}

	return keys, params, values
}

func (this *Base) formUpdateData(data map[string]interface{}, where map[string]interface{}) ([]string, []string, []interface{}) {
	dataLength := len(data)
	dataKeys := make([]string, dataLength)

	whereLength := len(where)
	whereKeys := make([]string, whereLength)

	values := make([]interface{}, dataLength+whereLength)

	index := 0
	for k, v := range data {
		if this.dbType == "mysql" {
			dataKeys[index] = k + "=?"
		} else {
			dataKeys[index] = k + "=$" + convert.ToStr(index)
		}

		values[index] = v
		index++
	}

	i := 0
	for k, v := range where {
		if this.dbType == "mysql" {
			whereKeys[i] = k + "=?"
		} else {
			whereKeys[i] = k + "=$" + convert.ToStr(index)
		}

		values[index] = v

		i++
		index++
	}

	return dataKeys, whereKeys, values
}
