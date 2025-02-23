package consul

import (
	"context"
	"fmt"
	"kenshop/goken/registry"
	"kenshop/pkg/log"
	"kenshop/pkg/queue"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type discover struct {
	dc  string
	cli *api.Client
	//缓存记录被加入监听的服务信息
	mapServices map[string]*serviceInfo
	lock        sync.RWMutex
	//服务发现中发现,更新服务信息的周期,默认为短轮询
	ttl time.Duration

	//用于将从consul得到的服务描述结构体转为包内的ServiceInstance
	serviceResolver ServiceResolveFunc
}

type DiscoverOption func(d *discover)

const (
	SingleDC string = "SINGLE"
	MultiDC  string = "MULTI"
)

func MustNewConsulDiscover(cli *api.Client, options ...DiscoverOption) *discover {
	d := &discover{
		dc:              SingleDC,
		cli:             cli,
		ttl:             time.Second * 10,
		mapServices:     make(map[string]*serviceInfo),
		serviceResolver: ServiceResove,
	}
	// 应用所有传入的选项
	for _, opt := range options {
		opt(d)
	}
	return d
}

// 为指定服务名称创建一个服务监听器,服务监听器会定时push服务的最新信息
func (d *discover) NewListener(ctx context.Context, serviceName string) (registry.Listener, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 在缓存中中查找服务实例集合
	service, ok := d.mapServices[serviceName]
	// 如果缓存不存在,就从Consul或其他服务注册中心加载服务信息
	if !ok {
		service = &serviceInfo{
			listener: make(map[*listener]struct{}),
			//services:    queue.Queue[],
			serviceName:  serviceName,
			serviceQueue: queue.NewQueue[*registry.ServiceInstance](),
		}
		d.mapServices[serviceName] = service
	}

	l := &listener{
		event: make(chan struct{}, 1),
	}
	// 继承父context并且自己再产生一个cancel用于调控
	l.ctx, l.cancel = context.WithCancel(ctx)

	// 关联监听器和service结构体
	l.serivce = service
	service.lock.Lock()
	service.listener[l] = struct{}{}
	service.lock.Unlock()
	if service.serviceQueue.Length() > 0 {
		//l.event <- struct{}{}
		// 如果服务实例集合中有新数据则唤醒listener
		service.storeAndAwake(nil)
	}

	// 这里是真正调用consul-api的函数
	if err := d.discover(ctx, service); err != nil {
		return nil, err
	}

	return l, nil
}

// 启动对consul服务发现的轮询,需保证每一个service有且只调用一次这个函数
func (d *discover) discover(ctx context.Context, ss *serviceInfo) error {
	//初始index设为0,以此获取最新服务数据
	services, idx, err := d.DiscoverService(ctx, ss.serviceName, 0, true)
	if err != nil {
		return err
	}

	//叫醒所有消费者
	if len(services) > 0 {
		ss.storeAndAwake(services)
	}
	//这里开始定时检查更新consul服务信息,上面部分则是调用刚开始时立马查询更新一次
	go func() {
		ticker := time.NewTicker(d.ttl)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//用idx进行监听,监听逻辑在此
				tmpService, tmpIdx, err := d.DiscoverService(ctx, ss.serviceName, idx, true)
				if err != nil {
					log.Errorf("[consul discover] 轮询中发现服务失败 err= %v", err)
					//注意这里出了错并没检验是什么错
					time.Sleep(time.Second)
					continue
				}
				if len(tmpService) != 0 && tmpIdx != idx {
					services = tmpService
					ss.storeAndAwake(services)
				}
				idx = tmpIdx
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// 从sonsul中发现发生改变的服务信息
func (d *discover) DiscoverService(ctx context.Context, serviceName string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	if d.dc == MultiDC {
		return d.multiDCService(ctx, serviceName, index, passingOnly)
	}

	opts := &api.QueryOptions{
		//这里的Index是:当consul发生每次Consul对服务数据进行修改时(例如对一个ID已存在的服务进行再次注册),
		//都会生成一个递增的Index,WaitIndex的作用是告知Consul客户端从哪个索引开始阻塞查询,客户端会等待直到当前索引变更,或者超时发生
		WaitIndex:  index,
		WaitTime:   time.Second * 55,
		Datacenter: "", //d.dc,
		//强制保证多节点数据一致性
		//RequireConsistent: true,
	}
	opts = opts.WithContext(ctx)

	// if d.dc == SingleDC {
	// 	opts.Datacenter = "dc1"
	// }

	entries, meta, err := d.singleDCEntries(serviceName, "", passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}
	return d.serviceResolver(ctx, entries), meta.LastIndex, nil
}

func (d *discover) multiDCService(ctx context.Context, service string, index uint64, passingOnly bool) ([]*registry.ServiceInstance, uint64, error) {
	if ctx == nil {
		fmt.Println(service, index, passingOnly)
	}
	// 	opts := &api.QueryOptions{
	// 		WaitIndex: index,
	// 		WaitTime:  time.Second * 55,
	// 	}
	// 	opts = opts.WithContext(ctx)

	// 	var instances []*registry.ServiceInstance

	// 	dcs, err := d.cli.Catalog().Datacenters()
	// 	if err != nil {
	// 		return nil, 0, err
	// 	}

	// 	for _, dc := range dcs {
	// 		opts.Datacenter = dc
	// 		e, m, err := d.singleDCEntries(service, "", passingOnly, opts)
	// 		if err != nil {
	// 			return nil, 0, err
	// 		}

	// 		ins := d.resolver(ctx, e)
	// 		for _, in := range ins {
	// 			if in.Metadata == nil {
	// 				in.Metadata = make(map[string]string, 1)
	// 			}
	// 			in.Metadata["dc"] = dc
	// 		}

	// 		instances = append(instances, ins...)
	// 		opts.WaitIndex = m.LastIndex
	// 	}

	// return instances, opts.WaitIndex, nil
	return nil, 0, nil
}

func (d *discover) singleDCEntries(service, tag string, passingOnly bool, opts *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	return d.cli.Health().Service(service, tag, passingOnly, opts)
}

func (d *discover) RegisrtyName() string {
	return "consul"
}

// 对被注册到consul的服务的抽象
type serviceInfo struct {
	//服务名称
	serviceName string
	//监听器集合,若serviceInfo发生变化(对应的consul服务发生变化),可通过这些listener通知监听程序
	listener map[*listener]struct{}
	//还没被消费掉的Service信息
	//services []*registry.ServiceInstance
	serviceQueue *queue.Queue[*registry.ServiceInstance]

	lock sync.RWMutex
}

// 把服务数据存入到serviceInfo中并唤醒所有消费者
func (s *serviceInfo) storeAndAwake(ss []*registry.ServiceInstance) {
	for _, service := range ss {
		s.serviceQueue.Push(service)
	}

	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.listener {
		select {
		//广播唤醒所有消费者(类似于生产了每个消费者所需的信息)
		case k.event <- struct{}{}:
		default:
		}
	}
}

// 示例的选项函数
func WithDataCenter(dc string) DiscoverOption {
	return func(d *discover) {
		d.dc = dc
	}
}

func WithTTL(ttl time.Duration) DiscoverOption {
	return func(d *discover) {
		d.ttl = ttl
	}
}

func WithServiceResolver(resolver ServiceResolveFunc) DiscoverOption {
	return func(d *discover) {
		d.serviceResolver = resolver
	}
}
