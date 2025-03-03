package config

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Loader struct {
	//Host用于适配配置中心
	Host string
	//从哪些地方读取
	Paths []string
	//是否要读取环境变量中的同名配置字段并覆盖原值
	EnableEnv bool
	EnvPrefix string
	v         *viper.Viper
}

func NewLoader(opts ...ConfigOption) *Loader {
	l := &Loader{}
	for _, opt := range opts {
		opt(l)
	}

	v := viper.New()
	path, _ := os.Getwd()
	if len(l.Paths) == 0 {
		v.AddConfigPath(path)
	}
	if l.EnableEnv {
		v.AutomaticEnv()
	}

	if l.EnvPrefix != "" {
		v.SetEnvPrefix(strings.Replace(strings.ToUpper(l.EnvPrefix), "-", "_", -1))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	for _, path = range l.Paths {
		v.AddConfigPath(path)
	}
	l.v = v
	return l
}

func (l *Loader) Load(file string, data interface{}) error {
	l.v.SetConfigFile(file)
	if err := l.v.ReadInConfig(); err != nil {
		return err
	}
	return l.v.Unmarshal(data)
}

func (l *Loader) LoadJSON(filename string, data interface{}) error {
	l.v.SetConfigName(filename)
	l.v.SetConfigType("json")
	if err := l.v.ReadInConfig(); err != nil {
		return err
	}
	return l.v.Unmarshal(data)
}

func (l *Loader) LoadYaml(filename string, data interface{}) error {
	l.v.SetConfigName(filename)
	l.v.SetConfigType("yaml")
	if err := l.v.ReadInConfig(); err != nil {
		return err
	}
	return l.v.Unmarshal(data)
}

// 监听配置文件是否改变,可以手动传入f来调节逻辑,f也可为nil
// 注意,若有多个配置文件被读入,则最终只能监听最后读入的文件
// 保存时多次触发回调函数是bug,修理TODO
func (l *Loader) Watch(data interface{}, f func(e fsnotify.Event)) {
	if f == nil {
		f = func(e fsnotify.Event) {
			log.Println("[config] 配置文件发生变动:", e.Name)
			if err := l.v.Unmarshal(data); err != nil {
				log.Println("[config] 重新加载配置文件失败:", err)
			}
		}
	}
	l.v.WatchConfig()
	l.v.OnConfigChange(f)
}
