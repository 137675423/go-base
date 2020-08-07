package main

import (
	"fmt"
	"os"
	"os/signal"
	"public/base"
	"public/web"
	"syscall"
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

	sig := make(chan os.Signal, 10)
	signal.Notify(sig)
	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGINT:
				fmt.Println("Ctrl+C Stop")
				//关闭web服务
				<-web.Server.Stop()
				fmt.Println("Web Stop")
				os.Exit(1)
			default:
				fmt.Println(s.String())
			}

		}
	}()

	web.Server.Listen()
	return

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
