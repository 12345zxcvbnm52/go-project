package resourse

import (
	"fmt"
	userdata "kenshop/api/user/internal/data"
	"kenshop/api/user/internal/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Conf.Mysql.UserName,
		Conf.Mysql.Password,
		Conf.Mysql.Ip,
		Conf.Mysql.Port,
		Conf.Mysql.DBName,
	)
	filePath, err := os.OpenFile(fmt.Sprintf("%s/log/db.log", Pwd), os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	log := logger.New(
		log.New(filePath, "", log.LstdFlags|log.Llongfile),
		logger.Config{
			SlowThreshold: 1 * time.Second,
			Colorful:      false,
			LogLevel:      logger.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		//让所有gorm适配的数据库的Err同步为gorm.Err类型
		TranslateError: true,
		Logger:         log,
	})
	if err != nil {
		panic(err)
	}
	UserData = userdata.MustNewGormUserData(db)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
}
