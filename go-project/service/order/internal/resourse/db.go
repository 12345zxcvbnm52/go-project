package resourse

import (
	"fmt"
	"kenshop/pkg/rockmq"
	gproto "kenshop/proto/goods"
	iproto "kenshop/proto/inventory"
	orderdata "kenshop/service/order/internal/data"
	model "kenshop/service/order/internal/model"
	"log"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
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

	Producer, err = rocketmq.NewProducer(
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", Conf.Rocketmq.Ip, Conf.Rocketmq.Port)}),
		producer.WithGroupName(Conf.Rocketmq.ProducerGroupName),
	)
	if err != nil {
		panic(err)
	}
	if err := Producer.Start(); err != nil {
		panic(err)
	}

	icli, err := InventoryClient.Dial()
	if err != nil {
		panic(err)
	}
	gcli, err := GoodsClient.Dial()
	if err != nil {
		panic(err)
	}
	iclient := iproto.NewInventoryClient(icli)
	gclient := gproto.NewGoodsClient(gcli)
	listener := orderdata.MustNewOrderListener(db, Logger, Conf.Rocketmq.TimeoutTopic, Conf.Rocketmq.RebackTopic, iclient, gclient, Producer)
	listener.Ctx = Ctx
	namesrv := fmt.Sprintf("%s:%d", Conf.Rocketmq.Ip, Conf.Rocketmq.Port)
	transProducer := rockmq.MustNewTransProducer(listener, []string{namesrv}, Conf.Rocketmq.TransProducerGroupName)

	OrderData = orderdata.MustNewGormOrderData(db, Logger, Conf.Rocketmq.TimeoutTopic, iclient, gclient, transProducer)

	TimeoutConsumer = rockmq.MustNewPushConsumer([]string{namesrv}, Conf.Rocketmq.ConsumerGroupName)
	TimeoutConsumer.Subscribe(Conf.Rocketmq.TimeoutTopic, consumer.MessageSelector{}, listener.OrderTimeout)
	if err := db.AutoMigrate(&model.Order{}, &model.Cart{}, &model.OrderGoods{}); err != nil {
		panic(err)
	}
}
