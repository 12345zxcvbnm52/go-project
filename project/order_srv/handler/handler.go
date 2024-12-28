package handler

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	gb "order_srv/global"
	"order_srv/model"
	pb "order_srv/proto"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

func CartToCartItemRes(c *model.Cart) *pb.CartItemInfoRes {
	return &pb.CartItemInfoRes{
		GoodsId:  c.GoodsId,
		GoodsNum: c.GoodsNums,
		Id:       c.ID,
		Selected: c.Selected,
	}
}

func OrderToOrderInfoRes(c *model.Order) *pb.OrderInfoRes {
	return &pb.OrderInfoRes{
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

func OrderGoodsToItemRes(c *model.OrderGoods) *pb.OrderItemRes {
	return &pb.OrderItemRes{
		Id:          c.ID,
		OrderId:     c.OrderId,
		GoodsId:     c.GoodsId,
		GoodsNum:    c.GoodsNum,
		GoodsName:   c.GoodsName,
		GoodsImages: c.GoodsImage,
		GoodsPrice:  c.GoodsPrice,
	}
}

func (us *OrderServer) GetCartInfoList(ctx context.Context, req *pb.UserInfoReq) (*pb.CartItemListRes, error) {
	logic := &model.Cart{}
	res, err := logic.FindByUserId(req.UserId)
	if err != nil {
		return nil, err
	}
	r := &pb.CartItemListRes{
		Total: res.Total,
	}
	for _, v := range res.Data {
		v.UserId = req.UserId
		r.Data = append(r.Data, CartToCartItemRes(v))
	}
	return r, nil
}

func (us *OrderServer) CreateCartItem(ctx context.Context, req *pb.CartItemInfoReq) (*pb.CartItemInfoRes, error) {
	u := &model.Cart{
		UserId:  req.UserId,
		GoodsId: req.GoodsId,
	}
	if err := u.InsertOne(); err != nil {
		return nil, err
	}
	return CartToCartItemRes(u), nil
}

func (us *OrderServer) DeleteCartItem(ctx context.Context, req *pb.DelCartItemReq) (*emptypb.Empty, error) {
	u := &model.Cart{}
	u.GoodsId = req.GoodsId
	u.ID = req.Id
	u.UserId = req.UserId
	if err := u.DeleteOneByUserGoodsIds(); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (us *OrderServer) GetOrderList(ctx context.Context, req *pb.OrderFliterReq) (*pb.OrderListRes, error) {
	logic := &model.Order{}
	res, err := logic.FindByOpt(&model.OrderFindOption{PagesNum: req.PagesNum, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	if res.Total == 0 {
		return nil, model.ErrOrderNotFound
	}
	r := &pb.OrderListRes{
		Total: res.Total,
	}
	for _, v := range res.Data {
		r.Data = append(r.Data, OrderToOrderInfoRes(v))
	}
	return r, nil
}

func (us *OrderServer) GetOrderInfo(ctx context.Context, req *pb.OrderInfoReq) (*pb.OrderDetailRes, error) {
	u := &model.Order{}
	u.ID = req.Id
	u.UserId = req.UserId
	if err := u.FindOne(); err != nil {
		return nil, err
	}
	res := &pb.OrderDetailRes{
		Id:           u.ID,
		UserId:       u.UserId,
		OrderSign:    u.OrderSign,
		Status:       int32(u.Status),
		PayWay:       u.PayWay,
		SignerName:   u.SignerName,
		SignerMobile: u.SignerMobile,
		Address:      u.Address,
		Cost:         u.Cost,
		Message:      u.Message,
	}
	//在逻辑上是能保证查找的订单是一定有相应的商品的
	logic := &model.OrderGoods{}
	if r, err := logic.FindByOrderId(u.UserId); err != nil {
		return nil, err
	} else {
		for _, v := range r.Data {
			res.Items = append(res.Items, OrderGoodsToItemRes(v))
		}
	}
	return res, nil
}

func GenOrderSign(userId uint32) string {
	now := time.Now()
	s := now.Format(time.DateTime) + fmt.Sprintf("%d", userId) + "@"
	if 32-len(s) < 0 {
		zap.S().Infow("记录一次超出意料的信息")
	}
	f := math.Pow(10, float64(32-len(s)))
	return s + fmt.Sprintf("%d", rand.IntN(int(f)))
}

// 新建订单需要计算商品的金额,不能由前端传入价格数据(防止爬虫等),
// 顾需要1.根据购物车商品跨微服务调用查询金额,2.跨微服务扣减库存,3生成订单,填充本服务的几个表
// 其中详细:根据购物车表获得商品的信息并写入OrderGoods表,并把信息同步到Order表中
func (us *OrderServer) CreateOrder(ctx context.Context, req *pb.OrderInfoReq) (*pb.OrderInfoRes, error) {
	cartLogic := &model.Cart{}
	goods, err := cartLogic.FindByOpt(&model.CartFindOption{UserId: req.UserId, Selected: true})
	if err != nil {
		if err == model.ErrCartNotFound {
			return nil, model.ErrCartNoSelected
		}
		return nil, err
	}
	goodsIds := []int32{}
	//这个map存goodsId对应的goodsNum
	goodsNumMap := make(map[int32]int32)
	for _, v := range goods.Data {
		goodsIds = append(goodsIds, v.GoodsId)
		goodsNumMap[v.GoodsId] = v.GoodsNums
	}

	res := &model.Order{}

	//跨微服务调用
	goodsConn, err := gb.DefaultDial(gb.ServerConfig.GoodsServerName)
	if err != nil {
		return nil, model.ErrInternalWrong
	}
	inventoryConn, err := gb.DefaultDial(gb.ServerConfig.InventoryServerName)
	if err != nil {
		return nil, model.ErrInternalWrong
	}
	goodsClient := pb.NewGoodsClient(goodsConn)
	goodsInfo, err := goodsClient.GetGoodsListById(context.Background(), &pb.BatchGoodsByIdReq{Id: goodsIds})
	if err != nil {
		return nil, model.ErrInternalWrong
	}
	//注意这里的orderGoods还没有填充orderId
	orderGoods := []*model.OrderGoods{}
	decrStock := []*pb.StockInfoReq{}
	for _, v := range goodsInfo.Data {
		res.Cost += v.SalePrice * float32(goodsNumMap[v.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
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
	inventory := pb.NewInventoryClient(inventoryConn)
	//这里需要在错误中更详细指明哪几个商品库存不足
	if _, err := inventory.DecrStock(context.Background(), &pb.DecrStockReq{DecrData: decrStock}); err != nil {
		return nil, err
	}
	res.Address = req.Address
	res.Message = req.Message
	res.SignerMobile = req.SignerMobile
	res.SignerName = req.SignerName
	res.UserId = req.UserId
	res.OrderSign = GenOrderSign(req.UserId)
	if err := res.InsertOne(); err != nil {
		return nil, err
	}
	for i := range orderGoods {
		orderGoods[i].OrderId = res.ID
	}

	return OrderToOrderInfoRes(res), nil
}
