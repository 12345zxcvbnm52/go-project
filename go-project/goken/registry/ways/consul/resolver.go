package consul

import (
	"context"
	"kenshop/goken/registry"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
)

// 自定义的用于将consul的ServiceEntry转化为ServiceInstance
type ServiceResolveFunc func(context.Context, []*api.ServiceEntry) []*registry.ServiceInstance

// 把得到的entries解析为对应的ServiceInstance
func ServiceResove(_ context.Context, entries []*api.ServiceEntry) []*registry.ServiceInstance {
	services := make([]*registry.ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		var version string
		//得到version信息
		for _, tag := range entry.Service.Tags {
			ss := strings.SplitN(tag, "=", 2)
			if len(ss) == 2 && ss[0] == "version" {
				version = ss[1]
			}
		}
		endpoints := make([]*url.URL, 0)
		//遍历一个服务注册的所有host,如果是局域网就不要暴露给外界
		for scheme, addr := range entry.Service.TaggedAddresses {
			if scheme == "lan_ipv4" || scheme == "wan_ipv4" || scheme == "lan_ipv6" || scheme == "wan_ipv6" {
				continue
			}
			url := &url.URL{
				Scheme: scheme,
				Host:   net.JoinHostPort(addr.Address, strconv.Itoa(addr.Port)),
			}
			endpoints = append(endpoints, url)
		}
		//如果TaggedAddresses内无有效已注册服务地址就使用Address里的地址
		if len(endpoints) == 0 && entry.Service.Address != "" && entry.Service.Port != 0 {
			url := &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(entry.Service.Address, strconv.Itoa(entry.Service.Port)),
			}
			endpoints = append(endpoints, url)
		}
		services = append(services, &registry.ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Version:   version,
			Endpoints: endpoints,
		})
	}
	return services
}
