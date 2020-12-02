package mr

import (
	"fmt"
	"log"
	"sync"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type IMaster interface {
	RegisterWorker(args *RegisterArgs, reply *RegisterReply) error
	FetchWorker(args *FetchArgs, reply *FetchReply) error
	ReportWorker(args *RegisterArgs, reply *RegisterReply) error
}

type workerStatus struct {
	status int
}

type Master struct {
	// Your definitions here.
	fileNames         []string
	workerMap         map[int]*workerStatus
	nReduce           int
	workerId          int
	outputFileMap     map[int][]string
	mapRequests       int
	reduceFinished    int
	workerIdMutex     sync.RWMutex
	workerMapMutex    sync.RWMutex
	workerOutputMutex sync.RWMutex
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
	if m.mapRequests < len(m.fileNames) {
		m.workerMapMutex.Lock()
		m.workerMap[args.Id].status = 1
		m.mapRequests++
		m.workerMapMutex.Unlock()
		reply.FileName = m.fileNames[m.mapRequests-1]
		reply.Status = 1
		return nil
	}

	return nil
}

func (m *Master) ReportWorker(args *RegisterArgs, reply *RegisterReply) error {
	panic("implement me")
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
		fileNames:      files,
		nReduce:        nReduce,
		workerMap:      make(map[int]*workerStatus),
		workerId:       0,
		outputFileMap:  make(map[int][]string),
		mapRequests:    0,
		reduceFinished: 0,
	}

	// Your code here.

	m.server()
	return &m
}
