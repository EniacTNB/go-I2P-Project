package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"encoding/json"

	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func Dockert() {
	cl, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	images, err := cl.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Println("=======ImageList===========")
	for _, image := range images {
		fmt.Println(image)
	}

	containers, err := cl.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("=======ContainerList===========")
	for _, container := range containers {
		fmt.Println(container)
		fmt.Printf("ContainerName: %s \n", container.Names)
		fmt.Printf("ContainerImage: %s \n", container.Image)
		fmt.Printf("ContainerCommand: %s \n", container.Command)
		fmt.Printf("ContainerCreated: %d \n", container.Created)
		fmt.Printf("ContainerPorts: %+v \n", container.Ports)
		fmt.Printf("ContainerState: %s \n", container.State)
		fmt.Printf("ContainerStatus: %s \n", container.Status)
	}
}

func GetDockerImageLog() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	options := types.ContainerLogsOptions{ShowStdout: true}

	out, err := cli.ContainerLogs(ctx, "453839f6d60b", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)

}

var sts types.Stats

type DockerState struct {
	Name        string //获取容器名称
	ID          string //容器id
	PrivateIP   string //指代机器的IP,不指代特定docker ip
	PublicIP    string //指代公网IP,不指代docker ip
	Logpath     string //获取容器路径
	Size        string //获取日志大小
	RestartTime string //获取容器启动时间
	//FinishedTime string //容器关闭时间，需要记录的是这一次的状态，而不是上一次的状态，所以删除。
	LogTime               string //启动日志收集容器时间本机utc时间-3分钟
	LogSize               string // docker日志文件大小
	LogCount              int    //计数，查过10次，也就是10分钟没人处理吧日志清空。
	Oversize              bool   //日志过大，超过预期
	OneMinute             bool   //运行状态是否小于一分钟
	IsUP                  bool   // 容器是否为up启动状态，down表示已经被删除。
	GetLogcomment         string //显示重启前三分钟日志方法
	Cpupercent            string //cpu占用百分比
	Memlitpercent         string //限制状态下，内存占用百分比
	Memallpercent         string //总内存状态，内存占用百分比
	Memusagesize          string //容器使用内存
	Memlimitsize          string //容器限制使用的内存
	Cputhresholdalter     bool   //cpu是否超过报警阈值
	Memlimithresholdalter bool   //mem是否超过限制使用内存阈值
	Memallthresholdalter  bool   //mem是否超过总使用内存阈值

}

var memtotal uint64
var memtotalsize string
var ncpu int

func CheckDockerStateTestF() *DockerState {
	fmt.Println("in checkDockerStateTestF")
	nproc, memtotal1, err := GetCPUAndMemPercent()
	ncpu = nproc
	memtotal = memtotal1
	if err != nil {
		fmt.Println(err)
	}
	memtotalsize = GetMemSize(memtotal)
	state := GetDockerState()

	// c := cron.New()
	// //定时任务一分钟重启一次
	// c.AddFunc("@every 1m", GetDockerState)
	// c.Start()
	// select {}
	return state
}

func GetDockerState() *DockerState {
	//获取容器状态，参考docker官方教程
	fmt.Println("in GetDockerState")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	ds := DockerState{}
	for _, container := range containers {
		//遍历容器
		//获取容器信息

		fmt.Println(strings.Trim(container.Names[0], "/"))
		if !strings.Contains(strings.Trim(container.Names[0], "/"), "i2p") {
			continue
		} else {
			fmt.Println("yesyesyes")
		}
		ds.ID = container.ID[:12]
		ds.Name = strings.Trim(container.Names[0], "/")

		cpupercent, memusage, memlimit, err := GetDockersts(ds.ID)
		if err != nil {
			fmt.Println(err)
		}

		//设置cpu超过0.75报警，这里是总cpu，不是单个cpu
		//设置mem内存与限制内存的比例超过0.8报警。
		//设置mem内存与限制内存的比例超过0.6报警
		cpupercentsize, memlitpercent, memallpercent, memusagesize, memlimitsize, cputhresholdalter, memlimithresholdalter, memallthresholdalter, err := GetDockerAlter(cpupercent, memusage, memlimit, memtotal, ncpu, 0.75, 0.75, 0.6)
		if err != nil {
			fmt.Println(err)
		}
		if cputhresholdalter {
			ds.Cputhresholdalter = true
		}
		if memlimithresholdalter {
			ds.Memlimithresholdalter = true
		}
		if memallthresholdalter {
			ds.Memallthresholdalter = true
		}
		ds.Memlitpercent = memlitpercent
		ds.Cpupercent = cpupercentsize
		ds.Memallpercent = memallpercent
		ds.Memusagesize = memusagesize
		ds.Memlimitsize = memlimitsize
		// ds.memtotalsize = memtotalsize

		url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=dxxxxxx"
		context := "容器名称: " + ds.Name + "\ncpu利用率： " + cpupercentsize + "\n限制内存利用率： " + memlitpercent + "\n正常内存利用率：" + memallpercent + "\n已用内存：" + memusagesize + "\n总限制内存：" + memlimitsize + "\n总内存：" + memtotalsize + "\ncpu是否超用：" + fmt.Sprint(cputhresholdalter) + "\n内存是否超用：" + fmt.Sprint(memlimithresholdalter, memallthresholdalter)
		fmt.Println("context")
		SendMessage(url, context)
		fmt.Printf("%+v", ds)
	}
	return &ds

}

func GetMemSize(u uint64) (size string) {
	//为了简单计算，小于1G的都设置为1G
	if u < 1024 {
		size = fmt.Sprintf("%.2fB", float64(u))
		return
	} else if float64(u) < 1024*1024 {
		size = fmt.Sprintf("%.2fKB", float64(u)/float64(1024))
		return
	} else if float64(u) < 1024*1024*1024 {
		size = fmt.Sprintf("%.2fMB", float64(u)/float64(1024*1024))
		return
	} else if float64(u) < 1024*1024*1024*1024 {
		size = fmt.Sprintf("%.2fGB", float64(u)/float64(1024*1024*1024))
		return
	} else {
		size = fmt.Sprintf("%.2fTB", float64(u)/float64(1024*1024*1024))
		return
	}

}

//本来想从client.info获取内存和，我发现这个内存在window上是限制内存，不是真是内存。
//所以还是通过外部工具获取
func GetCPUAndMemPercent() (ncpu int, memtotal uint64, err error) {
	ncpu, err = cpu.Counts(true)
	if err != nil {
		return 0, 0, err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return ncpu, memInfo.Total, nil
}

//cpu数据取值很复杂，建议使用exec直接调用docker stats命令读取(但是也是肉眼可见，它也是执行了1s多，应该也是抛弃了第一次的数据)
//后来考虑修改client.ContainerStatsoneshot中的函数中的query.Set("one-shot", "1") 将这个值更改为2试下，发现他的函数没有暴露。
//再后来就是通过client.ContainerStatsoneshot中获取的precpustatus是0，所以更改为流式获取，获取第二次的数据
//时间片1s=1000000000ns
//这里使用的是时间片的概念。源码说明unix系统就是ns为单位，windows是以100ns为单位
//大致概念
//unix 1s有1000000000时间片，即为1s总共可以执行那么多动作，谁拿的时间片多，谁占用时间片大
//windows 1s有10000000时间片，即为1s总共可以执行那么多动作，谁拿的时间片多，谁占用时间片大
//所以计算cpu比例如下
//cpu使用率获取思路
// cpuusage = (cpustats-precpustats)/1000000000(unix)
// cpuusage = (cpustats-precpustats)/10000000(windows)

func GetDockerAlter(cpupercent, memusage, memlimit, memtotal uint64, ncpu int, cputhreshold, memlitmitthreshold, memallthreshold float64) (cpupercentsize, memlitpercent, memallpercent, memusagesize, memlimitsize string, cputhresholdalter, memthresholdalter, memallthresholdalter bool, err error) {
	//本地考虑的是将windows和linux区分，再一次思考，windows用的是wsl的。底层运行的是linux内核。
	//sysType := runtime.GOOS
	var cpufloat float64
	cpufloat = float64(cpupercent) / 1000000000

	cpupercentsize = fmt.Sprintf("%.2f%%", cpufloat*100)
	if float64(cpufloat)/float64(ncpu) > cputhreshold {
		cputhresholdalter = true
	}
	memusagesize = GetMemSize(memusage)
	memlimitsize = GetMemSize(memlimit)
	if float64(memusage)/float64(memlimit) > memlitmitthreshold {
		memlitpercent = fmt.Sprintf("%.2f%%", float64(memusage)/float64(memlimit)*100)
		memthresholdalter = true
	}

	if float64(memusage)/float64(memtotal) > memallthreshold {
		memallpercent = fmt.Sprintf("%.2f%%", float64(memusage)/float64(memtotal)*100)
		memallthresholdalter = true
	}
	return cpupercentsize, memlitpercent, memallpercent, memusagesize, memlimitsize, cputhresholdalter, memthresholdalter, memallthresholdalter, nil

}
func GetDockersts(id string) (cpupercent, memusage, memlimit uint64, err error) {
	//获取容器状态，参考docker官方教程
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return 0, 0, 0, err
	}
	stats, err := cli.ContainerStats(ctx, id, true)
	if err != nil {
		return 0, 0, 0, err
	}
	defer stats.Body.Close()

	//获取第二次的数据进行解析
	//可能存在[]byte字节过少的情况
	var stsbody []byte
	buf := make([]byte, 4096)

	//只是为了获取第三次数据,第二次的数据有第一次的cpu信息。
	//这里存在一个问题，取2次的数据要2秒。所以使用这个函数尽量使用go并发处理。
	i := 0
	for {
		n, err := stats.Body.Read(buf)
		if err != nil {
			break
		}
		if i == 2 {
			stsbody = make([]byte, n)
			copy(stsbody, buf[:n])
			break
		}
		i++
	}

	//参考官方文档，请求访问的是/sys 发现获取的是docker.Stats类型
	err = json.Unmarshal(stsbody, &sts)
	if err != nil {
		return 0, 0, 0, err
	}
	cpupercent = sts.CPUStats.CPUUsage.TotalUsage - sts.PreCPUStats.CPUUsage.TotalUsage
	memusage = sts.MemoryStats.Usage
	memlimit = sts.MemoryStats.Limit

	//这里是为了测试精度，发现使用先转换为float精度更准确一点
	//fi := decimal.NewFromInt(int64(sts.MemoryStats.Usage)*100)
	//fy := decimal.NewFromInt(int64(sts.MemoryStats.Limit))
	//sub := fi.Div(fy)
	//
	//fmt.Printf("%.2f%%\n",float64(sts.MemoryStats.Usage/sts.MemoryStats.Limit)*100)
	//fmt.Printf("%.2f%%\n",float64(sts.MemoryStats.Usage)/float64(sts.MemoryStats.Limit)*100)
	//fmt.Println(sub)
	//fmt.Println(GetMemSize(sts.MemoryStats.Usage))
	//fmt.Println(GetMemSize(sts.MemoryStats.Limit))

	return cpupercent, memusage, memlimit, nil
}

//企业微信webhook发送
func SendMessage(url, msg string) {
	fmt.Println("url: " + url + " msg: " + msg)
}

//  "'ls", "-l", "/home/i2pd/data/netDb/'"
func GetNetDbFileChange() (string, error) {
	out2, err2 := exec.Command("/bin/bash", "-c", "ps").Output()
	if err2 != nil {
		fmt.Println("err3!")
		fmt.Println(err2)
	}
	fmt.Printf("%s", out2)
	testCmd := "docker exec -it i2pd /bin/sh -c 'ps'"
	out1, err1 := exec.Command("/bin/bash", "-c", testCmd).Output()
	if err1 != nil {
		fmt.Println("err1!")
		fmt.Println(err1)
	}
	fmt.Printf("%s", out1)
	baseCmdStr := "docker exec -it i2pd /bin/sh -c '"
	netDbPath := "/home/i2pd/data/netDb/"
	relativeCmdStr := "ls -l "

	// routerInfoMap := make(map[string]string)
	// r- r~ r0 -9 r a-z
	cmdStr := baseCmdStr + relativeCmdStr + netDbPath + "'"

	out, err := exec.Command("/bin/bash", "-c", cmdStr).Output()
	if err != nil {
		fmt.Println("err2!")
		fmt.Println(err)
	}

	// for()
	outStr := fmt.Sprintf("%s", out)
	return outStr, nil
}

var routerInfoMap map[string]string
var resBaseUpdateStr = " has been updated"
var resBaseAddStr = " has been added"
var resBaseDeleteStr = " has been deleted"

func GetRouterList() map[string][]string {
	relativeNetDbPath := "/home/i2pd/data/netDb/r"
	command := []string{"ls", "-l"}
	command = append(command, relativeNetDbPath+"-")
	// resStrList := []string{}
	routerInfoString, err := Exec(command)
	if err != nil {
		fmt.Println(err)
	}
	if routerInfoString == nil {
		fmt.Println("routerInfoString is Nil!")
	}

	tempList := GetRouterInfoList(routerInfoString)
	res := make(map[string][]string)

	res["r-"] = tempList

	command = command[:len(command)-1]
	command = append(command, relativeNetDbPath+"~")
	routerInfoString, err = Exec(command)
	if err != nil {
		fmt.Println(err)
	}
	if routerInfoString == nil {
		fmt.Println("routerInfoString is Nil!")
	}
	tempList = GetRouterInfoList(routerInfoString)
	res["r~"] = tempList
	for i := '0'; i <= '9'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		if routerInfoString == nil {
			fmt.Println("routerInfoString is Nil!")
		}
		tempList = GetRouterInfoList(routerInfoString)
		res["r"+string(i)] = tempList

	}
	for i := 'a'; i <= 'z'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		if routerInfoString == nil {
			fmt.Println("routerInfoString is Nil!")
		}
		tempList = GetRouterInfoList(routerInfoString)
		res["r"+string(i)] = tempList
	}

	for i := 'A'; i <= 'Z'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		if routerInfoString == nil {
			fmt.Println("routerInfoString is Nil!")
		}
		tempList = GetRouterInfoList(routerInfoString)
		res["r"+string(i)] = tempList
	}

	return res
}

func GetRouterInfoList(routerInfoString []string) []string {
	resList := []string{}
	for _, value := range routerInfoString {
		if strings.Contains(value, "routerInfo") {
			tempStringList := strings.Fields(value)
			fileName := tempStringList[8]
			if strings.Contains(fileName, "router") {
				resList = append(resList, fileName)
			}
		}
	}
	return resList
}
func ExecF() ([]string, error) {
	if routerInfoMap == nil {
		routerInfoMap = make(map[string]string)
	}
	tempRouterInfoMap := make(map[string]string)
	// fmt.Println(tempRouterInfoMap)
	relativeNetDbPath := "/home/i2pd/data/netDb/r"
	command := []string{"ls", "-l"}
	command = append(command, relativeNetDbPath+"-")
	resStrList := []string{}
	routerInfoString, err := Exec(command)
	if err != nil {
		fmt.Println(err)
	}
	if routerInfoString == nil {
		fmt.Println("routerInfoString is Nil!")
	}
	// fmt.Println(routerInfoString)
	tempRouterInfoMap = GetTempRouterInfoMap(tempRouterInfoMap, routerInfoString)
	command = command[:len(command)-1]
	command = append(command, relativeNetDbPath+"~")
	routerInfoString, err = Exec(command)
	if err != nil {
		fmt.Println(err)
	}
	if routerInfoString == nil {
		fmt.Println("routerInfoString is Nil!")
	}
	tempRouterInfoMap = GetTempRouterInfoMap(tempRouterInfoMap, routerInfoString)
	for i := '0'; i <= '9'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		if routerInfoString == nil {
			fmt.Println("routerInfoString is Nil!")
		}
		tempRouterInfoMap = GetTempRouterInfoMap(tempRouterInfoMap, routerInfoString)
	}

	for i := 'a'; i <= 'z'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		if routerInfoString == nil {
			fmt.Println("routerInfoString is Nil!")
		}
		tempRouterInfoMap = GetTempRouterInfoMap(tempRouterInfoMap, routerInfoString)
	}

	for i := 'A'; i <= 'Z'; i++ {
		command = command[:len(command)-1]
		command = append(command, relativeNetDbPath+string(i))
		routerInfoString, err = Exec(command)
		if err != nil {
			fmt.Println(err)
		}
		// if routerInfoString == nil {
		// 	fmt.Println("routerInfoString is Nil!")
		// }
		tempRouterInfoMap = GetTempRouterInfoMap(tempRouterInfoMap, routerInfoString)
	}
	resStrList = CheckMapDiff(tempRouterInfoMap)

	return resStrList, nil
}

func GetTempRouterInfoMap(tempRouterInfoMap map[string]string, routerInfoString []string) map[string]string {
	for _, value := range routerInfoString {
		if strings.Contains(value, "routerInfo") {
			tempStringList := strings.Fields(value)

			// 5 6 7 8
			month := tempStringList[5]
			day := tempStringList[6]
			hour := tempStringList[7]
			fileName := tempStringList[8]

			time := month + day + hour

			tempRouterInfoMap[fileName] = time
		}
	}
	return tempRouterInfoMap
}
func CheckMapDiff(tempRouterInfoMap map[string]string) []string {
	resStrList := []string{}
	for fileName, time := range tempRouterInfoMap {
		timePre, ok := routerInfoMap[fileName]
		if ok {
			if timePre == time {
				continue
			} else {
				tempRouterInfoMap[fileName] = time
				resStrList = append(resStrList, fileName+resBaseUpdateStr)
			}
		} else {
			tempRouterInfoMap[fileName] = time
			resStrList = append(resStrList, fileName+resBaseAddStr)
		}
	}

	for fileName, _ := range routerInfoMap {
		_, ok := tempRouterInfoMap[fileName]
		if !ok {
			resStrList = append(resStrList, fileName+resBaseDeleteStr)
		}
	}

	routerInfoMap = tempRouterInfoMap
	return resStrList
}

func Exec(command []string) ([]string, error) {
	ctx := context.Background()
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	defer docker.Close()
	containerID := "453839f6d60b"

	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          command,
	}

	execId, err := docker.ContainerExecCreate(ctx, containerID, config)

	respId, err := docker.ContainerExecAttach(context.Background(), execId.ID, types.ExecStartCheck{})
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(respId.Reader)
	resStr := []string{}
	for scanner.Scan() {
		tempStr := scanner.Text()
		// fmt.Println(tempStr)
		resStr = append(resStr, tempStr)
	}
	return resStr, nil
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

// func InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
// 	var execResult ExecResult
// 	docker, err := client.NewEnvClient()
// 	if err != nil {
// 		return execResult, err
// 	}
// 	defer docker.Close()

// 	resp, err := docker.ContainerExecAttach(ctx, id, types.ExecConfig{})
// 	if err != nil {
// 		return execResult, err
// 	}
// 	defer resp.Close()

// 	// read the output
// 	var outBuf, errBuf bytes.Buffer
// 	outputDone := make(chan error)

// 	go func() {
// 		// StdCopy demultiplexes the stream into two buffers
// 		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
// 		outputDone <- err
// 	}()

// 	select {
// 	case err := <-outputDone:
// 		if err != nil {
// 			return execResult, err
// 		}
// 		break

// 	case <-ctx.Done():
// 		return execResult, ctx.Err()
// 	}

// 	stdout, err := ioutil.ReadAll(&outBuf)
// 	if err != nil {
// 		return execResult, err
// 	}
// 	stderr, err := ioutil.ReadAll(&errBuf)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	res, err := docker.ContainerExecInspect(ctx, id)
// 	if err != nil {
// 		return execResult, err
// 	}

// 	execResult.ExitCode = res.ExitCode
// 	execResult.StdOut = string(stdout)
// 	execResult.StdErr = string(stderr)
// 	return execResult, nil
// }
