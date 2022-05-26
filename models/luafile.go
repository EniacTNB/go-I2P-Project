package models

import (
	"fmt"
	"time"
)

// var db *gorm.DB

type LuaFileMdb struct {
	ExperId       int       `json:"-" gorm:"column:exper_id;AUTO_INCREMENT"`
	FileName      string    `json:"file_name" gorm:"column:filename"`
	FilePath      string    `json:"file_path" gorm:"column:filePath"`
	ExperPoint    int       `json:"exper_type" gorm:"column:exper_point"`
	ExperTime     int       `json:"exper_time" gorm:"column:exper_time"`
	AutoStartTime time.Time `json:"auto_start_time" gorm:"column:auto_start_time"`
	IsPublic      bool      `json:"isPublic" gorm:"column:is_public"`
	StartTime     time.Time `json:"start_time" gorm:"column:start_time"`
	AddTime       time.Time `json"add_time" gorm:"column:add_time"`
	Desc          string    `json:"desc" gorm:"column:desc"`
	User          string    `json:"user" gorm:"column:user"`
	Status        int       `json:"status" gorm:"column:status"`
	IsAutoStart   int       `json:"is_auto_start" gorm:"column:is_auto_start"`
}

func (r LuaFileMdb) TableName() string {
	return "i2p_exper"
}

func GetExperList(pageSize int, page int) []LuaFileMdb {
	tempList := make([]LuaFileMdb, 0)
	Db := db
	Db = Db.Limit(pageSize).Offset((page - 1) * pageSize)
	err := Db.Table("i2p_exper").Find(&tempList).Error
	if err != nil {
		fmt.Println("get Exper list error!", err)
		return nil
	}

	return tempList
}

func InsertExper(exper *LuaFileMdb) {
	result := db.Create(exper)
	if result.Error != nil {
		fmt.Println("exper", result.Error)
		return
	}
}
