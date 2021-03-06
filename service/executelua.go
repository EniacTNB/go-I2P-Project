package service

import (
	"fmt"

	"gin-vue/service/docker"

	lua "github.com/yuin/gopher-lua"
)

func ExecuteLuaFile(luaFilePath string) {
	fmt.Println("Start Run Experiment in Service Package")
	luaFilePath = "/Users/eniac/Desktop/server/gin-vue/lua/test3.lua"

	L := lua.NewState()
	defer L.Close()
	// 加载fib.lua
	if err := L.DoFile(luaFilePath); err != nil {
		panic(err)
	}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("fib"), // 获取fib函数引用
		NRet:    1,                  // 指定返回值数量
		Protect: true,               // 如果出现异常，是panic还是返回err
	}, lua.LNumber(10)) // 传递输入参数n=10
	if err != nil {
		panic(err)
	}
	// 获取返回结果
	ret := L.Get(-1)
	// 从堆栈中扔掉返回结果
	L.Pop(1)
	// 打印结果
	res, ok := ret.(lua.LNumber)
	if ok {
		fmt.Println(int(res))
	} else {
		fmt.Println("unexpected result")
	}

	docker.Dockert()



}
