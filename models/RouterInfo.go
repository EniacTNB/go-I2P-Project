package models

import (
	"fmt"
	"strconv"
	"strings"
)

type RouterInfo struct {
	RouterId       int     `json:"-"`
	Pubkey         string  `json:"pubkey,omitempty" `
	Signkey        string  `json:"signkey,omitempty"`
	Options        *Option `json:"options,omitempty"`
	Addrs          []*Addr `json:"addrs,omitempty"`
	Cert           *Cert   `json:"cert,omitempty"`
	Published      int64   `json:"published,omitempty"`
	Signature      string  `json:"signature,omitempty"`
	RouterFilename string  `json:"filename,omitempty"`
}

type RouterInfoMdb struct {
	RouterId       int    `json:"-" gorm:"column:router_id;AUTO_INCREMENT"`
	Pubkey         string `json:"pubkey,omitempty" gorm:"column:pubkey"`
	Signkey        string `json:"signkey,omitempty" gorm:"column:signkey"`
	OptionId       int    `json:"options,omitempty" gorm:"column:options_id"`
	AddrId         string `json:"addrs,omitempty" gorm:"column:addrs_id"`
	CertId         int    `json:"cert,omitempty" gorm:"column:cert_id"`
	Published      int64  `json:"published,omitempty" gorm:"column:published"`
	Signature      string `json:"signature,omitempty" gorm:"column:signature"`
	RouterFilename string `json:"-" gorm:"column:router_filename"`
}

type Option struct {
	OptionId            int    `json:"-" gorm:"column:option_id;AUTO_INCREMENT"`
	Caps                string `json:"caps,omitempty" gorm:"column:caps"`
	NetID               string `json:"netId,omitempty" gorm:"column:net_id"`
	NetdbKnownLeaseSets string `json:"netdb.knownLeaseSets,omitempty" gorm:"column:netdb_known_lease_set"`
	NetdbKnownRouters   string `json:"netdb.knownRouters,omitempty" gorm:"column:netdb_known_routers"`
	RouterVersion       string `json:"router.version,omitempty" gorm:"column:router_version"`
}
type AddrOption struct {
	OptionId int    `json:"-" gorm:"column:option_id;AUTO_INCREMENT"`
	Host     string `json:"host,omitempty" gorm:"column:host"`
	I        string `json:"i,omitempty" gorm:"column:i"`
	Port     string `json:"port,omitempty" gorm:"column:port"`
	S        string `json:"s,omitempty" gorm:"column:s"`
	V        string `json:"v,omitempty" gorm:"column:v"`
	Caps     string `json:"caps,omitempty" gorm:"column:caps"`
	Key      string `json:"key,omitempty" gorm:"column:ikey"`
	Mtu      string `json:"mtu,omitempty" gorm:"column:imtu"`
}

type Addr struct {
	Cost      int         `json:"cost,omitempty"`
	Transport string      `json:"transport,omitempty"`
	Options   *AddrOption `json:"options,omitempty"`
	Expire    int         `json:"expire,omitempty"`
	Location  interface{} `json:"location,omitempty"`
}

type AddrMdb struct {
	AddrId    int         `gorm:"column:addr_id;AUTO_INCREMENT"`
	Cost      int         `json:"cost,omitempty" gorm:"column:cost"`
	Transport string      `json:"transport,omitempty" gorm:"column:transport"`
	OptionId  int         `json:"options,omitempty" gorm:"column:option_id"`
	Expire    int         `json:"expire,omitempty" gorm:"column:expire"`
	Location  interface{} `json:"location,omitempty" gorm:"column:location"`
}
type Cert struct {
	CertId        int    `json:"-" gorm:"column:cert_id;AUTO_INCREMENT"`
	SignatureType string `json:"signature_type,omitempty" gorm:"column:signature_type"`
	CryptoType    string `json:"crypto_type,omitempty" gorm:"cloumn:crypty_type"`
}

type MetaData struct {
	Pubkey         string          `json:"pubkey,omitempty" gorm:"column:pubkey"`
	Signkey        string          `json:"signkey,omitempty" gorm:"column:signkey"`
	Option         *Option         `json:"options,omitempty" gorm:"column:options_id"`
	Addrs          []*AddrMetaData `json:"addrs,omitempty" gorm:"column:addrs_id"`
	Cert           *Cert           `json:"cert,omitempty" gorm:"column:cert_id"`
	Published      int64           `json:"published,omitempty" gorm:"column:published"`
	Signature      string          `json:"signature,omitempty" gorm:"column:signature"`
	RouterFilename string          `json:"-" gorm:"column:router_filename"`
}
type AddrMetaData struct {
	Cost      int         `json:"cost,omitempty"`
	Transport string      `json:"transport,omitempty"`
	Options   *AddrOption `json:"options,omitempty"`
	Expire    int         `json:"expire,omitempty"`
	Location  interface{} `json:"location,omitempty"`
}

func (r AddrMdb) TableName() string {
	return "i2p_addrs"
}
func (r AddrOption) TableName() string {
	return "i2p_addr_options"
}
func (r Option) TableName() string {
	return "i2p_options"
}
func (r RouterInfoMdb) TableName() string {
	return "i2p_routerinfo"
}
func (cert Cert) TableName() string {
	return "i2p_cert"
}

func GetMetaData(pageSize int, page int) []*MetaData {
	tempMeta := make([]RouterInfoMdb, 0)
	Db := db
	Db = Db.Limit(pageSize).Offset((page - 1) * pageSize)
	err := Db.Table("i2p_routerinfo").Find(&tempMeta).Error
	if err != nil {
		fmt.Println("get meta data error", err)
		return nil
	}
	MetaDatas := make([]*MetaData, 0)

	addrMetaDatas := []*AddrMetaData{}
	for _, meta := range tempMeta {
		addrId := meta.AddrId
		addrIdArr := strings.Split(addrId, "_")
		metaData := &MetaData{}
		// Db := db
		for _, addrIdStr := range addrIdArr {
			var tempAddr AddrMdb
			addrIdInt, err := strconv.Atoi(addrIdStr)

			if err != nil {
				continue
			}
			// Db := db
			err = db.Table("i2p_addrs").Find(&tempAddr, "addr_id=?", addrIdInt).Error
			if err != nil {
				continue
			}
			optionId := tempAddr.OptionId
			addrMetaData := &AddrMetaData{
				Cost:      tempAddr.Cost,
				Transport: tempAddr.Transport,

				Expire:   tempAddr.Expire,
				Location: tempAddr.Location,
			}
			var tempAddrOption AddrOption

			err = db.Table("i2p_addr_options").Find(&tempAddrOption, "option_id=?", optionId).Error
			if err != nil {
				continue
			}

			addrMetaData.Options = &tempAddrOption
			addrMetaDatas = append(addrMetaDatas, addrMetaData)
		}
		metaData.Addrs = addrMetaDatas

		certId := meta.CertId
		var tempCert Cert
		err = db.Table("i2p_cert").Find(&tempCert, "cert_id=?", certId).Error
		if err == nil {
			metaData.Cert = &tempCert
		}

		optionsId := meta.OptionId
		var tempOption Option
		err = db.Table("i2p_options").Find(&tempOption, "option_id=?", optionsId).Error
		if err == nil {
			metaData.Option = &tempOption
		}

		metaData.Pubkey = meta.Pubkey
		metaData.Signkey = meta.Signkey
		metaData.Published = meta.Published
		metaData.Signature = meta.Signature
		metaData.RouterFilename = meta.RouterFilename

		MetaDatas = append(MetaDatas, metaData)

	}
	return MetaDatas
}
func CheckRouterInfoDump(pubKey string, signKey string) bool {
	var numRecord int

	db.Table("i2p_routerinfo").Where("pubKey=? AND signKey=?", pubKey, signKey).Count(&numRecord)
	if numRecord > 0 {
		return true
	} else {
		return false
	}
}

func GetMetaDataContent(pubkey string) interface{} {
	metaData := &MetaData{}
	var tempRouterInfo RouterInfoMdb
	db.Table("i2p_routerinfo").Find(&tempRouterInfo, "pubkey=?", pubkey)

	metaData.Pubkey = tempRouterInfo.Pubkey
	metaData.Signkey = tempRouterInfo.Signkey
	metaData.Published = tempRouterInfo.Published
	metaData.Signature = tempRouterInfo.Signature
	metaData.RouterFilename = tempRouterInfo.RouterFilename

	addrId := tempRouterInfo.AddrId
	addrIdArr := strings.Split(addrId, "_")
	addrMetaDatas := make([]*AddrMetaData, 0)
	for _, addrIdStr := range addrIdArr {
		addrIdInt, err := strconv.Atoi(addrIdStr)
		if err != nil {
			continue
		}

		var tempAddr AddrMdb
		err = db.Table("i2p_addrs").Find(&tempAddr, "addr_id=?", addrIdInt).Error
		if err != nil {
			continue
		}
		optionId := tempAddr.OptionId
		addrMetaData := &AddrMetaData{
			Cost:      tempAddr.Cost,
			Transport: tempAddr.Transport,

			Expire:   tempAddr.Expire,
			Location: tempAddr.Location,
		}
		var tempAddrOption AddrOption

		err = db.Table("i2p_addr_options").Find(&tempAddrOption, "option_id=?", optionId).Error
		if err != nil {
			continue
		}

		addrMetaData.Options = &tempAddrOption
		addrMetaDatas = append(addrMetaDatas, addrMetaData)
	}
	metaData.Addrs = addrMetaDatas

	return nil

}

func InsertRouterMsg(routerInfo *RouterInfo) {

	cert := routerInfo.Cert
	result := db.Create(cert)
	if result.Error != nil {
		fmt.Println("cert", result.Error)
		return
	}
	option := routerInfo.Options
	result = db.Create(option)
	if result.Error != nil {
		fmt.Println("options", result.Error)
		return
	}
	cert_id := cert.CertId
	option_id := option.OptionId
	fmt.Println("OpitonId is: ", option_id)

	addrs := routerInfo.Addrs
	// addrsMdb := make([]*DataStruct.AddrMdb, 0)
	addrsId := ""
	for index := range addrs {
		addrOption := addrs[index].Options
		result := db.Create(addrOption)
		if result.Error != nil {
			fmt.Println("addrOption", result.Error)
			return
		}
		addrOptionId := addrOption.OptionId

		addrMdb := &AddrMdb{
			Cost:      addrs[index].Cost,
			Transport: addrs[index].Transport,
			OptionId:  addrOptionId,
			Expire:    addrs[index].Expire,
		}
		result = db.Create(addrMdb)
		if result.Error != nil {
			fmt.Println("arrMdb", result.Error)
			return
		}
		addrId := addrMdb.AddrId
		fmt.Println("addrId: ", addrId)
		addrsId = addrsId + "_" + strconv.Itoa(addrId)
	}
	if len(addrsId) != 0 {
		addrsId = addrsId[1:]
	}

	routerInfoMdb := &RouterInfoMdb{}

	routerInfoMdb.CertId = cert_id
	routerInfoMdb.OptionId = option_id
	routerInfoMdb.Pubkey = routerInfo.Pubkey
	routerInfoMdb.Signature = routerInfo.Signature
	routerInfoMdb.Signkey = routerInfo.Signkey
	routerInfoMdb.Published = routerInfo.Published
	routerInfoMdb.AddrId = addrsId
	routerInfoMdb.RouterFilename = routerInfo.RouterFilename

	fmt.Println(routerInfoMdb.Published)
	result = db.Create(routerInfoMdb)

	if result.Error != nil {
		fmt.Println("routerINfoMdb", result.Error)
		return
	}

	fmt.Println("Insert RouterInfo Success!")

}
