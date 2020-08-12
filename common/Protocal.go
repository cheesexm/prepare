package common

import (
	"encoding/json"
	"fmt"
)

//Job 定时
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//http接口应答
type Response struct {
	Errno int         `json:"errno"` //错误no
	Msg   string      `json:"msg"`   //错误no
	Data  interface{} `json:"data"`  //data
}

//BuildResponse 返回json
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	//1、定义一个response
	fmt.Println("BuildResponse")
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data
	resp, err = json.Marshal(response)
	fmt.Println("resp", resp)
	fmt.Println("err", err)
	return
}
