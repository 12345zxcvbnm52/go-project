package logic

import (
	"context"
	"time"

	"srv/internal/pkg"
	__proto "srv/internal/proto"
	"srv/internal/svc"
	"srv/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *__proto.RegisterReq) (*__proto.RegisterRes, error) {
	userInfo, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}
	if userInfo != nil {
		return nil, ErrMobileRegisted
	}

	id := pkg.GenUid(l.svcCtx.Config.MysqlConf)
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.Users{
		Id:       id,
		Username: in.UserName,
		Gender:   in.Gender,
		Mobile:   in.Mobile,
		Avatar:   in.Avatar,
		Password: pkg.EncryptPassword(in.Password),
	})
	if err != nil {
		return nil, err
	}
	token, err := pkg.JWTAuthCreate(pkg.CustomClaims{ID: id, UserName: in.UserName}, l.svcCtx.Config.Jwt.AccessSecret)
	if err != nil {
		return nil, err
	}

	return &__proto.RegisterRes{Token: token, Expire: time.Now().Unix() + 1*24*60*60}, nil
}
