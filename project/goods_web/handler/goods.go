package handler

import (
	"context"
	"goods_web/form"
	gb "goods_web/global"
	pb "goods_web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 根据条件查找对应的商品
func GetGoodsList(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "ginCtx", c)
	client, err := gb.RpcPool.Value()
	if err != nil {
		zap.S().Errorw("池内获取的连接不可用", "msg", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		return
	}
	defer client.Close()
	cc := pb.NewGoodsClient(client)

	fliter := pb.GoodsFilterReq{}
	if minPrice := c.DefaultQuery("minp", "0"); minPrice != "0" {
		//哪怕转化失败,mp为0也不影响程序,故不检错
		mp, _ := strconv.Atoi(minPrice)
		fliter.MinPrice = int32(mp)
	}
	if maxPrice := c.DefaultQuery("maxp", "0"); maxPrice != "0" {
		mp, _ := strconv.Atoi(maxPrice)
		fliter.MaxPrice = int32(mp)
	}
	if isHot := c.DefaultQuery("ih", "0"); isHot == "1" {
		fliter.IsHot = true
	}
	if isNew := c.DefaultQuery("in", "0"); isNew == "1" {
		fliter.IsNew = true
	}
	if onTab := c.DefaultQuery("ot", "0"); onTab == "1" {
		fliter.OnTable = true
	}
	if pageSize := c.DefaultQuery("ps", "0"); pageSize != "0" {
		ps, _ := strconv.Atoi(pageSize)
		fliter.PageSize = int32(ps)
	}
	PagesNum := c.DefaultQuery("pn", "0")
	pn, _ := strconv.Atoi(PagesNum)
	fliter.PagesNum = int32(pn)
	if cid := c.DefaultQuery("cid", "0"); cid != "0" {
		cid, _ := strconv.Atoi(cid)
		fliter.CategyId = uint32(cid)
	}
	if bId := c.DefaultQuery("bid", "0"); bId != "0" {
		bid, _ := strconv.Atoi(bId)
		fliter.BrandId = uint32(bid)
	}
	fliter.KeyWords = c.DefaultQuery("kw", "")
	res, err := cc.GetGoodList(ctx, &fliter)
	if err != nil {
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": res.Total,
		"data":  res.Data,
	})
}

func CreateGoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := gb.RpcPool.Value()
	if err != nil {
		zap.S().Errorw("池内获取的连接不可用", "msg", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		return
	}
	defer client.Close()
	cc := pb.NewGoodsClient(client)

	u := &form.GoodsWriteForm{}
	if err := c.BindJSON(u); err != nil {
		ValidatorErrorHandle(c, err)
		return
	}

	req := &pb.WriteGoodsInfoReq{
		Name:        u.Name,
		GoodsSign:   u.GoodsSign,
		GoodsBrief:  u.GoodsBrief,
		MarketPrice: u.MarketPrice,
		SalePrice:   u.SalePrice,
		Images:      u.Images,
		FirstImage:  u.FirstImage,
		DescImages:  u.DescImages,
		TransFree:   u.TransFree,
		CategyId:    uint32(u.CategoryID),
		BrandId:     uint32(u.BrandID),
	}
	res, err := cc.CreateGoods(ctx, req)
	if err != nil {
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"商品id": res.Id,
	})
}
