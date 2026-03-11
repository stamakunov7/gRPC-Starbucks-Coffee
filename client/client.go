// An interface that clients use to make coffee orders (requests to the server)
package main

import (
	"context"
	"fmt"
	"io"
	"log" // timestamp and time (ex: 2026/03/08 12:00:00 in front of the message)
	"sort"
	"time"

	pb "grpc_starbuckscoffee/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	log.SetFlags(0) // remove timestamp from log output just to see clean output
}

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

	// Collect items from the server stream.
	itemsByID := map[string]*pb.Item{}
	for {
		resp, err := menuStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving menu items: %v", err)
		}
		for _, it := range resp.Items {
			// If server sends the full menu multiple times, dedupe by ID.
			itemsByID[it.Id] = it
		}
	}

	items := make([]*pb.Item, 0, len(itemsByID))
	for _, it := range itemsByID {
		items = append(items, it)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })

	if len(items) == 0 {
		log.Fatalf("no menu items received")
	}

	// Pretty-print menu.
	fmt.Println("Menu")
	fmt.Println("----")
	for _, it := range items {
		fmt.Printf("- %s: %s ($%.2f)\n  %s\n", it.Id, it.Name, it.Price, it.Description)
	}

	receipt, err := client.PlaceOrder(ctx, &pb.Order{
		Items: items,
	})
	if err != nil {
		log.Fatalf("error calling PlaceOrder: %v", err)
	}
	log.Printf("receipt: %s", receipt.Id)

	status, err := client.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("error calling GetOrderStatus: %v", err)
	}
	log.Printf("order status: %s (%s)", status.Status, status.OrderId)

}
