package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"public/base"
	"time"
)

var Server = NewWeb()

func NewWeb() *Web {
	return &Web{true, make(chan struct{}, 10240)}
}

type Web struct {
	//运行状态
	Run bool
	//Web请求数量
	Process chan struct{}
}

func (web *Web) SeeProcess() {
	fmt.Println(len(web.Process), " Process Running")
}

//监听服务关闭
func (web *Web) Stop() chan struct{} {
	web.Run = false
	for range time.Tick(time.Second) {
		wait := len(web.Process)
		fmt.Println(wait, " Process Running")
		if wait > 0 {
			continue
		}
		break
	}
	signal := make(chan struct{}, 1)
	signal <- struct{}{}
	return signal
}

func (web *Web) Listen() {
	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 30, //设置3秒的写超时
		Handler:      web,
	}

	log.Fatal(server.ListenAndServe())
}

func (web *Web) GetParam(req *http.Request) (map[string]interface{}, error) {
	b, e := ioutil.ReadAll(req.Body)
	if e != nil {
		return nil, e
	}
	var post map[string]interface{}

	return post, json.Unmarshal(b, &post)
}

func (web *Web) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	//每个请求都生成一条日志,途径路线日志全写入此日志对象
	logger := base.NewLogger()
	logger.Info(req.RequestURI, req.Method)
	defer logger.SaveFile()

	defer func() {
		logger.Err(recover())
	}()

	//服务没有运行
	if !web.Run {
		logger.Waring("Web Wait Close")
		return
	}

	panic("2")

	//计数
	web.Process <- struct{}{}
	fmt.Println(len(web.Process), " Process Running")

	//减少计数
	<-web.Process
}
