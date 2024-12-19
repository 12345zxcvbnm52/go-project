package logic

import (
	"context"

	__proto "srv/internal/proto"
	"srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *__proto.Req) (*__proto.Res, error) {
	// todo: add your logic here and delete this line

	return &__proto.Res{}, nil
}
