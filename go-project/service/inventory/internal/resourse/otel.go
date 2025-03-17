package resourse

import (
	"fmt"
	ktrace "kenshop/pkg/trace"
)

func InitOtel() {
	if err := ktrace.RegistorTP(Ctx, fmt.Sprintf("%s:%d", Conf.Otel.Ip, Conf.Otel.Port), ktrace.WithName(Conf.Otel.ServiceName)); err != nil {
		panic(err)
	}
}
