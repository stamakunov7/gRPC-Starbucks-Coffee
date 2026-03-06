// A server that handles coffee orders (requests from clients)
package main

import (
	"context"
	"log"
	"net"

	pb "grpc_starbuckscoffee/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{
		&pb.Item{
			Id:          "1",
			Name:        "Black Coffee",
			Description: "A cup of black colombian coffee",
			Price:       3.00,
		},
		&pb.Item{
			Id:          "2",
			Name:        "Vanilla Latte",
			Description: "A cup of vanilla latte",
			Price:       3.25,
		},
		&pb.Item{
			Id:          "3",
			Name:        "Matcha Latte",
			Description: "A cup of matcha latte",
			Price:       5.70,
		},
	}

	for i, _ := range items {
		if err := srv.Send(&pb.Menu{
			Items: items[0 : i+1],
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) MakeCoffee(context context.Context, coffeeRequest *pb.CoffeeRequest) (*pb.Coffee, error) {
	return &pb.Coffee{
		ItemName: coffeeRequest.ItemName,
		Size:     coffeeRequest.Size,
		Status:   "Ready",
	}, nil
}

func (s *server) PlaceOrder(context context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil
}

func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "In Progress",
	}, nil
}

func main() {
	// setup a listener on port 9001
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a new grpc server
	grpcServer := grpc.NewServer()
	// register the server
	pb.RegisterCoffeeShopServer(grpcServer, &server{})
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
