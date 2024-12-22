package handler

import (
	"context"
	gb "inventory_srv/global"
	"inventory_srv/model"
	pb "inventory_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServer
}

func (is *InventoryServer) SetStock(ctx context.Context, req *pb.WriteInvtReq) (*emptypb.Empty, error) {
	u := &model.Inventory{
		GoodsId:  req.GoodsId,
		GoodsNum: req.GoodsNum,
	}
	if err := u.InsertOne(); err != nil {
		if err == model.ErrDuplicated {
			err = u.UpdateOneByGoodsId()
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (is *InventoryServer) FindStock(ctx context.Context, req *pb.StockInfoReq) (*pb.StockInfoRes, error) {
	u := &model.Inventory{
		GoodsId: req.GoodsId,
	}
	if err := u.FindOneByGoodsId(); err != nil {
		return nil, err
	}
	return &pb.StockInfoRes{
		GoodsId:  u.GoodsId,
		GoodsNum: u.GoodsNum,
	}, nil
}

func (is *InventoryServer) DecrStock(ctx context.Context, req *pb.DecrStockReq) (*emptypb.Empty, error) {
	tx := gb.DB.Begin()
	ids := make([]int32, len(req.DecrData))
	for i, v := range req.DecrData {
		ids[i] = v.GoodsId
	}
	logic := &model.Inventory{}
	res, err := logic.FindByGoodsIds(ids...)
	if err != nil || int(res.Total) != len(req.DecrData) {
		tx.Rollback()
		return nil, err
	}

	for i, v := range res.Data {
		if v.GoodsNum < req.DecrData[i].GoodsNum {
			zap.S().Errorw("一次事务执行失败", "msg", "库存不足")
			tx.Rollback()
			return nil, model.ErrInvalidArgument
		}
		v.GoodsNum -= req.DecrData[i].GoodsNum
		r := tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).Update("goods_num", v.GoodsNum)
		if r.Error != nil {
			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
			tx.Rollback()
			return nil, model.ErrInternalWrong
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (is *InventoryServer) IncrStock(ctx context.Context, req *pb.IncrStockReq) (*emptypb.Empty, error) {
	tx := gb.DB.Begin()
	ids := make([]int32, len(req.IncrData))
	for i, v := range req.IncrData {
		ids[i] = v.GoodsId
	}
	logic := &model.Inventory{}
	res, err := logic.FindByGoodsIds(ids...)
	if err != nil || int(res.Total) != len(req.IncrData) {
		tx.Rollback()
		return nil, err
	}

	for i, v := range res.Data {
		v.GoodsNum += req.IncrData[i].GoodsNum
		r := tx.Model(&model.Inventory{}).Where("goods_id = ?", v.GoodsId).Update("goods_num", v.GoodsNum)
		if r.Error != nil {
			zap.S().Errorw("一次事务执行失败", "msg", r.Error.Error())
			tx.Rollback()
			return nil, model.ErrInternalWrong
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
