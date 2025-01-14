package rpcserver

import (
	"fmt"
	"goken/server/rpcserver/interceptors"
	"net"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerOption func(o *Server)

type Server struct {
	*grpc.Server
	address    string
	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
	listener   net.Listener

	timeout  time.Duration
	health   *health.Server
	endpoint *url.URL
}

func newDefaultOptions() *Server {
	s := &Server{address: "127.0.0.1:0"}
	s.timeout = time.Second * 1
	return s
}

func NewServer(opts ...ServerOption) (*Server, error) {
	s := newDefaultOptions()
	for _, v := range opts {
		v(s)
	}
	if err := s.listen(); err != nil {
		panic(err)
	}
	s.unaryInts = append(s.unaryInts, interceptors.UnaryTimeoutInterceptor(s.timeout))
	s.grpcOpts = append(s.grpcOpts, grpc.ChainUnaryInterceptor(s.unaryInts...))
	s.Server = grpc.NewServer(s.grpcOpts...)
	grpc_health_v1.RegisterHealthServer(s.Server, s.health)
	return s, nil
}

func WithAddress(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.timeout = timeout
	}
}

func WithListener(lis net.Listener) ServerOption {
	return func(o *Server) {
		o.listener = lis
	}
}

func WithUnaryInts(ui ...grpc.UnaryServerInterceptor) ServerOption {
	return func(o *Server) {
		o.unaryInts = ui
	}
}

func WithSteamInts(ui ...grpc.StreamServerInterceptor) ServerOption {
	return func(o *Server) {
		o.streamInts = ui
	}
}

func WithGrpcOptions(opts ...grpc.ServerOption) ServerOption {
	return func(o *Server) {
		o.grpcOpts = opts
	}
}

func (s *Server) listen() error {
	if s.listener == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.listener = lis

	}
	addr, err := extract(s.address, s.listener)
	if err != nil {
		return err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return nil
}

func isValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

// 获得合法的port
func getPort(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// 返回一个有效的host和port
func extract(address string, lis net.Listener) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil && lis == nil {
		return "", err
	}
	if lis != nil {
		if p, ok := getPort(lis); ok {
			port = strconv.Itoa(p)
		} else {
			return "", fmt.Errorf("这个地址无法获取有效的ip和port: %v", lis.Addr())
		}
	}
	if len(host) > 0 && (host != "0.0.0.0" && host != "[::]" && host != "::") {
		return net.JoinHostPort(host, port), nil
	}
	//?
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	lowest := int(^uint(0) >> 1)
	var result net.IP
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index < lowest || result == nil {
			lowest = iface.Index
		}
		if result != nil {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isValidIP(ip.String()) {
				result = ip
			}
		}
	}
	if result != nil {
		return net.JoinHostPort(result.String(), port), nil
	}
	return "", nil
}
