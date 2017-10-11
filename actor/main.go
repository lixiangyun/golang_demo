package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)
import "flag"
import "strconv"
import "runtime/pprof"

type Msg struct {
	send_fid  int       // 发送的FID编号
	recv_fid  int       // 接收的FID编号
	send_time time.Time // 时间戳
	send_id   int       // 消息ID
}

type FidInfo struct {
	fid      int      // FID编号
	recv_cnt int64    // 接收统计
	send_cnt int64    // 发送统计
	que      chan Msg // 消息接收缓存队列
}

type FidTable struct {
	fidnum  int              // FID数量
	quelen  int              // 队列长度
	oncemsg int              // 单次发送消息的数量
	fidinfo map[int]*FidInfo // FID信息
	waitgo  chan int         // 等待退出的通道
}

// 延时函数
func RunDelay(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
}

// 累计接收的消息总和
func RecvMsgNum(fidtable *FidTable) int64 {
	var totalnum int64 = 0

	for _, v := range fidtable.fidinfo {
		totalnum += v.recv_cnt
	}

	return totalnum
}

// 累计发送的消息总和
func SendMsgNum(fidtable *FidTable) int64 {
	var totalnum int64 = 0

	for _, v := range fidtable.fidinfo {
		totalnum += v.send_cnt
	}

	return totalnum
}

// 申请一个新FID实例
func NewFidTable(fidnum, quelen, oncemsg int) *FidTable {
	var fidtable FidTable

	fidtable.fidnum = fidnum
	fidtable.quelen = quelen
	fidtable.oncemsg = oncemsg

	fidinfo := make(map[int]*FidInfo, fidnum)

	for i := 0; i < fidnum; i++ {
		fidinfo[i] = &FidInfo{i, 0, 0, make(chan Msg, quelen)}
	}

	fidtable.fidinfo = fidinfo

	fidtable.waitgo = make(chan int, fidnum)

	return &fidtable
}

// 根据FID编号获取指定的FID实例
func GetFidInfo(fidtable *FidTable, fid int) *FidInfo {

	v, ok := fidtable.fidinfo[fid]

	if ok {
		return v
	} else {
		return nil
	}
}

// 向指定FID批量发送消息
func SendToMsg(fidtable *FidTable, send_fid, recv_fid, send_id, oncemsg int) {

	fidinfo := GetFidInfo(fidtable, recv_fid)

	if fidinfo == nil {
		fmt.Println("cat not found the fid : ", recv_fid)
		return
	}

	for i := 0; i < oncemsg; i++ {
		var m = Msg{send_fid, recv_fid, time.Now(), send_id}
		fidinfo.que <- m
	}
}

// FID接收处理协程
func RecvToMsg(fidtable *FidTable, fid int) {

	var recv_fid int

	fidinfo := GetFidInfo(fidtable, fid)

	if fidinfo == nil {
		fmt.Println("cat not found the fid : ", fid)
		return
	}

	if fid+1 == fidtable.fidnum {
		recv_fid = 0
	} else {
		recv_fid = fid + 1
	}

	fmt.Println("start recv msg : ", fid)

	for {

		m, ok := <-fidinfo.que
		if ok == false {
			fmt.Println("close fid recv msg.", fid)
			break
		}

		if m.send_id == -1 {
			break
		}

		fidinfo.recv_cnt++

		SendToMsg(fidtable, fid, recv_fid, 0, 1)

		fidinfo.send_cnt++
	}

	fmt.Println("end recv msg : ", fid)

	fidtable.waitgo <- fid
}

// 启动FID接收实例
func StartRecvMsg(fidtable *FidTable) {

	for i := 0; i < fidtable.fidnum; i++ {
		go RecvToMsg(fidtable, i)
	}

	for i := 0; i < fidtable.fidnum; i++ {
		SendToMsg(fidtable, 0, i, 0, fidtable.oncemsg)
	}
}

// 停止FID实例
func EndRecvMsg(fidtable *FidTable) {
	for i := 0; i < fidtable.fidnum; i++ {
		SendToMsg(fidtable, 0, i, -1, 1)
	}

	for i := 0; i < fidtable.fidnum; i++ {
		<-fidtable.waitgo
	}
}

// 启动性能测试
func StatPerfm(fidtable *FidTable, totalsec int, delaysec int) {

	fmt.Println("start stat fid proccess speed...")

	begin_time := time.Now()
	last_time := begin_time

	last_recvmsg := RecvMsgNum(fidtable)

	for {

		RunDelay(delaysec)

		end_time := time.Now()

		run_time := end_time.Sub(begin_time)

		if run_time.Seconds() > float64(totalsec) {
			break
		}

		run_time = end_time.Sub(last_time)

		now_recvmsg := RecvMsgNum(fidtable)

		fmt.Printf("RunTime : %.3f sec ", run_time.Seconds())
		fmt.Printf("Speed : %.3f kps \r\n", float64(now_recvmsg-last_recvmsg)/run_time.Seconds()/1000)

		last_time = end_time
		last_recvmsg = now_recvmsg
	}

	fmt.Println("total time out")
}

func Test2(fidnum, quelen, oncemsg, totalsec, delaysec int) {
	fidtable := NewFidTable(fidnum, quelen, oncemsg)

	StartRecvMsg(fidtable)

	StatPerfm(fidtable, totalsec, delaysec)

	EndRecvMsg(fidtable)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	args := flag.Args()

	if len(args) < 5 {
		fmt.Println("Usage : fidnum quelen oncemsg totalsec delaysec")
		return
	}

	runtime.SetCPUProfileRate(1000)

	fidnum, _ := strconv.Atoi(args[0])
	quelen, _ := strconv.Atoi(args[1])
	oncemsg, _ := strconv.Atoi(args[2])
	totalsec, _ := strconv.Atoi(args[3])
	delaysec, _ := strconv.Atoi(args[4])

	fmt.Println("FidNum   : ", fidnum)
	fmt.Println("QueLen   : ", quelen)
	fmt.Println("OnceMsg  : ", oncemsg)
	fmt.Println("TotalSec : ", totalsec)
	fmt.Println("DelaySec : ", delaysec)

	pfile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Create file failed!")
		return
	}
	pprof.StartCPUProfile(pfile)

	Test2(fidnum, quelen, oncemsg, totalsec, delaysec)

	pprof.StopCPUProfile()
}
