// A server that handles coffee orders (requests from clients)
package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "grpc_starbuckscoffee/proto"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
	firestoreClient *firestore.Client
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv pb.CoffeeShop_GetMenuServer) error {
	ctx := srv.Context()

	type drinkDoc struct {
		Name        string  `firestore:"name"`
		Description string  `firestore:"description"`
		Price       float64 `firestore:"price"`
		Category    string  `firestore:"category"`
	}

	iter := s.firestoreClient.Collection("drinks").Documents(ctx)
	defer iter.Stop()

	items := make([]*pb.Item, 0, 64)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		var d drinkDoc
		if err := doc.DataTo(&d); err != nil {
			return err
		}

		items = append(items, &pb.Item{
			Id:          doc.Ref.ID,
			Name:        d.Name,
			Description: d.Description,
			Price:       d.Price,
		})
	}

	if err := srv.Send(&pb.Menu{Items: items}); err != nil {
		return err
	}
	return nil
}

func (s *server) MakeCoffee(ctx context.Context, coffeeRequest *pb.CoffeeRequest) (*pb.Coffee, error) {
	return &pb.Coffee{
		ItemName: coffeeRequest.ItemName,
		Size:     coffeeRequest.Size,
		Status:   "Ready",
	}, nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
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
	// create a Firestore client
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "grpc-starbucks-coffee"
		log.Printf("GOOGLE_CLOUD_PROJECT not set, using %q", projectID)
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error creating Firestore Client: %v", err)
	}
	
	defer client.Close()
	// server implementation
	srvImpl := &server{firestoreClient: client}
	// register the server with Firestore
	pb.RegisterCoffeeShopServer(grpcServer, srvImpl)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
