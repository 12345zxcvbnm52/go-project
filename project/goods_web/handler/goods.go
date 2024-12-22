package handler

import (
	"context"
	pb "goods_web/proto"
	"goods_web/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var RpcPool util.Pooler = &util.Pool{}

func GetGoodsList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := RpcPool.Value()
	if err != nil {
		panic(err)
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
	if cid := c.DefaultQuery("cl", "0"); cid != "0" {
		cid, _ := strconv.Atoi(cid)
		fliter.CategyId = int32(cid)
	}
	if bId := c.DefaultQuery("bid", "0"); bId != "0" {
		bid, _ := strconv.Atoi(bId)
		fliter.Brand = int32(bid)
	}
	fliter.KeyWords = c.DefaultQuery("kw", "")
	res, err := cc.GetGoodList(ctx, &fliter)
	if err != nil {
		zap.S().Errorw("微服务调用失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"错误信息": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": res.Total,
		"data":  res.Data,
	})
}
