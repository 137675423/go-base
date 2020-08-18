package main

import (
	"fmt"
	"public/base"
	"public/web"
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

	Web := web.NewWeb()
	Web.AddRoute("/", test)
	go Web.Listen(1210)
	select {}

}

func test(engine *web.Engine) {

	time.Sleep(time.Second * 20)
	engine.Log.Info("DO TEST")
}

func DBExample() {
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

	//需要建立MySqlConn之后，才能调用MysqlQuery方法
	list, total, err := new(base.SqlModify).SetTable("tom_test").And("title", base.Like, "h").And("id", base.Equal, 1).Or("stat", base.Equal, 2).MysqlQuery(Model{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("total", total)
	for _, v := range list {
		switch v.(type) {
		case *Model:
			fmt.Println("elem", v.(*Model))
		}
	}
}
