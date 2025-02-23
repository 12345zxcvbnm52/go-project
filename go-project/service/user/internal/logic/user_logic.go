package userlogic

import (
	"context"
	proto "kenshop/proto/user"
	userdata "kenshop/service/user/internal/data"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Service层为主要逻辑层
// 用户服务
type UserService struct {
	UserData userdata.UserDataService
}

// 获得用户列表,可通过FliterReq过滤
func (s *UserService) GetUserListLogic(ctx context.Context, in *proto.UserFliterReq) (*proto.UserListRes, error) {
	return s.UserData.GetUserListDB(ctx, in)
}

// 通过用户id获取用户信息
func (s *UserService) GetUserByIdLogic(ctx context.Context, in *proto.UserIdReq) (*proto.UserInfoRes, error) {
	return s.UserData.GetUserByIdDB(ctx, in)
}

// 通过用户电话号码获取用户信息
func (s *UserService) GetUserByMobileLogic(ctx context.Context, in *proto.UserMobileReq) (*proto.UserInfoRes, error) {
	return s.UserData.GetUserByMobileDB(ctx, in)
}

// 创建一个用户
func (s *UserService) CreateUserLogic(ctx context.Context, in *proto.CreateUserReq) (*proto.CreateUserRes, error) {
	return s.UserData.CreateUserDB(ctx, in)
}

// 更新用户,传入的用户信息字段中无论是否为空都会完全覆盖原来的值
func (s *UserService) AbsUpdateUserLogic(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	return s.UserData.AbsUpdateUserDB(ctx, in)
}

// 局部更新设置了值的参数
func (s *UserService) UpdateUserLogic(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	return s.UserData.UpdateUserDB(ctx, in)
}

// 注销一个用户
func (s *UserService) DeleteUserLogic(ctx context.Context, in *proto.DelUserReq) (*emptypb.Empty, error) {
	return s.UserData.DeleteUserDB(ctx, in)
}

// 权限验证
func (s *UserService) CheckUserRoleLogic(ctx context.Context, in *proto.UserPasswordReq) (*proto.UserCheckRes, error) {
	return s.UserData.CheckUserRoleDB(ctx, in)
}
