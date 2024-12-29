package model

import (
	"context"
	"encoding/json"
	"fmt"
	gb "order_srv/global"
	pb "order_srv/proto"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
)

type OrderListener struct {
	Err  error
	ID   uint32
	Cost float32
	Ctx  context.Context
}

func (ol *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	//通过购物车找到对应的商品信息
	cartLogic := &Cart{}
	goodsIds := []uint32{}
	order := &Order{}
	json.Unmarshal(msg.Body, order)

	goods, err := cartLogic.FindByOpt(&CartFindOption{UserId: order.UserId, Selected: true})
	if err != nil {
		if err == ErrCartNotFound {
			ol.Err = ErrCartNoSelected
		} else {
			ol.Err = ErrInternalWrong
		}
		return primitive.RollbackMessageState
	}

	//这个map存goodsId对应的goodsNum,以便于后续到goods_srv中查找价格
	goodsNumMap := make(map[uint32]int32)
	for _, v := range goods.Data {
		goodsIds = append(goodsIds, v.GoodsId)
		goodsNumMap[v.GoodsId] = v.GoodsNums
	}

	//跨微服务调用
	goodsConn, err := gb.DefaultDial(gb.ServerConfig.GoodsServerName)
	if err != nil {
		zap.S().Errorw("商品服务连接失败", "msg", err.Error())
		ol.Err = ErrBadGoodsClient
		return primitive.RollbackMessageState
	}
	inventoryConn, err := gb.DefaultDial(gb.ServerConfig.InventoryServerName)
	if err != nil {
		zap.S().Errorw("库存服务连接失败", "msg", err.Error())
		ol.Err = ErrBadInventoryClient
		return primitive.RollbackMessageState
	}
	goodsClient := pb.NewGoodsClient(goodsConn)
	goodsInfo, err := goodsClient.GetGoodsListById(context.Background(), &pb.GoodsIdsReq{Id: goodsIds})
	if err != nil {
		zap.S().Errorw("跨微服务获得商品信息失败", "msg", err.Error())
		ol.Err = ErrBadGoodsClient
		return primitive.RollbackMessageState
	}

	//注意这里的orderGoods还没有填充orderId
	orderGoods := []*OrderGoods{}
	decrStock := []*pb.StockInfoReq{}
	for _, v := range goodsInfo.Data {
		order.Cost += v.SalePrice * float32(goodsNumMap[v.Id])
		orderGoods = append(orderGoods, &OrderGoods{
			GoodsId:    v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.FirstImage,
			GoodsPrice: v.SalePrice,
			GoodsNum:   goodsNumMap[v.Id],
		})
		decrStock = append(decrStock, &pb.StockInfoReq{
			GoodsId:  v.Id,
			GoodsNum: goodsNumMap[v.Id],
		})
	}

	inventoryClient := pb.NewInventoryClient(inventoryConn)
	if _, err := inventoryClient.DecrStock(context.Background(), &pb.DecrStockReq{DecrData: decrStock}); err != nil {
		zap.S().Errorw("跨微服务扣减失败", "msg", err.Error())
		//要考虑网络的问题(如出现误判,哪怕扣减成功,但因为网络的超时等导致返回的err是非nil),但是不知道怎么写好,先搁置在这
		//先简化模型,如果出错,则一定是扣减失败
		ol.Err = err
		return primitive.RollbackMessageState
	}

	tx := gb.DB.Begin()
	if res := tx.Create(order); res.Error != nil {
		tx.Rollback()
		zap.S().Errorw("创建订单失败", "msg", res.Error.Error())
		if res.Error == ErrOrderDuplicated {
			ol.Err = ErrOrderDuplicated
		}
		ol.Err = ErrOrderFailedCreate
		return primitive.CommitMessageState
	}
	ol.Cost = order.Cost
	ol.ID = order.ID
	//填充所有orderGoods的OrderId
	for i := range orderGoods {
		orderGoods[i].OrderId = order.ID
	}

	//把所有orderGoods插入到SQL中,失败的rollback动作会在Insert里执行
	orderGoodsLogic := &OrderGoods{}
	if err := orderGoodsLogic.Insert(orderGoods, tx); err != nil {
		ol.Err = err
		return primitive.CommitMessageState
	}
	//删除购物车内已经被购买的商品
	if res := tx.Model(&Cart{}).Where("selected = ", true).Delete(&Cart{}); res.RowsAffected == 0 || res.Error != nil {
		zap.S().Errorw("购物车清除失败", "msg", res.Error.Error())
		tx.Rollback()
		ol.Err = ErrInternalWrong
		return primitive.CommitMessageState
	}

	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMqConfig.Host, gb.ServerConfig.RockMqConfig.Port)}),
	)
	if err != nil {
		zap.S().Errorw("消息队列生成生产者失败", "msg", err.Error())
		tx.Rollback()
		ol.Err = ErrBadRockMq
		return primitive.CommitMessageState
	}
	if err = p.Start(); err != nil {
		zap.S().Errorw("消息队列生产者运行失败", "msg", err.Error())
		tx.Rollback()
		ol.Err = ErrBadRockMq
		return primitive.CommitMessageState
	}
	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {

		zap.S().Errorw("发送延时消息失败", "msg", err.Error())
		tx.Rollback()
		ol.Err = ErrBadRockMq
		return primitive.CommitMessageState
	}
	tx.Commit()
	ol.Err = nil
	return primitive.RollbackMessageState
}

// 我能想到的合理的检查是否扣减(最主要的是订单服务在扣减成功,订单事务创建前宕机难以检查)可以在库存服务内加一个表,把扣减的信息保存到表
// 后续查这张表即可,根据订单编号的唯一性和幂等性查询,这样便能很好保证整个过程的一致性

func (ol *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	order := &Order{}
	json.Unmarshal(msg.Body, order)
	if res := gb.DB.Where("order_sign = ?", order.OrderSign).First(&order); res.RowsAffected == 0 {
		return primitive.CommitMessageState //你并不能说明这里就是库存已经扣减了
	}
	return primitive.RollbackMessageState
}
