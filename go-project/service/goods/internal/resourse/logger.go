package resourse

import (
	"fmt"
	"kenshop/pkg/log"
)

func InitLogger() {
	Logger = log.MustNewOtelLogger(
		log.WithErrOutPaths(fmt.Sprintf("%s/log/error.log", Pwd)),
		log.WithOutputPaths(fmt.Sprintf("%s/log/info.log", Pwd)),
	)
}
