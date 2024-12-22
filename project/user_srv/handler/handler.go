package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"
	gb "user_srv/global"
	"user_srv/model"
	pb "user_srv/proto"

	passwd "github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

/*
	注意,web层已经过滤了一层数据,故数据是完整的
*/

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

func (us *UserServer) GetUserList(c context.Context, req *pb.UserFliterReq) (*pb.UserListRes, error) {
	logic := &model.User{}
	result, err := logic.FindByOpt(&model.FindOption{
		PagesNum: req.PagesNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	ulr := &pb.UserListRes{}
	ulr.Total = int32(result.Total)
	for _, u := range result.Data {
		i := UserToInfoRes(u)
		ulr.Data = append(ulr.Data, i)
	}
	return ulr, nil
}

func (us *UserServer) GetUserById(c context.Context, req *pb.UserIdReq) (*pb.UserInfoRes, error) {
	u := &model.User{}
	u.ID = uint(req.Id)
	if err := u.FindOneById(); err != nil {
		return nil, err
	}
	res := UserToInfoRes(u)
	return res, nil
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
	err := u.InsertOne()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (us *UserServer) UpdateUser(c context.Context, req *pb.WriteUserReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (us *UserServer) CheckUserRole(c context.Context, req *pb.UserPasswordReq) (*pb.UserCheckRes, error) {
	u := &model.User{}
	u.ID = uint(req.Id)
	result := gb.DB.Find(u)
	if result.Error != nil {
		return nil, model.ErrBadAuth
	}
	checkField := strings.Split(u.Password, "$")
	if !UnencryptPassword(&req.Password, &checkField[1], &checkField[2]) {
		return nil, model.ErrBadAuth
	}
	return &pb.UserCheckRes{
		Ok: true,
	}, nil
}

// 考虑到gorm内部采用了unlink删除,如果想延申逻辑则要么删除redis缓存要么彻底删除mysql数据
func (us *UserServer) DeleteUser(ctx context.Context, in *pb.DelUserReq) (*emptypb.Empty, error) {
	u := &model.User{}
	u.ID = uint(in.Id)
	err := u.DeleteById()
	if err != nil {
		zap.S().Errorw("用户删除失败", "msg", err.Error())
		return nil, err
	}
	return nil, nil
}
