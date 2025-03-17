package userdata

import (
	"context"
	"kenshop/pkg/common/paginate"
	"kenshop/pkg/encrypt"
	proto "kenshop/proto/user"
	model "kenshop/service/user/internal/model"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var _ UserDataService = &GormUserData{}

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

// data层为数据库交互层
// 用户服务
type GormUserData struct {
	DB *gorm.DB
}

func MustNewGormUserData(db *gorm.DB) *GormUserData {
	return &GormUserData{
		DB: db,
	}
}

func UserToResponse(u *model.User) *proto.UserInfoRes {
	return &proto.UserInfoRes{
		Id:       u.ID,
		Password: u.Password,
		Mobile:   u.Mobile,
		Role:     u.Role,
		Gender:   u.Gender,
		UserName: u.UserName,
		Birth:    u.Birth.Unix(),
	}
}

// 获得用户列表,可通过FliterReq过滤
func (s *GormUserData) GetUserListDB(ctx context.Context, in *proto.UserFliterReq) (*proto.UserListRes, error) {
	res := &proto.UserListRes{}
	data := []*model.User{}
	s.DB.Model(&model.User{}).Count(&res.Total)
	result := s.DB.Scopes(paginate.GormPaginate(int(in.PagesNum), int(in.PageSize))).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, v := range data {
		res.Data = append(res.Data, UserToResponse(v))
	}
	return res, nil
}

// 通过用户id获取用户信息
func (s *GormUserData) GetUserByIdDB(ctx context.Context, in *proto.UserIdReq) (*proto.UserInfoRes, error) {
	u := &model.User{}
	res := s.DB.Model(&model.User{}).First(u, in.Id)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return UserToResponse(u), nil
}

// 通过用户电话号码获取用户信息
func (s *GormUserData) GetUserByMobileDB(ctx context.Context, in *proto.UserMobileReq) (*proto.UserInfoRes, error) {
	u := &model.User{}
	res := s.DB.Model(&model.User{}).Where("mobile = ?", in.Mobile).Find(u)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return UserToResponse(u), nil
}

// 创建一个用户
func (s *GormUserData) CreateUserDB(ctx context.Context, in *proto.CreateUserReq) (*proto.CreateUserRes, error) {
	t := time.Unix(in.Birth, 0)
	u := &model.User{
		Password: encrypt.EncryptString(in.Password),
		Mobile:   in.Mobile,
		Gender:   in.Gender,
		UserName: in.UserName,
		Birth:    &t,
	}
	res := s.DB.Create(u)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return &proto.CreateUserRes{
		Mobile:   in.Mobile,
		Gender:   in.Gender,
		UserName: in.UserName,
		Birth:    in.Birth,
		Id:       u.ID,
	}, nil
}

// 更新用户,传入的用户信息字段中无论是否为空都会完全覆盖原来的值
func (s *GormUserData) AbsUpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	u := &model.User{
		Password: encrypt.EncryptString(in.Password),
		Role:     in.Role,
		UserName: in.UserName,
		Mobile:   in.Mobile,
		Gender:   in.Gender,
	}
	if in.Birth != 0 {
		*u.Birth = time.Unix(in.Birth, 0)
	}
	u.ID = in.Id
	res := s.DB.Save(u)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return &emptypb.Empty{}, nil
}

// 局部更新设置了值的参数
func (s *GormUserData) UpdateUserDB(ctx context.Context, in *proto.UpdateUserReq) (*emptypb.Empty, error) {
	u := &model.User{
		Password: encrypt.EncryptString(in.Password),
		Role:     in.Role,
		UserName: in.UserName,
		Mobile:   in.Mobile,
		Gender:   in.Gender,
	}
	if in.Birth != 0 {
		*u.Birth = time.Unix(in.Birth, 0)
	}
	u.ID = in.Id
	res := s.DB.Updates(u)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return &emptypb.Empty{}, nil
}

// 注销一个用户
func (s *GormUserData) DeleteUserDB(ctx context.Context, in *proto.DelUserReq) (*emptypb.Empty, error) {
	res := s.DB.Delete(&model.User{}, in.Id)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	return &emptypb.Empty{}, nil
}

// 权限验证
func (s *GormUserData) CheckUserRoleDB(ctx context.Context, in *proto.UserPasswordReq) (*proto.UserCheckRes, error) {
	u := &model.User{}
	u.ID = in.Id
	res := s.DB.Find(u)
	if res.Error != nil {
		return nil, GormErrHandle(res.Error)
	}
	ok := encrypt.UnencryptString(in.Password, u.Password)
	return &proto.UserCheckRes{
		Ok: ok,
	}, nil
}
