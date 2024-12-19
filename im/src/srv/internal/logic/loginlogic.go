package logic

import (
	"context"

	__proto "srv/internal/proto"
	"srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *__proto.LoginReq) (*__proto.LoginRes, error) {
	// todo: add your logic here and delete this line

	return &__proto.LoginRes{}, nil
}
