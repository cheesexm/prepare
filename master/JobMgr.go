package master

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

//JobMgr 调度job
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	//GJobMgr 单例
	GJobMgr *JobMgr
)

//InitJobMgr 初始化JobMgr
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   Gconfig.EtcdEndPoints,
		DialTimeout: time.Duration(Gconfig.EtcdDialTimeout) * time.Millisecond,
	}
	fmt.Println("Gconfig", Gconfig)
	if client, err = clientv3.New(config); err != nil {
		return
	}
	//得到KV和Lease的api自己
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	GJobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	fmt.Println("GJobMgr", GJobMgr)
	return
}
