package userdata

import (
	"context"
	proto "kenshop/proto/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UserDataService 是提供用户相关数据操作的接口
type UserDataService interface {
	GetUserListDB(ctx context.Context, in *proto.UserFliterReq) (*proto.UserListRes, error)
	GetUserByIdDB(ctx context.Context, in *proto.UserIdReq) (*proto.UserInfoRes, error)
	GetUserByMobileDB(ctx context.Context, in *proto.UserMobileReq) (*proto.UserInfoRes, error)
	CreateUserDB(ctx context.Context, in *proto.CreateUserReq) (*proto.CreateUserRes, error)
	AbsUpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error)
	UpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error)
	DeleteUserDB(ctx context.Context, in *proto.DelUserReq) (*emptypb.Empty, error)
	CheckUserRoleDB(ctx context.Context, in *proto.UserPasswordReq) (*proto.UserCheckRes, error)
}
