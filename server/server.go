// A server that handles coffee orders (requests from clients)
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "grpc_starbuckscoffee/proto"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if len(order.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	// Generate next order ID (001, 002, ...) in a transaction
	var orderID string
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ref := s.firestoreClient.Collection("counters").Doc("orders")
		doc, err := tx.Get(ref)
		var next int
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return err
			}
			next = 1
		} else {
			v, err := doc.DataAt("last")
			if err != nil {
				next = 1
			} else if n, ok := v.(int64); ok {
				next = int(n) + 1
			} else if n, ok := v.(float64); ok {
				next = int(n) + 1
			} else {
				next = 1
			}
		}
		orderID = fmt.Sprintf("%03d", next)
		if err := tx.Set(ref, map[string]interface{}{"last": next}); err != nil {
			return err
		}

		// Build order document (receipt data)
		var total float64
		itemsData := make([]map[string]interface{}, 0, len(order.Items))
		for _, it := range order.Items {
			itemsData = append(itemsData, map[string]interface{}{
				"id": it.Id, "name": it.Name, "price": it.Price,
			})
			total += it.Price
		}

		orderRef := s.firestoreClient.Collection("orders").Doc(orderID)
		return tx.Set(orderRef, map[string]interface{}{
			"items":     itemsData,
			"total":     total,
			"status":    "received",
			"createdAt": firestore.ServerTimestamp,
		})
	})
	if err != nil {
		return nil, err
	}

	// Simulate barista: update status to preparing -> ready after delays
	go s.simulateOrderProgress(context.Background(), orderID)

	log.Printf("Order %s created (status will advance to preparing → ready)", orderID)
	return &pb.Receipt{Id: orderID}, nil
}

func (s *server) simulateOrderProgress(ctx context.Context, orderID string) {
	ref := s.firestoreClient.Collection("orders").Doc(orderID)
	time.Sleep(3 * time.Second)
	_, _ = ref.Update(ctx, []firestore.Update{{Path: "status", Value: "preparing"}})
	time.Sleep(3 * time.Second)
	_, _ = ref.Update(ctx, []firestore.Update{{Path: "status", Value: "ready"}})
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	if receipt.Id == "" {
		return nil, fmt.Errorf("receipt id is required")
	}
	doc, err := s.firestoreClient.Collection("orders").Doc(receipt.Id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	data := doc.Data()
	status, _ := data["status"].(string)
	if status == "" {
		status = "unknown"
	}
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  status,
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
	emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if emulatorHost == "" {
		log.Printf("WARNING: FIRESTORE_EMULATOR_HOST not set — using real Firestore (or connection may fail)")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error creating Firestore Client: %v", err)
	}
	defer client.Close()

	srvImpl := &server{firestoreClient: client}
	pb.RegisterCoffeeShopServer(grpcServer, srvImpl)
	log.Printf("CoffeeShop server listening on :9001 (Firestore project=%s, emulator=%s)", projectID, emulatorHost)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
