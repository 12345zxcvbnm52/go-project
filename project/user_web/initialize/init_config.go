package initialize

import (
	"os"
	gb "user_web/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 读取主配置文件
func InitConfig() {
	prefix, _ := os.Getwd()
	viper.SetEnvPrefix("USER_SRV")
	viper.AutomaticEnv()
	flag := viper.GetBool("MODEL")
	configFileName := ""
	if !flag {
		configFileName = "/user_web_release.yaml"
	} else {
		configFileName = "/user_web_debug.yaml"
	}
	v := viper.New()
	v.SetConfigName(configFileName)
	v.SetConfigType("yaml")
	v.AddConfigPath(prefix)
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Errorw("viper读入配置文件失败", "msg", err.Error())
		os.Exit(1)
	}
	//viper读取配置文件时哪怕是指针也要取地址
	err = v.Unmarshal(&gb.ServerConfig)
	if err != nil {
		zap.S().Errorw("viper读出数据失败", "msg", err.Error())
		os.Exit(1)
	}

	//这里监视配置文件的改动
	//可以考虑全局切片化v,并且包装一层结构体来检查信息
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("config changed", e.Name)
	// })
	initConnPoolConfig()
}

func initConnPoolConfig() {
	v := viper.New()
	prefix, _ := os.Getwd()
	v.AddConfigPath(prefix)
	fileName := "conn_pool.yaml"
	//ConfigName不需要后缀

	v.SetConfigFile(fileName)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Errorw("读取连接池配置文件失败", "msg", err.Error())
		panic(err)
	}
	v.Unmarshal(&gb.ConnPoolConfig)
}
