package discover

import (
	"context"
	"fmt"
	"kenshop/goken/registry"
	"kenshop/pkg/errors"
	"net/url"
	"strings"

	"google.golang.org/grpc/resolver"
)

type KDtmOption func(*kDtmDriver)

func MustNewKDtmDriver(ctx context.Context, b resolver.Builder, opts ...KDtmOption) *kDtmDriver {
	k := &kDtmDriver{
		Ctx:     ctx,
		Builder: b,
	}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

type kDtmDriver struct {
	Register registry.Registor
	Builder  resolver.Builder
	Ctx      context.Context
}

func (s *kDtmDriver) GetName() string {
	return "dtm-driver-goken"
}

func (b *kDtmDriver) RegisterAddrResolver() {
	if b.Builder != nil {
		resolver.Register(b.Builder)
	}
}

func (b *kDtmDriver) RegisterService(target string, endpoint string) error {
	//先不注册,看看这个target和endpoint传入的是什么
	// if b.Register != nil {
	// 	b.Register.Register(b.Ctx,&registry.ServiceInstance{
	// 		ID: fmt.Sprintf("dtm-driver-goken"),
	// 	})
	// }
	fmt.Println(target, endpoint)
	return nil
}

func (b *kDtmDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", err
	}
	service := strings.Split(u.Path, "/")
	if len(service) != 4 {
		return "", "", errors.New("uri格式错误,无法解析出server-name与full method")
	}
	server = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, service[1])
	method = fmt.Sprintf("/%s/%s", service[2], service[3])
	return
}

func WithRegistor(r registry.Registor) KDtmOption {
	return func(kdd *kDtmDriver) {
		kdd.Register = r
	}
}
