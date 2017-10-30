package databaseHelper

import (
	"gentwolf/GolangHelper/convert"
	"gentwolf/GolangHelper/database"
	"gentwolf/GolangHelper/logger"
)

type DatabaseHelper struct {
	db *database.Base
}

// 检测表中某项值是否存在
func (this *DatabaseHelper) TableItemExists(table, key string, value interface{}, keyId int64) bool {
	name := table + "TableItemExists"

	sql := "SELECT id FROM " + table + " WHERE "
	if keyId > 0 {
		sql += " AND id!=" + convert.ToStr(keyId)

		name += "Key"
	}
	sql += " AND " + key + "=? LIMIT 1"

	id, err := this.db.PrepareQueryScalar(name, sql, "id", value)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
		return true
	}

	return id != ""
}

// 修改表中某值
func (this *DatabaseHelper) TableItemUpdate(table, key, value string, keyId int64) bool {
	name := table + "TableItemUpdate"

	sql := "UPDATE " + table + " SET " + key + "=? WHERE user_id=? AND id=?"
	_, err := this.db.PrepareExec(name, sql, value, keyId)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
		return false
	}

	return true
}

// 删除表中记录
func (this *DatabaseHelper) TableItemDelete(table string, keyId interface{}) bool {
	name := table + "TableItemDelete"

	sql := "DELETE FROM " + table + " WHERE id=?"
	_, err := this.db.PrepareExec(name, sql, keyId)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
		return false
	}

	return true
}

// 添加内容到表中
func (this *DatabaseHelper) TableItemInsert(name, table string, params map[string]interface{}) int64 {
	id, err := this.db.PrepareInsert(name, table, params)
	if err != nil {
		logger.Out.Error.Println("\t", err)
		return 0
	}

	return id
}

// 取表中记录 [简单查询]
func (this *DatabaseHelper) TableItems(name, sql string, params ...interface{}) []map[string]string {
	rows, err := this.db.PrepareQuery(name, sql, params...)
	if err != nil {
		logger.Out.Error.Println(sql)
		logger.Out.Error.Println("\t", err)
	}
	return rows
}
