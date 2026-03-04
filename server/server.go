// A server that handles coffee orders (requests from clients)
package main

import (
	"context"

	pb "github.com/stamakunov7/grpc-starbucks-coffee/coffeeshop"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(*pb.MenuRequest, pb.CoffeeShop_GetMenuServer) error {

}

func (s *server) MakeCoffee(context.Context, *pb.CoffeeRequest) (*pb.Coffee, error) {

}

func (s *server) PlaceOrder(context.Context, *pb.Order) (*pb.Receipt, error) {

}

func (s *server) GetOrderStatus(context.Context, *pb.Receipt) (*pb.OrderStatus, error) {

}
