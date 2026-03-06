// An interface that clients use to make coffee orders (requests to the server)
package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "grpc_starbuckscoffee/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// setup a connection to the server
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	// create a new client
	client := pb.NewCoffeeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// make a request to the server
	menuStream, err := client.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("error calling GetMenu function: %v", err)
	}

	done := make(chan bool)

	var items []*pb.Item
	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("error receiving menu items: %v", err)
			}
			items = resp.Items
			log.Printf("received menu items: %v", resp.Items)
		}
	}()

	<-done
	if len(items) == 0 {
		log.Fatalf("no menu items received")
	}
	receipt, err := client.PlaceOrder(ctx, &pb.Order{
		Items: items,
	})
	log.Printf("receipt: %v", receipt)

	status, err := client.GetOrderStatus(ctx, receipt)
	log.Printf("order status: %v", status)

}
