package orderconfig

import "sync"

type MysqlConf struct {
	Ip       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

type RedisConf struct {
	//要考虑以后的redis集群配置
	Ip       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type ConsulConf struct {
	Ip   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

// type LogConf struct {
// 	Level       string   `mapstructure:"level"`
// 	ErrLevel    string   `mapstructure:"error-level"`
// 	OutPaths    []string `mapstructure:"out-paths"`
// 	ErrOutPaths []string `mapstructure:"error-out-paths"`
// 	Format      string   `mapstructure:"format"`
// 	Development bool     `mapstructure:"development"`
// }

type OtelConf struct {
	Ip          string `mapstructure:"ip"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"service-name"`
}

type RocketmqConf struct {
	Ip   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
	//库存归还用的topic
	RebackTopic string `mapstructure:"reback-topic"`
	//用于检测订单支付超时的topic
	TimeoutTopic           string `mapstructure:"timeout-topic"`
	ProducerGroupName      string `mapstructure:"producer-group-name"`
	TransProducerGroupName string `mapstructure:"transaction-producer-group-name"`
	ConsumerGroupName      string `mapstructure:"consumer-group-name"`
}

type ServerConf struct {
	//这几个参数是用来服务注册的
	Name string   `mapstructure:"name"`
	Id   string   `mapstructure:"id"`
	Tags []string `mapstructure:"tags"`
	//ip和port不通过读取得到,而是通过动态获取,
	Ip      string
	Port    int
	Mtx     sync.RWMutex
	Version string `mapstructure:"version"`

	Mysql MysqlConf `mapstructure:"mysql"`
	//consul的ip和端口就静态绑定docker的consul地址
	Consul   ConsulConf   `mapstructure:"consul"`
	Redis    RedisConf    `mapstructure:"redis"`
	Otel     OtelConf     `mapstructure:"otel"`
	Rocketmq RocketmqConf `mapstructure:"rocketmq"`
}
