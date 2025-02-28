package sinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const healthCheckMethod = "/grpc.health.v1.Health/Check"

// 健康检查拦截器
func HealthCheckInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 判断请求是否为健康检查方法
		if info.FullMethod == healthCheckMethod {
			return &grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_SERVING,
			}, nil
		}
		return handler(ctx, req)
	}
}
