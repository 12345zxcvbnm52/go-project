package handler

import (
	pb "order_srv/proto"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}
