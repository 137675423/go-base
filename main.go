package main

import (
	"fmt"
	"public/base"
	"time"
)

//对应数据库表结构
type Model struct {
	Id         int       `sql:"id"`
	Title      string    `sql:"title"`
	CreateTime time.Time `sql:"create_time"`
	Stat       int       `sql:"stat"`
}

//例
func main() {
	db := new(base.DbConfig)
	db.Host = "127.0.0.1"
	db.Port = 3306
	db.DbName = "tom"
	db.Pwd = "tom.123456"
	db.User = "tom"
	db.MaxOpenConn = 20
	db.MaxIdleConn = 5
	db.ConnMaxLifetime = 60
	mysql, err := db.MySqlConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	mysql.Ping()

	sm := new(base.SqlModify).SetTable("tom_test").And("title", base.Like, "h").And("id", base.Equal, 1).Or("stat", base.Equal, 2).SetModel(Model{})
	list, _, _ := mysql.List(sm)
	for _, v := range list {
		switch v.(type) {
		case *Model:
			fmt.Println(v.(*Model))
		}
	}
}
