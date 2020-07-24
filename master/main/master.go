package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/cheesexm/prepar/master"
)

var (
	conFile string //配置文件路径
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	//master -config ./master.json
	flag.StringVar(&conFile, "config", "./master.json", "指定master.json")

}

func main() {
	var (
		err error
	)
	fmt.Println("hello world")
	initEnv()
	initArgs()
	//加载配置
	if err = master.InitConfig(conFile); err != nil {
		goto ERR
	}
	if err = master.InitAPIServe(); err != nil {
		goto ERR
	}
	return
ERR:
	fmt.Println(err)
}
