package config

type ConsulConfig struct {
	// 这个名字是用来服务发现的
	Name string `mapstructure:"name"`
	// 读取静态的consul的ip和port
	ConsulIp   string `mapstructure:"consul_ip"`
	ConsulPort int    `mapstructure:"consul_port"`
}

type HealthCheckConfig struct {
	//也是读取静态的健康检查ip与port,用于web服务注册
	HealthCheckIp   string `mapstructure:"health_check_ip"`
	HealthCheckPort int    `mapstructure:"health_check_port"`
}

type OrderServerConfig struct {
	OrderServerId   string
	OrderServerPort int
	OrderServerIp   string
	OrderServerTags []string
}

type ServerConfig struct {
	//consul的ip和端口就静态绑定docker的consul地址
	ConsulConfig      ConsulConfig      `mapstructure:"consul"`
	HealthCheckConfig HealthCheckConfig `mapstructure:"health_check"`
	//给web服务注册到consul
	Name string `mapstructure:"name"`
	//web服务使用静态绑定ip和port
	Ip   string   `mapstructure:"ip"`
	Port int      `mapstructure:"port"`
	Id   string   `mapstructure:"id"`
	Tags []string `mapstructure:"tags"`

	Version string `mapstructure:"version"`
	//用来存微服务层的数据
	OrderServerConfig OrderServerConfig
	JwtSign           string `mapstructure:"jwt_sign"`
}
