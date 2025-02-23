package discover

//一个自定义的gRPC服务发现resolver逻辑需要实现google.golang.org/grpc/resolver包中resolver.Builder和resolver.Resolver两个接口

//resolver.Builder负责解析服务名称并决定如何生成resolver.Resolver实例,
//使用时注册一个resolver.Builder,传入gRPC中以便根据给定的服务名称来构建自定义的resolver.Resolver

//resolver.Resolver是具体的服务发现逻辑,负责获取服务地址列表并维护服务变化的通知,
//你需要实现ResolveNow方法来查询和返回当前的服务地址信息
