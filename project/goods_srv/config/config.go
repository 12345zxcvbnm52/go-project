package config

type MysqlConfig struct {
	NetType string `mapstructure:"net_type"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

type ConsulConfig struct {
	ConsulIp   string `mapstructure:"consul_ip"`
	ConsulPort int    `mapstructure:"consul_port"`
}

type ServerConfig struct {
	//这个名字是用来服务注册的
	Name string `mapstructure:"name"`
	//ip和port不通过读取得到,而是通过动态获取,
	Ip          string
	Port        int
	Id          string      `mapstructure:"id"`
	Tags        []string    `mapstructure:"tags"`
	MysqlConfig MysqlConfig `mapstructure:"mysql"`
	//consul的ip和端口就静态绑定docker的consul地址
	ConsulConfig ConsulConfig `mapstructure:"consul"`
}
