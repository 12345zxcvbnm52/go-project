package usercontroller

import (
	"context"
	proto "kenshop/proto/user"
	userlogic "kenshop/service/user/internal/logic"

	gproto "google.golang.org/protobuf/proto"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Contoller层对外暴露grpc接口
// 用户服务
type UserServer struct {
	Service *userlogic.UserService
	Logger  *otelzap.Logger
	proto.UnimplementedUserServer
}

var ProtoJson = protojson.MarshalOptions{
	EmitUnpopulated: true,
}

func MethodInfoRecord(data gproto.Message) []byte {
	r, err := ProtoJson.Marshal(data)
	if err != nil {
		return []byte{}
	}
	return r
}

// 获得用户列表,可通过FliterReq过滤
func (s *UserServer) GetUserList(ctx context.Context, in *proto.UserFliterReq) (*proto.UserListRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetUserList调用,调用信息为: %s", string(info))
	res, err := s.Service.GetUserListLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetUserList失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 通过用户id获取用户信息
func (s *UserServer) GetUserById(ctx context.Context, in *proto.UserIdReq) (*proto.UserInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetUserById调用,调用信息为: %+v", info)
	res, err := s.Service.GetUserByIdLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetUserById失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 通过用户电话号码获取用户信息
func (s *UserServer) GetUserByMobile(ctx context.Context, in *proto.UserMobileReq) (*proto.UserInfoRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次GetUserByMobile调用,调用信息为: %+v", info)
	res, err := s.Service.GetUserByMobileLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用GetUserByMobile失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 创建一个用户
func (s *UserServer) CreateUser(ctx context.Context, in *proto.CreateUserReq) (*proto.CreateUserRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CreateUser调用,调用信息为: %+v", info)
	res, err := s.Service.CreateUserLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CreateUser失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 更新用户,传入的用户信息字段中无论是否为空都会完全覆盖原来的值
func (s *UserServer) AbsUpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次AbsUpdateUser调用,调用信息为: %+v", info)
	res, err := s.Service.AbsUpdateUserLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用AbsUpdateUser失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 局部更新设置了值的参数
func (s *UserServer) UpdateUser(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次UpdateUser调用,调用信息为: %+v", info)
	res, err := s.Service.UpdateUserLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用UpdateUser失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 注销一个用户
func (s *UserServer) DeleteUser(ctx context.Context, in *proto.DelUserReq) (*emptypb.Empty, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次DeleteUser调用,调用信息为: %+v", info)
	res, err := s.Service.DeleteUserLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用DeleteUser失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}

// 权限验证
func (s *UserServer) CheckUserRole(ctx context.Context, in *proto.UserPasswordReq) (*proto.UserCheckRes, error) {
	info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次CheckUserRole调用,调用信息为: %+v", info)
	res, err := s.Service.CheckUserRoleLogic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用CheckUserRole失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
}
