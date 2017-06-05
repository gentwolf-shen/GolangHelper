package database

import (
	"database/sql"
	"strings"
)

type Base struct {
	dbConn   *sql.DB
	stmtList map[string]*sql.Stmt
}

func (this *Base) Version() string {
	return "1.0"
}

func (this *Base) OpenDb(dbType string, dsn string, maxOpenConnections int, maxIdleConnections int) error {
	var err error
	this.dbConn, err = sql.Open(dbType, dsn)
	if err == nil {
		this.dbConn.SetMaxOpenConns(maxIdleConnections)
		this.dbConn.SetMaxIdleConns(maxIdleConnections)

		this.stmtList = make(map[string]*sql.Stmt)
	}

	return err
}

func (this *Base) Close() {
	this.dbConn.Close()
	for _, item := range this.stmtList {
		item.Close()
	}
}

func (this *Base) GetConn() *sql.DB {
	return this.dbConn
}

func (this *Base) Insert(sql string, args ...interface{}) (int64, error) {
	var id int64
	var err error

	if strings.Contains(sql, " RETURNING ") {
		row := this.dbConn.QueryRow(sql, args...)
		err = row.Scan(&id)
	} else {
		result, err1 := this.dbConn.Exec(sql, args...)
		err = err1
		if err1 == nil {
			id, _ = result.LastInsertId()
		}
	}

	return id, err
}

func (this *Base) Update(sql string, args ...interface{}) (int64, error) {
	var n int64
	var err error

	result, err := this.dbConn.Exec(sql, args...)
	if err == nil {
		n, _ = result.RowsAffected()
	}

	return n, err
}

func (this *Base) Delete(sql string, args ...interface{}) (int64, error) {
	return this.Update(sql, args...)
}

func (this *Base) Query(sql string, args ...interface{}) ([]map[string]string, error) {
	rows, err := this.dbConn.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return this.fetchRows(rows, err)
}

func (this *Base) QueryRow(sql string, args ...interface{}) (map[string]string, error) {
	rows, err := this.Query(sql, args...)

	if rows != nil && err == nil && len(rows) > 0 {
		return rows[0], err
	}

	return nil, err
}

func (this *Base) QueryScalar(sql string, key string, args ...interface{}) (string, error) {
	rows, err := this.Query(sql, args...)

	if rows != nil && err == nil && len(rows) > 0 {
		row := rows[0]
		if value, ok := row[key]; ok {
			return value, err
		}
	}
	return "", err
}

func (this *Base) fetchRows(rows *sql.Rows, err error) ([]map[string]string, error) {
	if rows == nil || err != nil {
		return nil, err
	}

	fields, _ := rows.Columns()
	for k, v := range fields {
		fields[k] = this.camelCase(v)
	}
	columnsLength := len(fields)

	values := make([]string, columnsLength)
	args := make([]interface{}, columnsLength)
	for i := 0; i < columnsLength; i++ {
		args[i] = &values[i]
	}

	index := 0
	listLength := 100
	lists := make([]map[string]string, listLength, listLength)
	for rows.Next() {
		if e := rows.Scan(args...); e == nil {
			row := make(map[string]string, columnsLength)
			for i, field := range fields {
				row[field] = string(values[i])
			}

			if index < listLength {
				lists[index] = row
			} else {
				lists = append(lists, row)
			}
			index++
		}
	}

	return lists[0:index], nil
}

func (this *Base) camelCase(str string) string {
	if strings.Contains(str, "_") {
		items := strings.Split(str, "_")
		arr := make([]string, len(items))
		for k, v := range items {
			if 0 == k {
				arr[k] = v
			} else {
				arr[k] = strings.Title(v)
			}
		}
		str = strings.Join(arr, "")
	}

	return str
}

func (this *Base) PrepareSql(name, sql string) (*sql.Stmt, error) {
	stmt, bl := this.stmtList[name]
	if !bl {
		var err error
		stmt, err = this.dbConn.Prepare(sql)
		if err != nil {
			return nil, err
		}

		this.stmtList[name] = stmt
	}

	return stmt, nil
}

func (this *Base) PrepareQuery(name, sql string, args ...interface{}) ([]map[string]string, error) {
	stmt, err := this.PrepareSql(name, sql)
	if err != nil {
		return nil, err
	}

	rows, err1 := stmt.Query(args...)
	defer rows.Close()

	return this.fetchRows(rows, err1)
}

func (this *Base) PrepareQueryRow(name, sql string, args ...interface{}) (map[string]string, error) {
	rows, err := this.PrepareQuery(name, sql, args...)

	if rows != nil && err == nil && len(rows) > 0 {
		return rows[0], err
	}

	return nil, err
}

func (this *Base) PrepareQueryScalar(name, sql string, args ...interface{}) (string, error) {
	stmt, err := this.PrepareSql(name, sql)
	if err != nil {
		return "", err
	}

	var value string
	rows, err1 := stmt.Query(args...)
	if err1 != nil {
		return "", err1
	}

	defer rows.Close()
	if rows.Next() {
		rows.Scan(&value)
	}

	return value, err
}

func (this *Base) PrepareExec(name, sql string, args ...interface{}) (int64, error) {
	stmt, err := this.PrepareSql(name, sql)
	if err != nil {
		return 0, err
	}

	result, err1 := stmt.Exec(args...)
	if err1 != nil {
		return 0, err1
	}

	n := int64(0)
	if "INSERT" == strings.ToUpper(string(sql[0:6])) {
		//n, err = result.LastInsertId()
	} else {
		n, err = result.RowsAffected()
	}

	return n, err
}
