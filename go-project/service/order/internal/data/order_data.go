package orderdata

import (
	"context"
	"encoding/json"
	"kenshop/pkg/errors"
	"kenshop/pkg/rockmq"
	gproto "kenshop/proto/goods"
	iproto "kenshop/proto/inventory"
	proto "kenshop/proto/order"
	model "kenshop/service/order/internal/model"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// OrderDataService是提供Order底层相关数据操作的接口
type OrderDataService interface {
	//获得用户购物车信息
	GetUserCartItemsDB(ctx context.Context, in *proto.UserInfoReq) (*proto.CartItemListRes, error)
	//为购物车添加商品
	CreateCartItemDB(ctx context.Context, in *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error)
	//修改购物车的一条记录
	UpdateCartItemDB(ctx context.Context, in *proto.UpdateCartItemReq) (*emptypb.Empty, error)

	DeleteCartItemDB(ctx context.Context, in *proto.DelCartItemReq) (*emptypb.Empty, error)

	CreateOrderDB(ctx context.Context, in *proto.CreateOrderReq) (*proto.OrderInfoRes, error)

	GetOrderListDB(ctx context.Context, in *proto.OrderFliterReq) (*proto.OrderListRes, error)

	GetOrderInfoDB(ctx context.Context, in *proto.OrderInfoReq) (*proto.OrderDetailRes, error)

	UpdateOrderStatusDB(ctx context.Context, in *proto.OrderStatusReq) (*emptypb.Empty, error)
}

// Order服务中的Data层,是数据操作的具体逻辑
type GormOrderData struct {
	DB                 *gorm.DB
	Logger             *otelzap.Logger
	OrderTransProducer *rockmq.TransProducer
	GoodsRpcCli        gproto.GoodsClient
	InventoryRpcCli    iproto.InventoryClient
	RebackTopic        string
}

type OrderListener struct {
	DB                 *gorm.DB
	Logger             *otelzap.Logger
	TimeoutRebackTopic string
	RebackTopic        string
	GoodsRpcCli        gproto.GoodsClient
	InventoryRpcCli    iproto.InventoryClient
	TimeoutProducer    rocketmq.Producer
	Ctx                context.Context
}

func MustNewOrderListener(db *gorm.DB, l *otelzap.Logger, t string, r string,
	ir iproto.InventoryClient, gr gproto.GoodsClient, p rocketmq.Producer) *OrderListener {
	return &OrderListener{DB: db, Logger: l, GoodsRpcCli: gr, RebackTopic: r,
		InventoryRpcCli: ir, TimeoutProducer: p, TimeoutRebackTopic: t}
}

var _ OrderDataService = (*GormOrderData)(nil)

func MustNewGormOrderData(db *gorm.DB, l *otelzap.Logger, r string,
	ir iproto.InventoryClient, gr gproto.GoodsClient, p *rockmq.TransProducer) *GormOrderData {
	return &GormOrderData{DB: db, Logger: l, GoodsRpcCli: gr, InventoryRpcCli: ir, OrderTransProducer: p, RebackTopic: r}
}

// 获得用户购物车信息
func (s *GormOrderData) GetUserCartItemsDB(ctx context.Context, in *proto.UserInfoReq) (*proto.CartItemListRes, error) {
	return nil, errors.New("this method is not implemented")
}

// 为购物车添加商品
func (s *GormOrderData) CreateCartItemDB(ctx context.Context, in *proto.CreateCartItemReq) (*proto.CartItemInfoRes, error) {
	return nil, errors.New("this method is not implemented")
}

// 修改购物车的一条记录
func (s *GormOrderData) UpdateCartItemDB(ctx context.Context, in *proto.UpdateCartItemReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormOrderData) DeleteCartItemDB(ctx context.Context, in *proto.DelCartItemReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *OrderListener) OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		order := &model.Order{}
		order.OrderSign = string(msg.Body)

		if res := s.DB.Model(&model.Order{}).Where("order_sign = ?", order.OrderSign).First(order); res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return consumer.ConsumeSuccess, nil
			}
			return consumer.ConsumeRetryLater, res.Error
		}
		//如果状态小于等于订单未支付就归还
		if order.Status < model.StatusPaid {
			tx := s.DB.Begin()
			defer func() {
				if err := recover(); err != nil {
					tx.Rollback()
					s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
					return
				}
			}()
			order.Status = model.StatusCancelled
			res := tx.Model(&model.Order{}).Where("order_sign = ?", order.OrderSign).Update("status", model.StatusCancelled)
			if res.Error != nil {
				s.Logger.Sugar().Errorf("订单状态更新失败 err = %s", res.Error.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, res.Error
			}

			orderGoods := []*model.OrderGoods{}
			res = tx.Model(&model.OrderGoods{}).Where("order_id = ?", order.ID).Find(&orderGoods)
			if res.Error != nil {
				s.Logger.Sugar().Errorf("订单内商品查找失败 err = %s", res.Error.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, res.Error
			}

			if res.RowsAffected == 0 {
				return consumer.ConsumeSuccess, nil
			}
			//这里暂时只考虑使用gorm软删除下的购物车恢复功能
			goodsIds := []uint32{}
			for _, v := range orderGoods {
				goodsIds = append(goodsIds, v.GoodsId)
			}
			r := tx.Exec("update carts set deleted_at = null where goods_id in (?)", goodsIds)
			if r.Error != nil {
				tx.Rollback()
				zap.S().Errorw("购物车恢复失败", "msg", r.Error.Error())
				return consumer.ConsumeRetryLater, r.Error
			}
			if r.RowsAffected == 0 {
				tx.Rollback()
				zap.S().Errorw("购物车恢复失败", "msg", "不存在要恢复的数据")
				return consumer.ConsumeSuccess, nil
			}

			_, err := s.TimeoutProducer.SendSync(ctx, primitive.NewMessage(s.RebackTopic, msg.Body))
			if err != nil {
				tx.Rollback()
				zap.S().Errorw("rockmq中producer发送事务消息失败", "msg", err.Error())
				return consumer.ConsumeRetryLater, err
			}
			tx.Commit()
		}
	}
	return consumer.ConsumeSuccess, nil
}

func (s *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func (s *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	carts := []*model.Cart{}
	order := &model.Order{}
	json.Unmarshal(msg.Body, order)

	//先检查一下存不存在该订单,存在则直接退出,实现部分幂等性(但是仍有可能因为并发出现多次调用都进入逻辑,需要在后面检查)
	res := s.DB.Model(&model.Order{}).Where("order_sign = ?", order.OrderSign).First(order)
	if (res.Error != nil && res.Error != gorm.ErrRecordNotFound) || res.Error == nil {
		s.Logger.Sugar().Warnf("幂等性检查: err = %+v", res.Error)
		return primitive.RollbackMessageState
	}

	//通过购物车找到对应的商品信息
	res = s.DB.Model(&model.Cart{}).Where("user_id = ? and selected = true", order.UserId).Find(&carts)
	if res.Error != nil {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(res.Error))))
		return primitive.RollbackMessageState
	}
	//TODO
	if res.RowsAffected == 0 {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(gorm.ErrRecordNotFound))))
	}

	goodsIds := []uint32{}
	//这个map存goodsId对应的goodsNum,以便于后续到goods_srv中查找价格
	goodsNumMap := make(map[uint32]int32)
	for _, v := range carts {
		goodsIds = append(goodsIds, v.GoodsId)
		goodsNumMap[v.GoodsId] = v.GoodsNums
	}

	//从商品服务中获得商品信息,注意,前面这部分只涉及到数据查询,不用开启事务
	goodsInfo, err := s.GoodsRpcCli.GetGoodsListById(s.Ctx, &gproto.GoodsIdsReq{Ids: goodsIds})
	if err != nil {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(err))))
		return primitive.RollbackMessageState
	}
	//准备记录要写到orderGoods表的记录和扣减库存记录
	//注意这里的orderGoods还没有填充orderId
	orderGoods := []*model.OrderGoods{}
	decrStock := &iproto.UpdateStockReq{}
	for _, v := range goodsInfo.Data {
		//记录总的花费
		order.Cost += v.SalePrice * float32(goodsNumMap[v.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			GoodsId:    v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.FirstImage,
			GoodsPrice: v.SalePrice,
			GoodsNum:   goodsNumMap[v.Id],
		})
		decrStock.DecrData = append(decrStock.DecrData, &iproto.UpdateInventoryReq{
			GoodsId:  v.Id,
			GoodsNum: goodsNumMap[v.Id],
		})
	}
	decrStock.OrderSign = order.OrderSign

	//TODO 一旦扣减完成就要开启事务,这里出现err可能还得考虑真的是否没有扣减,防止出现问题
	// 可以考虑的实现方式有数据库表标志位或缓存检查,
	if _, err := s.InventoryRpcCli.DecrStock(s.Ctx, decrStock); err != nil {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(err))))
		return primitive.RollbackMessageState
	}

	tx := s.DB.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			s.Logger.Sugar().Errorf("gorm遇到异常无法正常执行, err=%+v", err)
			return
		}
	}()

	//考虑到order_sign是唯一标志,直接创建是不用考虑幂等的
	if res := tx.Model(&model.Order{}).Omit("pay_time").Create(order); res.Error != nil {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(res.Error))))
		tx.Rollback()
		return primitive.CommitMessageState
	}

	for i := range orderGoods {
		orderGoods[i].OrderId = order.ID
	}

	//把所有orderGoods插入到SQL中,失败的rollback动作会在Insert里执行
	if err := tx.Model(&model.OrderGoods{}).CreateInBatches(&orderGoods, 100).Error; err != nil {
		tx.Rollback()
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(err))))
		return primitive.CommitMessageState
	}

	// 删除购物车内已经被购买的商品,这里考虑到只有有选中的购物车商品才能创建订单,故无需判断是不是不存在select的订单
	// TODO 注意这里暂时用的gorm软删除,像购物车这种应当考虑完全删除
	res = tx.Model(&model.Cart{}).Where("selected = ? and user_id = ?", true, order.UserId).Delete(&model.Cart{})
	if res.Error != nil {
		tx.Rollback()
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(res.Error))))
		return primitive.CommitMessageState
	}

	//把订单号消息发送到库存服务中即可,具体超时归还逻辑由那边处理即可
	timeoutMsg := primitive.NewMessage(s.TimeoutRebackTopic, []byte(order.OrderSign))
	timeoutMsg.WithDelayTimeLevel(4)
	_, err = s.TimeoutProducer.SendSync(s.Ctx, timeoutMsg)
	if err != nil {
		msg.WithProperty("error", string(errors.MarshalCodeError(GormOrderErrHandle(err))))
		tx.Rollback()
		return primitive.CommitMessageState
	}
	tx.Commit()
	//可以考虑如果整个函数在这里失败(宕机)会发生什么
	return primitive.CommitMessageState
}

func (s *GormOrderData) CreateOrderDB(ctx context.Context, in *proto.CreateOrderReq) (*proto.OrderInfoRes, error) {
	order := &model.Order{
		UserId:       in.UserId,
		Address:      in.Address,
		SignerName:   in.SignerName,
		SignerMobile: in.SignerMobile,
		Status:       model.StatusCreated,
		PayWay:       in.PayWay,
		Message:      in.Message,
		OrderSign:    in.OrderSign,
	}

	msgBody, _ := json.Marshal(order)
	msg := primitive.NewMessage(s.RebackTopic, msgBody)
	_, err := s.OrderTransProducer.SendMessageInTransaction(ctx, msg)

	if err != nil {
		return nil, errors.WithCoder(err, errors.CodeBadRockmq, "")
	}
	errMsg := msg.GetProperty("error")
	if err = errors.UnmarshalCodeError(errMsg); err != nil {
		return nil, err
	}

	orderMsg := msg.GetProperty("order")
	json.Unmarshal([]byte(orderMsg), order)
	return OrderToOrderInfoRes(order), nil
}

func (s *GormOrderData) GetOrderListDB(ctx context.Context, in *proto.OrderFliterReq) (*proto.OrderListRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormOrderData) GetOrderInfoDB(ctx context.Context, in *proto.OrderInfoReq) (*proto.OrderDetailRes, error) {
	return nil, errors.New("this method is not implemented")
}

func (s *GormOrderData) UpdateOrderStatusDB(ctx context.Context, in *proto.OrderStatusReq) (*emptypb.Empty, error) {
	return nil, errors.New("this method is not implemented")
}

func OrderToOrderInfoRes(c *model.Order) *proto.OrderInfoRes {
	return &proto.OrderInfoRes{
		Id:           c.ID,
		UserId:       c.UserId,
		OrderSign:    c.OrderSign,
		Status:       int32(c.Status),
		PayWay:       c.PayWay,
		SignerName:   c.SignerName,
		SignerMobile: c.SignerMobile,
		Address:      c.Address,
		Cost:         c.Cost,
		Message:      c.Message,
	}
}
