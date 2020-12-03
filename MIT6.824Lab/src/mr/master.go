package mr

import (
	"fmt"
	"log"
	"sync"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type IMaster interface {
	RegisterWorker(args *RegisterArgs, reply *RegisterReply) error
	FetchWorker(args *FetchArgs, reply *FetchReply) error
	ReportWorker(args *RegisterArgs, reply *RegisterReply) error
	checkWorkerRunStatus()
}

type workerStatus struct {
	status        int
	fileIndex     int
	fetchWorkTime time.Time
}

type Master struct {
	// Your definitions here.
	fileNames            []string              // 文件名列表
	workerMap            map[int]*workerStatus // 工作状态
	nReduce              int                   // reduce数量
	workerId             int                   // id
	outputFileMap        [][]string            // 输出文件
	mapStart             int                   // map请求数量
	reduceStart          int                   // reduce开始数量
	workerIdMutex        sync.RWMutex          // id锁
	workerMapReduceMutex sync.RWMutex          // mr锁
}

// Your code here -- RPC handlers for the worker to call.

// RegisterWorker 注册worker 用于派发id
func (m *Master) RegisterWorker(args *RegisterArgs, reply *RegisterReply) error {
	m.workerIdMutex.Lock()
	reply.Id = m.workerId
	m.workerId++
	m.workerIdMutex.Unlock()
	fmt.Println("Worker注册成功！Id:", reply.Id)
	m.workerMap[reply.Id] = &workerStatus{
		status: 0,
	}
	return nil
}

func (m *Master) FetchWorker(args *FetchArgs, reply *FetchReply) error {
	if m.mapStart < len(m.fileNames) {
		m.workerMapReduceMutex.Lock()
		m.workerMap[args.Id].status = 1
		m.mapStart++
		reply.FileNames = []string{m.fileNames[m.mapStart-1]}
		m.workerMapReduceMutex.Unlock()
		reply.Status = 1
		reply.NReduce = m.nReduce
		return nil
	}
	if m.reduceStart < m.nReduce {
		m.workerMapReduceMutex.Lock()
		m.workerMap[args.Id].status = 2
		m.reduceStart++
		reply.FileNames = m.outputFileMap[m.reduceStart-1]
		m.workerMapReduceMutex.Unlock()
		reply.Status = 2
		reply.NReduce = m.nReduce
		return nil
	}
	reply.Status = 0
	reply.NReduce = m.nReduce
	return nil
}

func (m *Master) ReportWorker(args *RegisterArgs, reply *RegisterReply) error {
	panic("implement me")
}

func (m *Master) checkWorkerRunStatusAsync() {
	for !m.Done() {

		time.Sleep(10 * time.Second)
	}
}

func (m *Master) checkWorkerRunStatus() {
	go m.checkWorkerRunStatusAsync()
}

//
// 启动线程监听worker
//
func (m *Master) server() {
	// 注册服务
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{
		fileNames:   files,
		nReduce:     nReduce,
		workerMap:   make(map[int]*workerStatus),
		workerId:    0,
		mapStart:    0,
		reduceStart: 0,
	}

	// Your code here.

	m.server()
	m.checkWorkerRunStatus()

	return &m
}
