package initialize

import (
	"fmt"
	gb "goods_srv/global"
	"os"

	"goods_srv/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	prefix, _ := os.Getwd()
	viper.SetEnvPrefix("GOODS_SRV")
	viper.AutomaticEnv()
	flag := viper.GetBool("MODEL")
	configFileName := ""
	if !flag {
		configFileName = "/goods_srv_release.yaml"
	} else {
		configFileName = "/goods_srv_debug.yaml"
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

	err = v.Unmarshal(&gb.ServerConfig)
	if err != nil {
		zap.S().Errorw("viper读出数据失败", "msg", err.Error())
		os.Exit(1)
	}

	//把动态的地址记录到配置文件中
	addr := util.NewTcpAddr()
	gb.ServerConfig.Ip = addr.IP.String()
	gb.ServerConfig.Port = addr.Port

	//这里监视配置文件的改动
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config changed", e.Name)
	})

}
