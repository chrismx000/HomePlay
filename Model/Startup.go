package model

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&VideoInfo{})   //创建posts表时指存在引擎
	db.AutoMigrate(&VideoDetail{}) //创建posts表时指存在引擎
}

func Startup() {

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:123@tcp(127.0.0.1:3306)/video?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize: 256,                                                                        // string 类型字段的默认长度
		//DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		//DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		//DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		PrepareStmt:                              true, // 全局模式，所有 DB 操作都会创建并缓存预编译语句
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger, //注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
		NamingStrategy:schema.NamingStrategy{
			SingularTable: true,
		},
		
	})
	if err != nil {
		panic(err)
	}
	Migrate(DB)
}
