package handler

import (
	"context"
	"net/http"
	"order_web/form"
	gb "order_web/global"
	pb "order_web/proto"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateOrder(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "gin-ctx", c)
	client, err := gb.RpcPool.Value()
	if err != nil {
		zap.S().Errorw("池内获取的连接不可用", "msg", err.Error())
		RpcErrorHandle(c, status.Error(codes.Internal, ""))
		c.Abort()
		return
	}
	defer client.Close()
	cc := pb.NewOrderClient(client)

	form := &form.OrderCreateForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		ValidatorErrorHandle(c, err)
		c.Abort()
		return
	}

	res, err := cc.CreateOrder(ctx, &pb.OrderInfoReq{
		UserId:       form.UserId,
		Address:      form.Address,
		SignerName:   form.SignerName,
		SignerMobile: form.SignerMobile,
		Message:      form.Message,
		PayWay:       form.PayWay,
	})
	if err != nil {
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}
	//后续支付还得做成一个函数
	c.JSON(http.StatusOK, gin.H{
		"msg":  "订单创建成功,请在15min内支付",
		"data": res,
	})
}
