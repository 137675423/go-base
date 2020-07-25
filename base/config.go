package base

import (
	"database/sql"
	"fmt"
	"time"
)

//数据库连接配置
type DbConfig struct {
	Host            string
	Port            int
	User            string
	Pwd             string
	DbName          string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

func (d *DbConfig) MySqlConn() (mysql *MySql, err error) {
	if mySqlHandle.db != nil {
		return mySqlHandle, nil
	}
	source := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8`, d.User, d.Pwd, d.Host, d.Port, d.DbName)
	db, err := sql.Open("mysql", source)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(d.MaxOpenConn)
	db.SetMaxIdleConns(d.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(d.ConnMaxLifetime))
	mySqlHandle.db = db
	return mySqlHandle, nil
}

func (d *DbConfig) PSqlConn() error {
	if PSqlHandle.db != nil {
		return nil
	}
	source := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8`, d.User, d.Pwd, d.Host, d.Port, d.DbName)
	mysql, err := sql.Open("psql", source)
	if err != nil {
		return err
	}
	mysql.SetMaxOpenConns(d.MaxOpenConn)
	mysql.SetMaxIdleConns(d.MaxIdleConn)
	mysql.SetConnMaxLifetime(time.Second * time.Duration(d.ConnMaxLifetime))
	PSqlHandle.db = mysql
	return nil
}
