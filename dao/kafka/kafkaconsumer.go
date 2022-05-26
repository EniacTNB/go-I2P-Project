package kafka

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gin-vue/models"
	"gin-vue/pkg/setting"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Shopify/sarama"
)

func init() {
	// go RouterInfoLinstenr()
}
func RouterInfoLinstenr() {
	//TODO: 将其放到配置文件中
	topic, brokers, config := setting.InitKafkaConfig()
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		// log.Fatal("NewClient err: ", err)
		//TODO: Log Fail
		fmt.Println("new kafka client error", err)
	}
	defer client.Close()
	// offsetManager 用于管理每个 consumerGroup的 offset
	// 根据 groupID 来区分不同的 consumer，注意: 每次提交的 offset 信息也是和 groupID 关联的
	offsetManager, _ := sarama.NewOffsetManagerFromClient("test-consumer-group", client) // 偏移量管理器
	defer offsetManager.Close()
	// 每个分区的 offset 也是分别管理的，demo 这里使用 0 分区，因为该 topic 只有 1 个分区
	partitionOffsetManager, _ := offsetManager.ManagePartition(topic, 0) // 对应分区的偏移量管理器
	defer partitionOffsetManager.Close()
	defer offsetManager.Commit()
	consumer, _ := sarama.NewConsumerFromClient(client)
	// 根据 kafka 中记录的上次消费的 offset 开始+1的位置接着消费
	nextOffset, _ := partitionOffsetManager.NextOffset() // 取得下一消息的偏移量作为本次消费的起点
	pc, _ := consumer.ConsumePartition(topic, 0, nextOffset)
	defer pc.Close()
	for message := range pc.Messages() {
		value := string(message.Value)
		fmt.Println(fmt.Printf("[Consumer] partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, value))
		// 每次消费后都更新一次 offset,这里更新的只是程序内存中的值，需要 commit 之后才能提交到 kafka
		partitionOffsetManager.MarkOffset(message.Offset+1, "modified metadata") // MarkOffset 更新最后消费的 offset
		go ParseRouterInfo(value)
		// go consumeMessage(value)
	}
}
func consumeMessage(value string) {
	fmt.Println(value)
}
func ParseRouterInfo(kafkaValue string) {
	routerInfo := &models.RouterInfo{}
	err := json.Unmarshal([]byte(kafkaValue), routerInfo)
	if err != nil {
		//Log UnRouter File
		//TODO: Log Unmarshal Error
		fmt.Println("Not Right File", err)
		return
	}

	GetIpInfoInAddrs(routerInfo)

	pubKey := routerInfo.Pubkey
	signKey := routerInfo.Signkey
	if !models.CheckRouterInfoDump(pubKey, signKey) {
		models.InsertRouterMsg(routerInfo)
	}

}
func GetIpInfoInAddrs(routerInfo *models.RouterInfo) {
	addrs := routerInfo.Addrs
	for index := range addrs {
		if addrs[index].Location == nil {
			getIpInfo(addrs[index].Options.Host, routerInfo.RouterFilename)
		}
	}

}

func getIpInfo(ipStr string, filename string) {
	location := QueryIpInfoFromBS(ipStr)
	if location == nil {
		return
	}
	locationMdb := &models.IpAddressInfoAliYunMdb{
		Ip:       location.IP,
		EnShort:  location.Result.EnShort,
		EnName:   location.Result.EnName,
		Nation:   location.Result.Nation,
		Province: location.Result.Province,
		City:     location.Result.City,
		District: location.Result.District,
		Adcode:   location.Result.Adcode,
		Lat:      location.Result.Lat,
		Lng:      location.Result.Lng,
		Filename: filename,
	}
	models.InsertLocationMdb(locationMdb)
}

// 在这里放弃百度IP查询的接口，准确率太低了
func QueryIpInfoFromBS(ipStr string) *models.IpAddressInfoAliYun {
	return QueryIpInfoFromAli(ipStr)
	// url := fmt.Sprintf("http://api.map.baidu.com/location/ip?ak=SxyLgdi40Yj8U0UBhaxMHuoLTaRvm6p6&ip=%s&coor=bd09ll", ipStr)
	// resp, err := http.Get(url)
	// fmt.Println("Url is =====>", url)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if resp.StatusCode == 200 {
	// 	fmt.Println(string(body))
	// 	routerAddressInfo := &models.IpAddressInfoBaidu{}
	// 	err := json.Unmarshal([]byte(body), &routerAddressInfo)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return QueryIpInfoFromAli(ipStr)
	// 	}
	// 	if routerAddressInfo.Status != 0 {
	// 		return QueryIpInfoFromAli(ipStr)

	// 	} else {
	// 		fmt.Println(routerAddressInfo.Content.AddressDetail.City)
	// 		return QueryIpInfoFromAli(ipStr)
	// 		// return routerAddressInfo.Content.AddressDetail.City
	// 	}
	// } else {
	// 	//TODO: 错误日志处理逻辑
	// 	return nil
	// }
}

func QueryIpInfoFromAli(ipStr string) *models.IpAddressInfoAliYun {
	client := &http.Client{}
	url := fmt.Sprintf("https://ips.market.alicloudapi.com/iplocaltion?ip=%s", ipStr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Set("Authorization", "APPCODE 9b41840f989d462da77b3d0307407bd9")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("bodytest is=======>", string(bodyText))
	// fmt.Printf("%s\n", bodyText)
	routerAddressInfo := &models.IpAddressInfoAliYun{}
	err = json.Unmarshal([]byte(bodyText), &routerAddressInfo)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// fmt.Println(routerAddressInfo)
	if routerAddressInfo.Code != 100 {
		return nil
	} else {
		return routerAddressInfo
	}
}

func QueryIpInfoFromNewAli(ipStr string) *models.IpAddressInfoAliYun {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	urlStr := "https://jumipaddup.market.alicloudapi.com"
	queryPath := "/ip/address-query-v2"
	client := &http.Client{Transport: tr}
	appCode := "9b41840f989d462da77b3d0307407bd9"
	// add body
	params := url.Values{}
	params.Add("ip", ipStr)
	reqBody := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", urlStr+queryPath, reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	req.Header.Add("Authorization", "APPCODE "+appCode)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// reqBody, _ := ioutil.ReadAll(req.Body)
	// fmt.Println("%+v", string(reqBody))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	bodyStr := string(body)
	routerNewAddressInfo := &models.IpAddressInfoNewAliyun{}
	err = json.Unmarshal([]byte(bodyStr), &routerNewAddressInfo)
	if err != nil {
		return nil
	}
	if !routerNewAddressInfo.Success {
		return nil
	}

	routerAddressInfo := &models.IpAddressInfoAliYun{}
	routerAddressInfo.IP = ipStr
	routerAddressInfo.Result.Province = routerNewAddressInfo.Data.Province
	// enshort enname
	routerAddressInfo.Result.Nation = routerNewAddressInfo.Data.Nation
	routerAddressInfo.Result.City = routerNewAddressInfo.Data.City
	routerAddressInfo.Result.Adcode = routerNewAddressInfo.Data.Code
	routerAddressInfo.Result.District = routerNewAddressInfo.Data.District
	routerAddressInfo.Result.Lat = routerNewAddressInfo.Data.Latitude
	routerAddressInfo.Result.Lng = routerNewAddressInfo.Data.Longitude

	fmt.Println("%+v", routerAddressInfo)
	// add header

	return routerAddressInfo
}
