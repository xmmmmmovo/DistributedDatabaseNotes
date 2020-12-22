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
	outputFileMap        [][]string            // 输出文件
	workerMap            map[int]*workerStatus // 工作状态
	nReduce              int                   // reduce数量
	workerId             int                   // id
	taskId               int                   // taskId
	mapFinished          int                   // map结束
	mapStart             int                   // map开始
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
	reply.NReduce = m.nReduce
	m.workerMap[reply.Id] = &workerStatus{}
	return nil
}

func (m *Master) FetchWorker(args *FetchArgs, reply *FetchReply) error {
	m.workerMapReduceMutex.Lock()
	defer m.workerMapReduceMutex.Unlock()

	if m.mapStart < len(m.fileNames) {
		m.workerMap[args.Id].status = MapStatus
		m.workerMap[args.Id].fileIndex = m.mapStart
		m.workerMap[args.Id].fetchWorkTime = time.Now()
		reply.FileNames = []string{m.fileNames[m.mapStart]}
		m.mapStart++
		reply.Status = MapStatus
		reply.TaskId = m.taskId
		m.taskId++
		return nil
	} else if m.mapFinished == len(m.fileNames) && m.reduceStart < m.nReduce {
		m.workerMap[args.Id].status = ReduceStatus
		m.workerMap[args.Id].fileIndex = m.reduceStart
		m.workerMap[args.Id].fetchWorkTime = time.Now()
		reply.FileNames = m.outputFileMap[m.reduceStart]
		m.reduceStart++
		reply.Status = ReduceStatus
		reply.TaskId = m.taskId
		m.taskId++
		return nil
	}
	m.workerMap[args.Id].status = IdleStatus
	m.workerMap[args.Id].fetchWorkTime = time.Now()
	reply.Status = IdleStatus
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
		reduceStart: 0,
	}

	// Your code here.

	m.server()
	m.checkWorkerRunStatus()

	return &m
}
