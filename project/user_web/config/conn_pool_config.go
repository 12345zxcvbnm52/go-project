package config

type ConnPoolConfig struct {
	//用于指定客户端尝试建立与服务器的连接时的最大等待时间,
	//默认单位为秒
	DialTimeout int32 `mapstructure:"dial_timeout"`
	//最大尝试退避策略来尝试重连,默认单位为秒
	BackoffMaxTry int32 `mapstructure:"backoff_max_try"`
	//多少时间内连接未被访问过则尝试ping一次服务端以测试可用性
	//默认单位为秒
	KeepAliveCheck int32 `mapstructure:"keep_alive_check"`
	//每一次ping的最大等待时间,超过则断开连接,默认单位为秒
	KeepAliveTimeout int32 `mapstructure:"keep_alive_timeout"`
	//指定了客户端在连接级别上的初始窗口大小
	InitStreamWindowSize int32
	//指定了客户端在流级别的初始窗口大小
	InitConnWindowSize int32
}
