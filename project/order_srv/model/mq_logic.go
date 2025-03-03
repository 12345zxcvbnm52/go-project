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
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type OrderListener struct {
	Err error

	OrderId uint32
	//总的花费,调用库存服务时在本地事务中计算得到
	Cost float32
	//用于链路追踪,与业务逻辑无关
	Ctx context.Context
}

func (ol *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	cartLogic := &Cart{}
	goodsIds := []uint32{}
	order := &Order{}
	json.Unmarshal(msg.Body, order)
	fspan := opentracing.SpanFromContext(ol.Ctx)

	//通过购物车找到对应的商品信息
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

	//跨微服务调用,通过商品服务获得价格,通过库存服务扣减库存,
	//在本地事务开启前,只要出错就rollback库存归还半消息
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

	//获得商品价格信息
	cspan := opentracing.GlobalTracer().StartSpan("ordersrv调用goodssrv", opentracing.ChildOf(fspan.Context()))
	goodsClient := pb.NewGoodsClient(goodsConn)
	goodsInfo, err := goodsClient.GetGoodsListById(context.Background(), &pb.GoodsIdsReq{Id: goodsIds})
	if err != nil {
		zap.S().Errorw("跨微服务获得商品信息失败", "msg", err.Error())
		ol.Err = ErrBadGoodsClient
		return primitive.RollbackMessageState
	}
	cspan.Finish()

	//准备记录要写到orderGoods表的记录和扣减库存记录
	//注意这里的orderGoods还没有填充orderId
	orderGoods := []*OrderGoods{}
	decrStock := []*pb.WriteInvtReq{}
	for _, v := range goodsInfo.Data {
		order.Cost += v.SalePrice * float32(goodsNumMap[v.Id])
		orderGoods = append(orderGoods, &OrderGoods{
			GoodsId:    v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.FirstImage,
			GoodsPrice: v.SalePrice,
			GoodsNum:   goodsNumMap[v.Id],
		})
		decrStock = append(decrStock, &pb.WriteInvtReq{
			GoodsId:  v.Id,
			GoodsNum: goodsNumMap[v.Id],
		})
	}

	//要考虑网络的问题(如出现误判,哪怕扣减成功,但因为网络的超时等导致返回的err是非nil),但是不知道怎么写好,先搁置在这
	//先简化模型,如果出错,则一定是扣减失败
	cspan = opentracing.GlobalTracer().StartSpan("ordersrv调用inventorysrv", opentracing.ChildOf(fspan.Context()))
	inventoryClient := pb.NewInventoryClient(inventoryConn)
	if _, err := inventoryClient.DecrStock(
		context.Background(),
		&pb.DecrStockReq{DecrData: decrStock, OrderSign: order.OrderSign},
	); err != nil {
		zap.S().Errorw("跨微服务扣减失败", "msg", err.Error())
		ol.Err = err
		return primitive.RollbackMessageState
	}
	cspan.Finish()
	//这里如果库存扣减完成,但是服务宕机在后续出错后commit或rollback前是只能靠回查来实现reback的
	//要注意的是,上述这种情况本地事务是安全的

	//事务开启后出了错就一定要发送出库存归还的半消息,补上上面的库存扣减
	cspan = opentracing.GlobalTracer().StartSpan("localTrascation", opentracing.ChildOf(fspan.Context()))
	tx := gb.DB.Begin()
	if res := tx.Model(&Order{}).Omit("pay_time").Create(order); res.Error != nil {
		tx.Rollback()
		zap.S().Errorw("创建订单失败", "msg", res.Error.Error())
		ol.Err = ErrOrderFailedCreate
		if res.Error == ErrOrderDuplicated {
			ol.Err = ErrOrderDuplicated
		}
		return primitive.CommitMessageState
	}
	ol.OrderId = order.ID
	ol.Cost = order.Cost
	//填充所有orderGoods的OrderId
	for i := range orderGoods {
		orderGoods[i].OrderId = order.ID
	}

	//把所有orderGoods插入到SQL中,失败的rollback动作会在Insert里执行
	orderGoodsLogic := &OrderGoods{}
	if err := orderGoodsLogic.Insert(orderGoods, tx); err != nil {
		tx.Rollback()
		ol.Err = err
		return primitive.CommitMessageState
	}
	//删除购物车内已经被购买的商品,这里考虑到只有有选中的购物车商品才能创建订单,故无需判断是不是不存在select的订单
	if res := tx.Model(&Cart{}).Where("selected = ? and user_id = ?", true, order.UserId).Delete(&Cart{}); res.RowsAffected == 0 || res.Error != nil {
		zap.S().Errorw("购物车清除失败", "msg", res.Error.Error())
		tx.Rollback()

		ol.Err = ErrInternalWrong
		return primitive.CommitMessageState
	}

	//发送延时消息用于延迟归还,这里考虑创建一个新的消息队列
	//go的逻辑是默认按照进程id创建,这样同一个进程默认无法生成第二个消息队列,暂时未做
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", gb.ServerConfig.RockMq.Host, gb.ServerConfig.RockMq.Port)}),
		producer.WithGroupName(order.OrderSign+"-"+gb.ServerConfig.RockMq.TimeoutTopic),
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
	defer p.Shutdown()
	//发送延时消息,话说不确定不同组的producer是不是可以放心shutdown
	//这时候发过去的是只有订单号的order
	msg = primitive.NewMessage(gb.ServerConfig.RockMq.TimeoutTopic, msg.Body)
	msg.WithDelayTimeLevel(6)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorw("发送延时消息失败", "msg", err.Error())
		tx.Rollback()
		ol.Err = ErrBadRockMq
		return primitive.CommitMessageState
	}
	zap.S().Infoln("延时消息发送成功")
	tx.Commit()
	cspan.Finish()
	ol.Err = nil
	return primitive.RollbackMessageState
}

// 我能想到的合理的检查是否扣减(最主要的是订单服务在扣减成功,订单事务创建前宕机难以检查)可以在库存服务内加一个表,把扣减的信息保存到表
// 后续查这张表即可,根据订单编号的唯一性和幂等性查询,这样便能很好保证整个过程的一致性
func (ol *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	order := &Order{}
	zap.S().Errorw("错在这")
	json.Unmarshal(msg.Body, order)
	//这里可能要跨微服务调用,调用库存服务得到record看库存是否处于待扣减还是已扣减状态
	if res := gb.DB.Model(&Order{}).Where("order_sign = ?", order.OrderSign).First(order); res.RowsAffected == 0 {
		zap.S().Errorw("错在这")
		return primitive.CommitMessageState //你并不能说明这里就是库存已经扣减了
	}
	zap.S().Errorw("错在这")
	return primitive.RollbackMessageState
}
