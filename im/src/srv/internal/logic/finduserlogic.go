package logic

import (
	"context"

	__proto "srv/internal/proto"
	"srv/internal/svc"
	"srv/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *__proto.FindUserReq) (*__proto.UserInfoRes, error) {
	res := []*model.Users{}

	return &__proto.UserInfoRes{}, nil
}
