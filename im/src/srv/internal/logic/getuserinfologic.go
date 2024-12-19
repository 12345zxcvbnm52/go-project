package logic

import (
	"context"

	__proto "srv/internal/proto"
	"srv/internal/svc"
	"srv/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *__proto.UserInfoReq) (*__proto.UserInfoRes, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, ErrUserNotFind
		}
		return nil, err
	}

	return &__proto.UserInfoRes{
		UserName: u.Username,
		Id:       u.Id,
		Avatar:   u.Avatar,
		Mobile:   u.Mobile,
		Status:   int32(u.Status),
		Gender:   u.Gender,
	}, nil
}
