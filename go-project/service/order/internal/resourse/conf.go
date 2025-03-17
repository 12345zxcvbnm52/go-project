package resourse

import (
	"fmt"
	"kenshop/pkg/config"
)

func InitConf() {
	ConfLoader = config.NewLoader(
		config.WithEnableEnv(true),
		config.WithPaths([]string{fmt.Sprintf("%s/etc", Pwd)}),
	)
	if err := ConfLoader.LoadYaml("order", &Conf); err != nil {
		panic(err)
	}
}
