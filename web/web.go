package web

import (
	"encoding/json"
	"fmt"
	"github.com/137675423/public/base"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewWeb() *Web {
	return &Web{true, make(chan struct{}, 10240), make(map[string]handle)}
}

type Engine struct {
	Req *http.Request
	Res http.ResponseWriter
	Log *base.Logger
}

type handle func(*Engine)

type Web struct {
	//运行状态
	Run bool
	//Web请求数量
	Process chan struct{}
	routes  map[string]handle
}

func (web *Web) AddRoute(path string, r handle) {
	web.routes[path] = r
}

//监听正在处理的进程完成
func (web *Web) Wait() chan struct{} {
	web.Run = false
	for range time.Tick(time.Second) {
		wait := len(web.Process)

		if wait > 0 {
			fmt.Println(fmt.Sprintf("Wait [%d] Process Running", wait))
			continue
		}
		fmt.Println("All Process Finish,Success Close Http Server!")
		break
	}
	signal := make(chan struct{}, 1)
	signal <- struct{}{}
	return signal
}

func (web *Web) Listen(port int) {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 30, //设置3秒的写超时
		Handler:      web,
	}

	sig := make(chan os.Signal, 10)
	signal.Notify(sig)
	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGINT:
				//等待进行中的进程全部完成
				<-web.Wait()
				server.Close()
				time.Sleep(time.Second * 5)
				return
			default:
				fmt.Println(s.String())
			}
		}
	}()

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
	logger.Info(fmt.Sprintf("Received HTTP %s Request URL:%s", req.Method, req.RequestURI))

	defer func() {
		err := recover()
		if err != nil {
			logger.Err(err)
		}

		logger.SaveFile()
	}()

	//服务没有运行
	if !web.Run {
		logger.Waring("Web Wait Close")
		return
	}

	route, ok := web.routes[req.RequestURI]
	//计数
	web.Process <- struct{}{}
	logger.Info(fmt.Sprintf("Process Number +1,Curr Process %d", len(web.Process)))

	if ok {
		route(&Engine{req, rw, logger})
	} else {
		logger.Waring("Route Not Found")
	}
	//减少计数
	<-web.Process
	logger.Info(fmt.Sprintf("Process Number -1,Curr Process %d", len(web.Process)))
}
