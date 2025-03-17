package goodsconfig

import "sync"

// type MysqlConf struct {
// 	Ip       string `mapstructure:"ip"`
// 	Port     int    `mapstructure:"port"`
// 	UserName string `mapstructure:"username"`
// 	Password string `mapstructure:"password"`
// 	DBName   string `mapstructure:"db_name"`
// }

type RedisConf struct {
	//要考虑以后的redis集群配置
	Ip       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	UserName string `mapstructure:"username"`
}

type ConsulConf struct {
	Ip         string `mapstructure:"ip"`
	Port       int    `mapstructure:"port"`
	HealthIp   string `mapstructure:"health-ip"`
	HealthPort int    `mapstructure:"health-port"`
}

// type LogConf struct {
// 	Level       string   `mapstructure:"level"`
// 	ErrLevel    string   `mapstructure:"error-level"`
// 	OutPaths    []string `mapstructure:"out-paths"`
// 	ErrOutPaths []string `mapstructure:"error-out-paths"`
// 	Format      string   `mapstructure:"format"`
// 	Development bool     `mapstructure:"development"`
// }

type JwtConf struct {
	Key string `mapstructure:"key"`
}

// 用于服务发现或直接连接
type GoodsServerConf struct {
	Name string   `mapstructure:"name"`
	Id   string   `mapstructure:"id"`
	Ip   string   `mapstructure:"ip"`
	Port int      `mapstructure:"port"`
	Tags []string `mapstructure:"tags"`
}

type OtelConf struct {
	Ip          string `mapstructure:"ip"`
	Port        int    `mapstructure:"port"`
	TracerName  string `mapstructure:"tracer-name"`
	ServiceName string `mapstructure:"service-name"`
}

type ServerConf struct {
	//这几个参数是用来服务注册的
	Name    string   `mapstructure:"name"`
	Id      string   `mapstructure:"id"`
	Tags    []string `mapstructure:"tags"`
	Ip      string   `mapstructure:"ip"`
	Port    int      `mapstructure:"port"`
	Version string   `mapstructure:"version"`

	//	Mysql MysqlConf `mapstructure:"mysql"`
	//consul的ip和端口就静态绑定docker的consul地址
	Consul   ConsulConf      `mapstructure:"consul"`
	Redis    RedisConf       `mapstructure:"redis"`
	Jwt      JwtConf         `mapstructure:"jwt"`
	Otel     OtelConf        `mapstructure:"otel"`
	GoodsSrv GoodsServerConf `mapstructure:"goods-server"`
	Mtx      sync.RWMutex
}
