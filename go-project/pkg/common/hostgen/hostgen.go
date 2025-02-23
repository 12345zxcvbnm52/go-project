package hostgen

import (
	"net"
	"strconv"
)

func IsValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

func GetUsagePort() int {
	//把字符串解析为tcp端点
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	//tcp协议中,如果端口为0则在listen,dial这些函数中会默认给它分配一个空闲的端口
	lis, _ := net.ListenTCP("tcp", addr)
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port

}

// 解析传入的host和lis,选取(若无适合的ip或port则自动生成一个可用的)两者中适合的地址(即可使用的ip与port),
func ResolveHost(host string) (string, error) {
	ip, port, err := net.SplitHostPort(host)
	if err != nil {
		return "", err
	}

	//如果port为无效值则自动获取port
	if uport, err := strconv.Atoi(port); uport <= 0 || err != nil {
		uport = GetUsagePort()
		port = strconv.Itoa(uport)
	}

	//如果ip与port是可用的就直接返回
	if len(ip) > 0 && (ip != "0.0.0.0" && ip != "[::]" && ip != "::") {
		return net.JoinHostPort(ip, port), nil
	}

	//通过以下逻辑获得可使用的ip
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
			if IsValidIP(ip.String()) {
				result = ip
			}
		}
	}
	if result != nil {
		return net.JoinHostPort(result.String(), port), nil
	}
	return "", nil
}
