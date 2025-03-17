package main

import (
	"context"
	"fmt"
	"kenshop/goken/registry/registor"
	"kenshop/goken/registry/ways/consul"
	"kenshop/goken/server/rpcserver"
	proto "kenshop/test_srv/proto"
	"log"
	"net"
	"os"
	"time"

	dtmgrpc "github.com/dtm-labs/client/dtmgrpc"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Count struct {
	gorm.Model
	Money int32 `gorm:""`
}

type ConnServer struct {
	proto.UnimplementedDmtServer
	DB *gorm.DB
}

func (s *ConnServer) ATrans(ctx context.Context, in *proto.A) (*emptypb.Empty, error) {
	tx := s.DB.Begin()
	u := &Count{}
	if res := tx.Model(&Count{}).First(u, in.Id); res.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}

	if u.Money < in.Decr {
		tx.Rollback()
		return nil, status.Error(codes.FailedPrecondition, "A余额不足")
	}

	u.Money -= in.Decr
	if res := tx.Model(&Count{}).Where("id = ?", in.Id).Update("money", u.Money); res.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *ConnServer) ACom(ctx context.Context, in *proto.A) (*emptypb.Empty, error) {
	if res := s.DB.Model(&Count{}).Where("id = ?", in.Id).Update("money", gorm.Expr("money + ? ", in.Decr)); res.Error != nil {
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *ConnServer) BTrans(ctx context.Context, in *proto.B) (*emptypb.Empty, error) {
	tx := s.DB.Begin()
	u := &Count{}
	if res := tx.Model(&Count{}).First(u, in.Id); res.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}

	u.Money += in.Incr
	if res := tx.Model(&Count{}).Where("id = ?", in.Id).Update("money", u.Money); res.Error != nil {
		tx.Rollback()
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (s *ConnServer) BCom(ctx context.Context, in *proto.B) (*emptypb.Empty, error) {
	if res := s.DB.Model(&Count{}).Where("id = ?", in.Id).Update("money", gorm.Expr("money - ? ", in.Incr)); res.Error != nil {
		return nil, status.Error(codes.Aborted, res.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *ConnServer) Try(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	At := "consul://192.168.199.128:8500/ken/dmt/ATrans"
	Ac := "consul://192.168.199.128:8500/ken/dmt/ACom"
	Bt := "consul://192.168.199.128:8500/ken/dmt/BTrans"
	Bc := "consul://192.168.199.128:8500/ken/dmt/BCom"
	A := &proto.A{Id: 1, Decr: 50}
	B := &proto.B{Id: 2, Incr: 50}
	gid, _ := uuid.NewV7()
	saga := dtmgrpc.NewSagaGrpc("192.168.199.128:36790", gid.String()).
		Add(At, Ac, A).
		Add(Bt, Bc, B)
	saga.WaitResult = true
	if err := saga.Submit(); err != nil {
		fmt.Println(err.Error())
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func main() {
	lis, _ := net.Listen("tcp", "192.168.199.128:33333")

	cfg := api.DefaultConfig()
	cfg.Address = "192.168.199.128:8500"
	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	r := registor.MustNewRegister(
		consul.MustNewConsulRegistor(cli,
			consul.WithEnableHealthCheck(true),
			consul.WithDeregisterCriticalServiceAfter("30s"),
			consul.WithHealthcheckInterval("10s"),
		),
	)

	s := rpcserver.MustNewServer(context.Background(),
		rpcserver.WithHost("127.0.0.1:33333"),
		rpcserver.WithServiceName("ken"),
		rpcserver.WithListener(lis),
		rpcserver.WithRegistor(r),
		//rpcserver.WithUnaryInts(sinterceptors.UnaryTracingInterceptor),
	)

	dsn := "root:123@tcp(192.168.199.128:3307)/ken?charset=utf8mb4&parseTime=True&loc=Local"

	log := logger.New(
		log.New(os.Stdout, "", log.LstdFlags|log.Llongfile),
		logger.Config{
			SlowThreshold: 1 * time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		//让所有gorm适配的数据库的Err同步为gorm.Err类型
		TranslateError: true,
		Logger:         log,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Count{})
	cs := &ConnServer{DB: db}

	proto.RegisterDmtServer(s.Server, cs)
	if err := s.Serve(); err != nil {
		panic(err)
	}
}
