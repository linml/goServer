// PZ_GameServer project main.go
//
// 品峥麻将 游戏服务器
//
//
//	目录结构:
//		base  基础类
//		game  游戏类
//		msg   消息类(RPC, JSON, Protobuf)
//		net   网络类(Websocket)
//		data  数据类(DBServer->MySql)

package main

import (
	"fmt"
	//"time"
	//"log" Z
	//"net/http"

	"flag"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	//"runtime/pprof"
	"syscall"

	logf "PZ_GameServer/log"
	"PZ_GameServer/server"
	//"PZ_GameServer/server/game"
	//"PZ_GameServer/server/game/room/xiangshan"
	//"PZ_GameServer/server/game/roombase"
)

var (
	version    = "0.0.1"
	buildStamp = "no timestamp set"
	commit     = "no hash set"
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func sendchan(i int) {

	fmt.Println("sendchan ")
	_, isClose := <-out
	if !isClose {
		out <- i
		fmt.Println("out...")
	}
}

var out = make(chan int)

func main() {
	fmt.Printf("Version:%s\n", version)
	fmt.Println("buildstamp is:", buildStamp)
	fmt.Println("commit is:", commit)

	runtime.GOMAXPROCS(runtime.NumCPU())

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	//	arr := []string{"4", "5", "6"}
	//	arr = append(arr[1:])

	//	arr[0] = arr[len(arr)-1]

	//	fmt.Println("0:" + arr[0] + " 2:" + arr[len(arr)-1])
	//	//endChan := make(chan int)
	//	end := false
	//	state := 0
	//	go func() {
	//		for {

	//			if end {
	//				fmt.Println("end it")
	//				break
	//			}
	//			time.Sleep(10)
	//			switch state {
	//			case 1:

	//				fmt.Println("state 1")
	//				state = -1
	//			case 2:
	//				fmt.Println("state 2")
	//				state = -1
	//			}
	//			//			select {
	//			//			case s := <-out:
	//			//				fmt.Println("s = ", s)
	//			//				if s == 1 {
	//			//					fmt.Println("break")
	//			//					break
	//			//				}
	//			//			}
	//			//			break
	//		}
	//		//fmt.Println("close")
	//		//close(out)
	//	}()
	//	for i := 0; i < 3; i++ {
	//		time.Sleep(2 * time.Second)
	//		state = 1
	//		state = 2
	//		state = 1
	//		time.Sleep(2 * time.Second)
	//		state = 2
	//	}

	//	//	sendchan(1)
	//	time.Sleep(2 * time.Second)
	//	end = true

	//	fmt.Println("  end ")

	//	//	sendchan(1)
	//	//	time.Sleep(1000)
	//	//	fmt.Println("end")
	//	//endChan <- 1
	//	//<-endChan
	//	return

	//	logFile, _ := os.OpenFile("../log/server.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	//	syscall.Dup2(int(logFile.Fd()), 1)
	//	syscall.Dup2(int(logFile.Fd()), 2)
	//	go func() {
	//		//log.Println(http.ListenAndServe("116.62.209.46:8081", nil))  //debug/pprof/
	//		http.ListenAndServe("localhost:8081", nil)
	//		//log.Println(http.ListenAndServe("localhost:8081", nil))
	//	}()

	//	f, err := os.OpenFile("./cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer f.Close()
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()

	//WebSocket
	server.Start()

	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-exitChan

	server.Stop()

	// 注意，有时候 defer f.Close()， defer pprof.StopCPUProfile() 会执行不到，这时候我们就会看到 prof 文件是空的， 我们需要在自己代码退出的地方，增加上下面两行，确保写文件内容了。
	//	pprof.StopCPUProfile()
	//	f.Close()

	logf.Debug("Server End ")
}

//CPU信息
//简单的获取方法fmt.Sprintf("Num:%d Arch:%s\n", runtime.NumCPU(), runtime.GOARCH)
//func GetCpuInfo() string {
//	var size uint32 = 128
//	var buffer = make([]uint16, size)
//	var index = uint32(copy(buffer, syscall.StringToUTF16("Num:")) - 1)
//	nums := syscall.StringToUTF16Ptr("NUMBER_OF_PROCESSORS")
//	arch := syscall.StringToUTF16Ptr("PROCESSOR_ARCHITECTURE")
//	r, err := syscall.GetEnvironmentVariable(nums, &buffer[index], size-index)
//	if err != nil {
//		return ""
//	}
//	index += r
//	index += uint32(copy(buffer[index:], syscall.StringToUTF16(" Arch:")) - 1)
//	r, err = syscall.GetEnvironmentVariable(arch, &buffer[index], size-index)
//	if err != nil {
//		return syscall.UTF16ToString(buffer[:index])
//	}
//	index += r
//	return syscall.UTF16ToString(buffer[:index+r])
//}
