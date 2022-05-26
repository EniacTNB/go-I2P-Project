package i2p

import (
	"fmt"
	"gin-vue/dao/kafka"
	"gin-vue/models"
	"gin-vue/pkg/e"
	"gin-vue/pkg/setting"
	"gin-vue/pkg/util"
	"gin-vue/service"
	"gin-vue/service/docker"
	"gin-vue/viewModels"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDockerStateTest(c *gin.Context) {
	state := docker.CheckDockerStateTestF()
	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": state,
	})
}

func GetDockerLogs(c *gin.Context) {
	docker.GetDockerImageLog()
	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": "test2",
	})

}
func GetNetDbFileList(c *gin.Context) {
	res := docker.GetRouterList()
	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": res,
	})
}

func GetNetDbChangeInDocker(c *gin.Context) {
	fileChangeStr, err := docker.ExecF()
	if err != nil {
		code := e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": "Get Docker Status Error",
		})
	}

	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": fileChangeStr,
	})

}
func GetFileList(c *gin.Context) {
	path := path.Join(setting.RootPath, setting.LuaPath)
	files, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Println(err)
		//TODO: Error Return
		return
	}
	res := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), ".lua") {
			res = append(res, f.Name())
		}
	}
	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": res,
	})

}

// 查看文件内容
func GetFileContent(c *gin.Context) {
	fileName := c.Query("filename")
	pathStr := path.Join(setting.RootPath, setting.LuaPath)
	pathStr = path.Join(pathStr, fileName)
	fmt.Println(pathStr)
	data, err := ioutil.ReadFile(pathStr)
	if err != nil {
		code := e.ERROR_EXIST_TAG
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
		})
	}

	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": string(data),
	})
}

func SaveLuaFile(c *gin.Context) {
	var upFile viewModels.LuaFile
	c.Bind(&upFile)

	fileName := upFile.FileName
	fileContent := upFile.FileContent
	if len(fileName) == 0 {
		timeUnix := time.Now().Unix()
		fileName = fmt.Sprintf("default%d", timeUnix)
	}

	if !strings.Contains(fileName, ".lua") {
		fileName = fmt.Sprintf("%s%s", fileName, ".lua")
	} else {
		if strings.LastIndex(fileName, ".lua") != len(fileName)-4 {
			fileName = fmt.Sprintf("%s%s", fileName, ".lua")
		}
	}

	fmt.Printf("filename is %s \n", fileName)
	fmt.Printf("filecontent is %s \n", fileContent)
	fmt.Printf("Lua File is %+v", upFile)

	pathStr := path.Join(setting.RootPath, setting.LuaPath)
	filePath := pathStr
	pathStr = path.Join(pathStr, fileName)

	isFileExist := util.CheckFileExist(pathStr)

	if isFileExist {
		// 文件已经存在
		c.JSON(http.StatusOK, gin.H{
			"code": e.ERROR,
			"msg":  "filename dumplicated",
		})
		return
	}
	// 可以创建文件
	f, err := os.Create(pathStr) //创建文件
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": e.ERROR,
			"msg":  "create lua file error",
		})
	}
	defer f.Close()
	_, err = io.WriteString(f, fileContent)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": e.ERROR,
			"msg":  "write file failed",
		})
	}
	// start save into mysql
	saveLuaExper2Db(upFile, fileName, filePath)

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
	})

}

func saveLuaExper2Db(upFile viewModels.LuaFile, fileName string, filePath string) {
	fmt.Println("=====Start save function")
	// LuaFileMdb :=
	experMdb := &models.LuaFileMdb{}
	// experMdb
	experMdb.FileName = fileName
	experMdb.FilePath = filePath
	experMdb.ExperPoint = 1
	experMdb.ExperTime = upFile.ExperTime
	experMdb.AddTime = time.Now().UTC()
	experMdb.IsPublic = upFile.IsPublic
	experMdb.Desc = upFile.Desc
	experMdb.User = upFile.User

	if len(upFile.Date1) == 0 || len(upFile.Date2) == 0 {
		experMdb.AutoStartTime = time.Now().UTC()
		experMdb.IsAutoStart = 0
	} else {
		autoStartTime := strings.Split(upFile.Date1, "T")[0] + " " + strings.Split(upFile.Date2, "T")[1]
		autoStartTime = autoStartTime[0 : len(autoStartTime)-5]
		fmt.Println("autoStartTime is " + autoStartTime)
		loc, err := time.LoadLocation("Local")
		if err != nil {
			fmt.Println("Load location error: ", err)
		}
		experMdb.AutoStartTime, err = time.ParseInLocation("2006-01-02 15:04:05", autoStartTime, loc)
		if err != nil {
			fmt.Println("error at parsein location: ", err)
		}
		experMdb.IsAutoStart = 1
	}

	experMdb.StartTime = time.Now().UTC()

	fmt.Printf("%+v\n", experMdb)

	models.InsertExper(experMdb)
	fmt.Println("=====End save function")

}

func KafkaConsumeSample(c *gin.Context) {
	testKafkaMesg := "{\"pubkey\": \"BtLD~gAH1vP2O1I1cZ87MHhrMRVYaPmxB0owlflXTP8=\", \"signkey\": \"~gkz-04bNOWHKT56EihMR0umin2q2z7Vf7WwPgLu9GU=\", \"options\": {\"caps\": \"NfR\", \"netId\": \"2\", \"netdb.knownLeaseSets\": \"37\", \"netdb.knownRouters\": \"3859\", \"router.version\": \"0.9.31\"}, \"addrs\": [{\"cost\": 4, \"transport\": \"SSU\", \"options\": {\"caps\": \"BC\", \"host\": \"114.214.191.255\", \"key\": \"3WxV3BXFYF2vxI~CXKC4WNkEVomB1tntZ5p2lC6kKOo=\", \"mtu\": \"1488\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 6, \"transport\": \"SSU\", \"options\": {\"caps\": \"BC\", \"host\": \"2.219.3.198\", \"key\": \"3WxV3BXFYF2vxI~CXKC4WNkEVomB1tntZ5p2lC6kKOo=\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 9, \"transport\": \"NTCP\", \"options\": {\"host\": \"2a02:c7d:202f:ae00:bdb2:fb70:8aeb:90a\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}, {\"cost\": 10, \"transport\": \"NTCP\", \"options\": {\"host\": \"2.219.3.198\", \"port\": \"26387\"}, \"expire\": 0, \"location\": null}], \"cert\": {\"signature_type\": \"EdDSA_SHA512_Ed25519\", \"crypto_type\": \"ElGamal\"}, \"published\": 1635859281314, \"signature\": \"Pj~JQgamEKDGHHgDqTgaUZxpdZXrrpVeTcFFstJ1CfE=\",\"filename\":\"routerInfo-JaMocLiyYtoeP~jw7EllxiVIM-o3xivw93PLtgkn2Go=.dat\"}"
	kafka.ParseRouterInfo(testKafkaMesg)
}

func GetMetaDataContent(c *gin.Context) {
	pubkey, isExist := c.GetQuery("pubkey")
	if !isExist {
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
	res := models.GetMetaDataContent(pubkey)
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": res,
	})

}

func GetMetaData(c *gin.Context) {
	//TODO: 添加条件搜索
	size, isExist := c.GetQuery("size")
	if !isExist {
		size = "10"
	}

	page, isExist := c.GetQuery("page")
	if !isExist {
		page = "1"
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 10
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	res := models.GetMetaData(sizeInt, pageInt)
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": res,
	})
}

func GetExperList(c *gin.Context) {
	size, isExist := c.GetQuery("size")
	if !isExist {
		size = "10"
	}
	page, isExist := c.GetQuery("page")
	if !isExist {
		page = "1"
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 10
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	res := models.GetExperList(sizeInt, pageInt)
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": res,
	})
}

func RunLuaFile(c *gin.Context) {
	luaFileName, isExist := c.GetQuery("filename")
	if !isExist {
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.ERROR_NOT_EXIST_ARTICLE,
			"msg":  e.GetMsg(e.ERROR_NOT_EXIST_ARTICLE),
			"data": "未找到指定lua文件",
		})
		return
	}

	pathStr := path.Join(setting.RootPath, setting.LuaPath)
	pathStr = path.Join(pathStr, luaFileName)

	isFileExist := util.CheckFileExist(pathStr)

	if !isFileExist {
		// 文件已经存在
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.ERROR_NOT_EXIST_ARTICLE,
			"msg":  e.GetMsg(e.ERROR_NOT_EXIST_ARTICLE),
			"data": "未找到指定lua文件",
		})
		return
	}

	service.ExecuteLuaFile(pathStr)

}
