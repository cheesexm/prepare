package master

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cheesexm/prepar/common"
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

func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	//把任务保存到/cron/jobs/任务名
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	fmt.Println("job.Name", job.Name)
	//etcd的保存key
	jobKey = "/cron/jobs/" + job.Name
	//任务信息json
	fmt.Println(1)
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	fmt.Println(2)
	//保存到etcd
	if putResp, err = JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	fmt.Println(3)
	fmt.Println("putResp.PrevKv", putResp.PrevKv)
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {

			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	fmt.Println(4)
	return
}
func (JobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	//etcd中任务的key
	jobKey = "/cron/jobs/" + name
	//从etcd中删除它
	fmt.Println("DeleteJobname", jobKey)
	if delResp, err = JobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	fmt.Println("DeleteJobdelResp", delResp)
	fmt.Println("DeleteJobdelRespPrevKvs", delResp.PrevKvs)
	//返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		fmt.Println("DeleteJobdelRespValue", delResp.PrevKvs[0].Value)
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return

}
