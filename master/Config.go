package master

import (
	"encoding/json"
	"io/ioutil"
)

//Config 读取json配置
type Config struct {
	APIPort         int `json:"APIPort"`
	APIReadTimeout  int `json:"APIReadTimeout"`
	APIWriteTimeout int `json:"APIWriteTimeout"`
}

//config单例
var (
	Gconfig *Config
)

//InitConfig 初始化单例
func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	//1、把配置读取进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//2、做json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	Gconfig = &conf
	return
}
