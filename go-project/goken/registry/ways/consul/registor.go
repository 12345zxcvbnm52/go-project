package consul

import (
	"context"
	"fmt"
	"kenshop/goken/registry"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
)

// Client是对consul client的一层封装
type registor struct {
	cli     *api.Client
	timeout string
	//健康检查的时间间隔
	healthcheckInterval            string
	deregisterCriticalServiceAfter string
	// 用户指定的自定义需要检查的地址
	serviceChecks api.AgentServiceChecks
	//心跳检查标志
	heartBeat         bool
	enableHealthCheck bool
	ttlTimeout        string
}

type RegistorOption func(*registor)

// 这里为了灵活性选择让调用者自己传入api.Client
func MustNewConsulRegistor(apiClient *api.Client, opts ...RegistorOption) *registor {
	r := &registor{
		cli:                            apiClient,
		timeout:                        "20s",
		healthcheckInterval:            "20s",
		deregisterCriticalServiceAfter: "1m",
		enableHealthCheck:              true,
		ttlTimeout:                     "5s",
		heartBeat:                      false,
	}

	for _, o := range opts {
		o(r)
	}
	return r
}

// 服务注册接口
func (r *registor) Register(ctx context.Context, ins *registry.ServiceInstance) error {
	//这里map的key是web应用层协议,后续consul注册也是通过协议名注册
	addresses := make(map[string]api.ServiceAddress, len(ins.Endpoints))
	//记录所有需要健康检查的URL
	checkAddresses := make([]*url.URL, 0, len(ins.Endpoints))
	//这里通过一系列操作用于记录对应的信息
	for _, endpoint := range ins.Endpoints {
		port, _ := strconv.ParseInt(endpoint.Port(), 10, 16)
		checkAddresses = append(checkAddresses, endpoint)
		addresses[endpoint.Scheme] = api.ServiceAddress{Address: endpoint.Hostname(), Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:   ins.ID,
		Name: ins.Name,
		Meta: ins.Metadata,
		Tags: []string{fmt.Sprintf("version=%s", ins.Version)},
		//TaggedAddresses 用于一次注册多个地址
		TaggedAddresses: addresses,
	}

	//拿第一个地址做默认地址
	if len(checkAddresses) > 0 {
		asr.Address = checkAddresses[0].Hostname()
		port, _ := strconv.ParseInt(checkAddresses[0].Port(), 10, 16)
		asr.Port = int(port)
	}

	if r.enableHealthCheck {
		for _, address := range checkAddresses {
			switch address.Scheme {
			case "grpc":
				asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
					GRPC:                           address.Host,
					Interval:                       r.healthcheckInterval,
					DeregisterCriticalServiceAfter: r.deregisterCriticalServiceAfter,
					Timeout:                        r.timeout,
				})

			case "http", "https":
				asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
					HTTP:                           address.Host + "/health",
					Interval:                       r.healthcheckInterval,
					DeregisterCriticalServiceAfter: r.deregisterCriticalServiceAfter,
					Timeout:                        r.timeout,
				})
			}
		}
		//把用户需要检查的URL也记录下来
		if r.serviceChecks != nil {
			asr.Checks = append(asr.Checks, r.serviceChecks...)
		}
	}

	//相比于上面的检查模式,TTL模式要求服务主动向consul发送请求来确定服务是健康的
	if r.heartBeat {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			CheckID:                        "service:" + ins.ID,
			TTL:                            r.ttlTimeout,
			DeregisterCriticalServiceAfter: r.deregisterCriticalServiceAfter,
		})
	}

	err := r.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	return nil
}

// func (r *registor) HeatBeat(ctx context.Context,ins *registry.ServiceInstance) {
// 	// 在首次运行时稍微等待1秒钟,确保其他启动过程已完成
// 	time.Sleep(time.Second)

// 	// 初始化服务健康检查状态为pass,表示服务健康
// 	err := r.cli.Agent().UpdateTTL("service:"+ins.ID, "pass", "pass")
// 	if err != nil {
// 		log.Errorf("Consul心跳检查机制初始化服务状态失败 err= %v", err)
// 	}

// 	// 创建一个定时器,根据healthcheckInterval来控制心跳更新的间隔
// 	ttlTimeout, _ := time.ParseDuration(r.ttlTimeout)
// 	ttlTimeout = ttlTimeout / 2
// 	ticker := time.NewTicker(time.Second * time.Duration(ttlTimeout))
// 	defer ticker.Stop() // 确保在函数退出时停止ticker,避免内存泄漏

// 	// 启动一个死循环,不断检查心跳,并处理服务的健康检查更新
// 	for {
// 		select {
// 		// 如果上下文已取消或超时,注销服务并退出
// 		case <-ctx.Done():
// 			_ = r.cli.Agent().ServiceDeregister(ins.ID)
// 			return
// 		// 如果ticker时间到,更新TTL心跳
// 		case <-ticker.C:
// 			// 一层保险,如果上下文已被取消或者超时,注销服务并退出
// 			if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
// 				_ = r.cli.Agent().ServiceDeregister(ins.ID)
// 				return
// 			}

// 			// 更新服务的 TTL(Time-To-Live)状态为pass
// 			err = r.cli.Agent().UpdateTTLOpts("service:"+ins.ID, "pass", "pass", new(api.QueryOptions).WithContext(ctx))
// 			//同样是一层保险
// 			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
// 				// 如果上下文被取消或超时,注销服务并退出
// 				_ = r.cli.Agent().ServiceDeregister(ins.ID)
// 				return
// 			}

// 			// 如果更新TTL时发生了错误,尝试重新注册该服务
// 			if err != nil {
// 				log.Errorf("Consul心跳检查机制更新服务状态失败 err=%v", err)

// 				// 如果更新失败,稍微等待一下,然后尝试重新注册服务
// 				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

// 				// 重新注册服务
// 				if err := r.cli.Agent().ServiceRegister(asr); err != nil {
// 					log.Errorf("Consul心跳检查机制重新注册失败 err=%v", err)
// 				} else {
// 					log.Warn("Consul心跳检查机制重新注册成功")
// 				}
// 			}
// 		}
// 	}
// }

// 移除服务中心的服务
func (r *registor) Deregister(ctx context.Context, serviceID string) error {
	return r.cli.Agent().ServiceDeregister(serviceID)
}

// 设置Timeout,
func WithTimeout(timeout string) RegistorOption {
	if _, valid := time.ParseDuration(timeout); valid != nil {
		return func(r *registor) {}
	}
	return func(r *registor) {
		r.timeout = timeout
	}
}

func WithTTLtimeout(ttlTimeout string) RegistorOption {
	if _, valid := time.ParseDuration(ttlTimeout); valid != nil {
		return func(r *registor) {}
	}
	return func(r *registor) {
		r.ttlTimeout = ttlTimeout
	}
}

// 设置healthcheckInterval
func WithHealthcheckInterval(interval string) RegistorOption {
	if _, valid := time.ParseDuration(interval); valid != nil {
		return func(r *registor) {}
	}
	return func(r *registor) {
		r.healthcheckInterval = interval
	}
}

// 设置deregisterCriticalServiceAfter
func WithDeregisterCriticalServiceAfter(after string) RegistorOption {
	if _, valid := time.ParseDuration(after); valid != nil {
		return func(r *registor) {}
	}
	return func(r *registor) {
		r.deregisterCriticalServiceAfter = after
	}
}

// 设置自定义的服务检查
func WithServiceChecks(checks api.AgentServiceChecks) RegistorOption {
	return func(r *registor) {
		r.serviceChecks = checks
	}
}

// 启用心跳检查
func WithHeartBeat(enabled bool) RegistorOption {
	return func(r *registor) {
		r.heartBeat = enabled
	}
}

// 启用健康检查
func WithEnableHealthCheck(enabled bool) RegistorOption {
	return func(r *registor) {
		r.enableHealthCheck = enabled
	}
}
