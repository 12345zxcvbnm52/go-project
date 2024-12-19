package handler

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"strings"
	"time"
	gb "user_srv/global"
	"user_srv/model"
	pb "user_srv/proto"

	passwd "github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedUserServer
}

func EncryptPassword(pd string) string {
	options := &passwd.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, codePwd := passwd.Encode(pd, options)
	passWord := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, codePwd)
	return passWord
}

func UnencryptPassword(raw, salt, encryptPassword *string) bool {
	opt := &passwd.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	return passwd.Verify(*raw, *salt, *encryptPassword, opt)
}

func UserToInfoRes(u *model.User) *pb.UserInfoRes {
	return &pb.UserInfoRes{
		UserName: u.UserName,
		Mobile:   u.Mobile,
		Id:       uint32(u.ID),
		Role:     u.Role,
		Password: u.Password,
		Birth:    u.Birth.Unix(),
	}
}

func InfoResToUser(u *pb.UserInfoRes) *model.User {
	birth := time.Unix(u.Birth, 0)
	return &model.User{
		UserName: u.UserName,
		Mobile:   u.Mobile,
		Model:    model.Model{ID: uint(u.Id)},
		Role:     u.Role,
		Password: u.Password,
		Birth:    &birth,
	}

}

// gorm给出的分页函数的最佳实践
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (us *UserServer) GetUserList(c context.Context, req *pb.UserFliterReq) (*pb.UserListRes, error) {
	var users []model.User = make([]model.User, req.PageSize)
	result := gb.DB.Scopes(Paginate(int(req.PagesNum), int(req.PageSize))).Find(users)
	if result.Error != nil {
		zap.S().Errorw("获得UserList失败", "msg", result.Error)
		return nil, result.Error
	}
	ulr := &pb.UserListRes{}

	ulr.Total = int32(result.RowsAffected)
	for _, u := range users {
		i := UserToInfoRes(&u)
		ulr.Data = append(ulr.Data, i)
	}
	return ulr, nil
}

func (us *UserServer) GetUserById(c context.Context, req *pb.UserIdReq) (*pb.UserInfoRes, error) {
	return nil, nil
}

func (us *UserServer) GetUserByMobile(c context.Context, req *pb.UserMobileReq) (*pb.UserInfoRes, error) {
	return nil, nil
}

func (us *UserServer) CreateUser(c context.Context, req *pb.WriteUserReq) (*pb.UserInfoRes, error) {
	res := &pb.UserInfoRes{
		Password: EncryptPassword(req.Password),
		UserName: req.UserName,
		Gender:   req.Gender,
		Birth:    req.Birth,
		Role:     req.Role,
		Mobile:   req.Mobile,
	}
	u := InfoResToUser(res)
	result := gb.DB.Create(u)
	if result.Error != nil {
		zap.S().Errorw("UserCreate失败", "msg", result.Error)
		return nil, result.Error
	}
	return res, nil
}

func (us *UserServer) UpdateUser(c context.Context, req *pb.WriteUserReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (us *UserServer) CheckUserRole(c context.Context, req *pb.UserPasswordReq) (*pb.UserCheckRes, error) {
	zap.S().Infoln("someone call once")
	u := &model.User{}
	u.ID = uint(req.Id)
	result := gb.DB.Find(u)
	if result.Error != nil {
		return nil, errors.New("用户名或密码错误")
	}
	checkField := strings.Split(u.Password, "$")
	if !UnencryptPassword(&req.Password, &checkField[1], &checkField[2]) {
		return nil, errors.New("Password错误")
	}
	return &pb.UserCheckRes{
		Ok: true,
	}, nil
}
