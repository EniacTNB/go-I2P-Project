package setting

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize    int
	IdentityKey string

	RootPath string

	LuaPath string
)

func init() {
	fmt.Println("开始执行Init！！！")
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	getRootPath()
	LoadBase()
	LoadServer()
	LoadApp()

}

// 获取根目录
func getRootPath() {
	rootPath, err := os.Getwd()
	fmt.Println(rootPath)
	if err != nil {
		//panic?
		fmt.Println("err at bootstrap")
	}
	// rootPathSlice := strings.Split(rootPath, "/")
	// rootPath = ""

	// for index := 0; index < len(rootPathSlice)-2; index++ {
	// 	fmt.Println(rootPathSlice[index])
	// 	rootPath = fmt.Sprintf("%s/%s", rootPath, rootPathSlice[index])
	// }

	// fmt.Println("RootPath is: %s", rootPath)

	RootPath = rootPath
	LuaPath = Cfg.Section("pat").Key("SCRIPT_PATH").MustString("lua")

}
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	IdentityKey = sec.Key("IDENTITY_KEY").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func InitKafkaConfig() (string, []string, *sarama.Config) {
	//TODO: write into config file
	config := sarama.NewConfig()
	// 配置开启自动提交 offset，这样 samara 库会定时帮我们把最新的 offset 信息提交给 kafka
	config.Consumer.Offsets.AutoCommit.Enable = true              // 开启自动 commit offset
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动 commit时间间隔

	kafkaTopic := "routerinfo"
	kafkaBrokers := []string{"207.148.64.55:9092"}

	return kafkaTopic, kafkaBrokers, config
}
