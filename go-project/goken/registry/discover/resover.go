package discover

import (
	"context"
	"errors"
	"time"

	"kenshop/goken/registry"
	"kenshop/pkg/log"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// 服务发现的具体逻辑由该结构体调用
// resolver用于解析客户端传入的address(dial系函数的第一个参数)
// 功能类似于http域名解析,通过传入的域名(也可以是直接地址)获得对应的连接
// grpc默认的解析器一般只能解析tcp连接(直接地址)
type Resolver struct {
	listener registry.Listener
	cc       resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	insecure bool
}

// 服务发现的核心逻辑,同样是失败无限重试
func (r *Resolver) listen() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		//这里的逻辑类似一直监听对应的注册中心,阻塞或等到有新消息发出
		ins, err := r.listener.ListenAndGet(context.TODO())
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Errorf("服务发现失败, err= %v", err)
			time.Sleep(time.Second)
			continue
		}
		if len(ins) == 0 {
			continue
		}
		//一旦发现新的节点信息就更新到grpc内部中
		r.update(ins)
	}
}

// 用于更新服务发现的实例列表
func (r *Resolver) update(ins []*registry.ServiceInstance) {
	addrs := make([]resolver.Address, 0)
	// 使用一个map来存储已处理的新发现的endpoint地址,避免冗余
	endpoints := make(map[string]struct{})

	for _, in := range ins {
		// ParseEndpoint 这里暂时先只考虑grpc,后续开放更多类型的协议
		// 这里似乎函数有问题

		for _, e := range in.Endpoints {
			endpoint := e.Host
			// 如果解析结果为空的endpoint,跳过此服务实例
			if endpoint == "" {
				continue
			}

			// 如果这个endpoint已经存在,跳过此地址(避免冗余)
			if _, ok := endpoints[endpoint]; ok {
				continue
			}

			endpoints[endpoint] = struct{}{}
			// 构建一个resolver.Address对象,用于表示一个服务的地址信息
			addr := resolver.Address{
				ServerName: in.Name,
				//该字段类似于metadata,记录任何需要的元数据
				Attributes: parseAttributes(in.Metadata),
				Addr:       endpoint,
			}

			// 可以为resolver.Address添加一个自定义属性,存储原始的ServiceInstance信息
			//addr.Attributes = addr.Attributes.WithValue("rawServiceInstance", in)
			addrs = append(addrs, addr)
		}

	}

	// 如果没有有效的地址信息但是却触发了update函数,需要记录警告日志
	if len(addrs) == 0 {
		log.Warnf("resolver未发现有效的地址信息, instances= %v", ins)
		return
	}

	// grpc给出的接口,是grpc内部获得,更新连接的核心代码
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Errorf("resolver服务更新失败, err= %v", err)
	}

	//可以在最后打日志记录:更新/修改了某个服务以便快速得知信息变更
}

func (r *Resolver) Close() {
	r.cancel()
	err := r.listener.StopListen(context.TODO())
	if err != nil {
		log.Errorf("resolver停止监听服务失败, err= %s", err)
	}
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {}

func parseAttributes(md map[string]string) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range md {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}
