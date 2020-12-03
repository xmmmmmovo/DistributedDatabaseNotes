package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

// Add your RPC definitions here.

// 注册参数
type RegisterArgs struct {
}

// 注册返回
type RegisterReply struct {
	Id int
}

// 获取任务参数
type FetchArgs struct {
	Id int
}

// 获取任务返回
type FetchReply struct {
	Status    int
	FileNames []string
	NReduce   int
}

// 报告参数
type ReportArgs struct {
}

// 报告返回
type ReportReply struct {
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
