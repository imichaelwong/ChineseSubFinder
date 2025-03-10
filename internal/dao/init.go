package dao

import (
	"errors"
	"fmt"
	"github.com/allanpk716/ChineseSubFinder/internal/models"
	"github.com/allanpk716/ChineseSubFinder/internal/pkg/log_helper"
	"github.com/allanpk716/ChineseSubFinder/internal/pkg/my_util"
	"github.com/allanpk716/ChineseSubFinder/internal/pkg/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// GetDb 获取数据库实例
func GetDb() *gorm.DB {
	if db == nil {
		once.Do(func() {
			err := InitDb()
			if err != nil {
				log_helper.GetLogger().Errorln("dao.InitDb()", err)
				panic(err)
			}
		})
	}
	return db
}

// DeleteDbFile 删除 Db 文件
func DeleteDbFile() error {

	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Close()
		if err != nil {
			return err
		}
	}

	// 这里需要考虑是 Windows 的时候就是在本程序的允许目录下新建数据库即可
	// 如果是 Linux 则在 /config 目录下
	nowDbFileName := getDbName()

	if my_util.IsFile(nowDbFileName) == true {
		return os.Remove(nowDbFileName)
	}
	return nil
}

// InitDb 初始化数据库
func InitDb() error {
	var err error
	// 新建数据库
	nowDbFileName := getDbName()
	if nowDbFileName == "" {
		return errors.New(fmt.Sprintf(`InitDb().getDbName() is empty, not support this OS.
you need implement getDbName() in file: internal/dao/init.go `))
	}

	dbDir := filepath.Dir(nowDbFileName)
	if my_util.IsDir(dbDir) == false {
		_ = os.MkdirAll(dbDir, os.ModePerm)
	}
	db, err = gorm.Open(sqlite.Open(nowDbFileName), &gorm.Config{})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect database, %s", err.Error()))
	}
	// 迁移 schema
	err = db.AutoMigrate(&models.HotFix{}, &models.SubFormatRec{})
	if err != nil {
		return errors.New(fmt.Sprintf("db AutoMigrate error, %s", err.Error()))
	}
	return nil

}

func getDbName() string {
	nowDbFileName := ""
	sysType := runtime.GOOS
	if sysType == "linux" {
		nowDbFileName = dbFileNameLinux
	}
	if sysType == "windows" {
		nowDbFileName = dbFileNameWindows
	}
	if sysType == "darwin" {
		home, _ := os.UserHomeDir()
		nowDbFileName = home + "/.config/chinesesubfinder/" + dbFileNameDarwin
	}
	return nowDbFileName
}

var (
	db   *gorm.DB
	once sync.Once
)

const (
	dbFileNameLinux   = "/config/settings.db"
	dbFileNameWindows = "settings.db"
	dbFileNameDarwin  = "settings.db"
)
