package resourse

import (
	"fmt"
	"kenshop/pkg/redlock"
	inventorydata "kenshop/service/inventory/internal/data"
	model "kenshop/service/inventory/internal/model"
	"log"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
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
	InitRocketmq()

	InventoryData = inventorydata.MustNewGormInventoryData(db,
		redlock.MustNewRedLock([]string{fmt.Sprintf("%s:%d", Conf.Redlock.Ip, Conf.Redlock.Port)},
			redlock.WithPassword(Conf.Redlock.Password),
		),
	)
	Consumer.Subscribe(Conf.Rocketmq.RebackTopic, consumer.MessageSelector{}, InventoryData.Reback)
	if err := db.AutoMigrate(&model.Inventory{}, &model.OrderDecrRecord{}); err != nil {
		panic(err)
	}
}
