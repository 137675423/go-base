package base

import "database/sql"

var PSqlHandle = new(PSql)

type PSql struct {
	db *sql.DB
}

func (p *PSql) List(sm *SqlModify) []interface{} {
	return nil
}

func (p *PSql) Update(sm *SqlModify) (change int, err error) {
	return
}

func (p *PSql) Delete(sm *SqlModify) (change int, err error) {
	return
}

func (p *PSql) Insert(sm *SqlModify) (id int, err error) {
	return
}
