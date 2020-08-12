package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cheesexm/prepar/common"
	//"github.com/cheesexm/prepar/common"
)

//APIServer 任务的http接口
type APIServer struct {
	httpServer *http.Server
}

//创建单例
var (
	GAPIServer *APIServer
)

//保存任务接口
//POST job {name:"job1","command":"echo hello","cronExpr":"*****"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		ReadAll []byte
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	//2、取表单中的job字段
	ReadAll, _ = ioutil.ReadAll(req.Body)
	postJob = string(ReadAll)
	//3、反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	//4、保存到etcd
	if oldJob, err = GJobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	//5、返回正常应答（{"errno":0,"msg":"","data"""}）
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return

	//6、返回异常应答
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}

}

func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	//删除的任务名
	fmt.Println("删除", req.PostForm)
	fmt.Println("删除name", req.PostForm.Get("name"))
	name = req.PostForm.Get("name")
	//去删除任务
	if oldJob, err = GJobMgr.DeleteJob(name); err != nil {
		fmt.Println("删除6")
		goto ERR
	}
	fmt.Println("删除2")
	fmt.Println("删除oldJob", oldJob)
	//正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		fmt.Println("删除4")
		resp.Write(bytes)
	}
	return
ERR:
	fmt.Println("删除3")
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		fmt.Println("删除5")
		resp.Write(bytes)
	}

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
	mux.HandleFunc("/job/delete", handleJobDelete)
	fmt.Println("mux", mux)

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(Gconfig.APIPort)); err != nil {
		return
	}
	fmt.Println("listener", listener.Addr())

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
