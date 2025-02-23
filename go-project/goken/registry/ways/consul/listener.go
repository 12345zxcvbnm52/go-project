package consul

import (
	"context"
	"kenshop/goken/registry"
)

// listener扮演消费者的形象,讲consul中服务的变化信息当做可消费信息,并将信息push到各个客户端中
type listener struct {
	//类似于条件变量的channel
	event   chan struct{}
	serivce *serviceInfo

	ctx    context.Context
	cancel context.CancelFunc
}

func (l *listener) StopListen(_ context.Context) error {
	l.cancel()
	l.serivce.lock.Lock()
	defer l.serivce.lock.Unlock()
	delete(l.serivce.listener, l)
	return nil
}

// 具体的消费行为,即把发现得到的ServiceInstance更新到服务器中
func (w *listener) ListenAndGet(_ context.Context) ([]*registry.ServiceInstance, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
		//消费
	case <-w.event:
	}
	w.serivce.lock.Lock()
	defer w.serivce.lock.Unlock()
	services := make([]*registry.ServiceInstance, 0)
	for {
		if w.serivce.serviceQueue.Length() > 0 {
			services = append(services, w.serivce.serviceQueue.Pop())
		} else {
			return services, nil
		}
	}
}
