package models

import "fmt"

type RouterAddressInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

type IpAddressInfoBaidu struct {
	Address string `json:"address"`
	Content struct {
		Address       string `json:"address"`
		AddressDetail struct {
			Adcode       string `json:"adcode"`
			City         string `json:"city"`
			CityCode     string `json:"city_code"`
			District     string `json:"district"`
			Province     string `json:"province"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"address_detail"`
		Point struct {
			X string `json:"x"`
			Y string `json:"y"`
		} `json:"point"`
	} `json:"content"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type IpAddressInfoAliYun struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	IP      string `json:"ip"`
	Result  struct {
		EnShort  string  `json:"en_short"`
		EnName   string  `json:"en_name"`
		Nation   string  `json:"nation"`
		Province string  `json:"province"`
		City     string  `json:"city"`
		District string  `json:"district"`
		Adcode   int     `json:"adcode"`
		Lat      float64 `json:"lat"`
		Lng      float64 `json:"lng"`
	} `json:"result"`
}
type IpAddressInfoAliYunMdb struct {
	// Id int `gorm:"column:"`
	Ip       string  `gorm:"column:ip"`
	EnShort  string  `gorm:"column:en_short"`
	EnName   string  `gorm:"column:en_name"`
	Nation   string  `gorm:"column:nation"`
	Province string  `gorm:"column:province"`
	City     string  `gorm:"column:city"`
	District string  `gorm:"column:district"`
	Adcode   int     `gorm:"column:adcode"`
	Lat      float64 `gorm:"column:lat"`
	Lng      float64 `gorm:"column:lng"`
	Filename string  `gorm:"column:filename"`
}

func (r IpAddressInfoAliYunMdb) TableName() string {
	return "i2p_ip_info"
}

type IpCityInfo struct {
	Country string
	City    string
	Num     int
	Lat     string
	Lng     string
}

type EnCnNameDs struct {
	Name   string `json:"name"`
	CnName string `json:"cn_name"`
}

type IpAddressInfoNewAliyun struct {
	Data struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Nation    string  `json:"nation"`
		Province  string  `json:"province"`
		City      string  `json:"city"`
		District  string  `json:"district"`
		Code      int     `json:"code"`
	} `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	TaskNo  string `json:"taskNo"`
}

func InsertLocationMdb(locationMdb *IpAddressInfoAliYunMdb) {
	// if global.Db == nil {
	// 	global.Db = ConnectMdb()
	// 	if global.Db == nil {
	// 		return
	// 	}
	// }
	result := db.Create(locationMdb)
	if result.Error != nil {
		fmt.Println("locationMdb Error")
		return
	}
}
