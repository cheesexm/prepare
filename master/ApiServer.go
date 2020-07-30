package master

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

//APIServer 任务的http接口
type APIServer struct {
	httpServer *http.Server
}

//创建单例
var (
	GAPIServer *APIServer
)

func handleJobSave(w http.ResponseWriter, r *http.Request) {
	//233
	fmt.Println("233")
}

//InitAPIServe 初始化api服务
func InitAPIServe() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	fmt.Println(Gconfig)

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(Gconfig.APIPort)); err != nil {
		return
	}

	//创建一个http服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(Gconfig.APIReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Gconfig.APIWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	GAPIServer = &APIServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)
	return
}
