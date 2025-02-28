package userdata

import (
	"context"
	proto "kenshop/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ UserDataService = &GrpcUserData{}

// data层为数据库交互层
// 用户服务
type GrpcUserData struct {
	Cli proto.UserClient
}

func MustNewGrpcUserData(cli *grpc.ClientConn) *GrpcUserData {
	c := proto.NewUserClient(cli)
	return &GrpcUserData{
		Cli: c,
	}
}

// 获得用户列表,可通过FliterReq过滤
func (s *GrpcUserData) GetUserListDB(ctx context.Context, in *proto.UserFliterReq) (*proto.UserListRes, error) {
	return s.Cli.GetUserList(ctx, in)
}

// 通过用户id获取用户信息
func (s *GrpcUserData) GetUserByIdDB(ctx context.Context, in *proto.UserIdReq) (*proto.UserInfoRes, error) {
	return s.Cli.GetUserById(ctx, in)
}

// 通过用户电话号码获取用户信息
func (s *GrpcUserData) GetUserByMobileDB(ctx context.Context, in *proto.UserMobileReq) (*proto.UserInfoRes, error) {
	return s.Cli.GetUserByMobile(ctx, in)
}

// 创建一个用户
func (s *GrpcUserData) CreateUserDB(ctx context.Context, in *proto.CreateUserReq) (*proto.CreateUserRes, error) {
	return s.Cli.CreateUser(ctx, in)
}

// 更新用户,传入的用户信息字段中无论是否为空都会完全覆盖原来的值
func (s *GrpcUserData) AbsUpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	return s.Cli.AbsUpdateUser(ctx, in)
}

// 局部更新设置了值的参数
func (s *GrpcUserData) UpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	return s.Cli.UpdateUser(ctx, in)
}

// 注销一个用户
func (s *GrpcUserData) DeleteUserDB(ctx context.Context, in *proto.DelUserReq) (*emptypb.Empty, error) {
	return s.Cli.DeleteUser(ctx, in)
}

// 权限验证
func (s *GrpcUserData) CheckUserRoleDB(ctx context.Context, in *proto.UserPasswordReq) (*proto.UserCheckRes, error) {
	return s.Cli.CheckUserRole(ctx, in)
}
