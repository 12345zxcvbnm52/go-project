package initialize

import (
	"context"
	"fmt"
	gb "goods_srv/global"
	"goods_srv/model"
	"strconv"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type Logger struct{}

func (l Logger) Printf(format string, v ...interface{}) {
	zap.S().Infof(format, v...)
}

func InitEs() {
	dsn := fmt.Sprintf("http://%s:%d", gb.ServerConfig.EsConfig.Host, gb.ServerConfig.EsConfig.Port)
	var err error
	gb.EsConn, err = elastic.NewClient(elastic.SetURL(dsn), elastic.SetSniff(false), elastic.SetTraceLog(Logger{}))
	if err != nil {
		panic(err)
	}
	logic := model.EsGoods{}
	exist, err := gb.EsConn.IndexExists(logic.IndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exist {
		_, err := gb.EsConn.CreateIndex(logic.IndexName()).BodyString(logic.IndexMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	//DebugTest()
}

func DebugTest() {
	goods := []*model.Goods{}
	gb.DB.Find(&goods)
	logic := model.EsGoods{}
	for _, v := range goods {
		d := model.EsGoods{
			ID:          v.ID,
			BrandID:     v.BrandID,
			CategoryID:  v.CategoryID,
			IsHot:       v.IsHot,
			IsNew:       v.IsNew,
			GoodsSign:   v.GoodSign,
			GoodsBrief:  v.GoodsBrief,
			TransFree:   v.TransFree,
			OnSale:      v.OnSale,
			Name:        v.Name,
			ClickNum:    v.ClickNum,
			SoldNum:     v.SoldNum,
			FavorNum:    v.FavorNum,
			SalePrice:   v.SalePrice,
			MarketPrice: v.MarketPrice,
		}
		_, err := gb.EsConn.Index().Index(logic.IndexName()).BodyJson(d).Id(strconv.Itoa(int(d.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
