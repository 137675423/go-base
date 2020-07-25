package base

import (
	"database/sql"
)

var mySqlHandle = new(MySql)

type MySql struct {
	db *sql.DB
}

func (m *MySql) Ping() error {
	return m.db.Ping()
}

func (m *MySql) List(sm *SqlModify) ([]interface{}, int, error) {
	sql, val := sm.QueryList()
	list, t, e := sm.queryList(m.db, sql, val...)
	return sm.Out(list), t, e
}

func (m *MySql) Update(sm *SqlModify) (change int, err error) {
	return
}

func (m *MySql) Delete(sm *SqlModify) (change int, err error) {
	return
}

func (m *MySql) Insert(sm *SqlModify) (id int, err error) {
	return
}
