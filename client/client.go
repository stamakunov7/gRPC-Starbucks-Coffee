// An interface that clients use to make coffee orders (requests to the server)
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
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

	client := pb.NewCoffeeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 1) Request menu
	menuStream, err := client.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("error calling GetMenu: %v", err)
	}

	itemsByID := make(map[string]*pb.Item)
	for {
		resp, err := menuStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving menu: %v", err)
		}
		for _, it := range resp.Items {
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

	// 2) Display menu (numbered)
	fmt.Println("\n--- Menu ---")
	for i, it := range items {
		fmt.Printf("  %2d. %s — %s ($%.2f)\n", i+1, it.Name, it.Description, it.Price)
	}
	fmt.Println()

	// 3) Ask user to select drinks (by number or ID)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter drink numbers or IDs (e.g. 1 3 5 or latte cold_brew): ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		log.Fatal("no items selected")
	}

	selected := selectItems(items, itemsByID, strings.Fields(line))
	if len(selected) == 0 {
		log.Fatal("no valid items selected")
	}

	// 4) Show cart and total
	var total float64
	fmt.Println("\n--- Your order ---")
	for _, it := range selected {
		fmt.Printf("  • %s — $%.2f\n", it.Name, it.Price)
		total += it.Price
	}
	fmt.Printf("  Total: $%.2f\n", total)
	fmt.Print("\nConfirm order? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	c := strings.TrimSpace(strings.ToLower(confirm))
	if c != "y" && c != "yes" {
		log.Fatal("order cancelled")
	}

	// 5) Place order
	receipt, err := client.PlaceOrder(ctx, &pb.Order{Items: selected})
	if err != nil {
		log.Fatalf("error placing order: %v", err)
	}
	fmt.Printf("\nOrder placed. Receipt ID: %s\n", receipt.Id)

	// 6) Poll order status until ready or max polls (server may never return Ready yet)
	fmt.Println("\nWaiting for your order...")
	const pollInterval = 2 * time.Second
	const maxPolls = 5
	for i := 0; i < maxPolls; i++ {
		status, err := client.GetOrderStatus(ctx, receipt)
		if err != nil {
			log.Printf("error getting status: %v", err)
			time.Sleep(pollInterval)
			continue
		}
		fmt.Printf("  Status: %s\n", status.Status)
		if isReady(status.Status) {
			fmt.Println("\nYour order is ready! Enjoy.")
			return
		}
		time.Sleep(pollInterval)
	}
	fmt.Printf("\nOrder received (receipt: %s). Server is still reporting In Progress; when the backend returns Ready, this client would stop here.\n", receipt.Id)
}

// selectItems parses user input (numbers 1-based or drink IDs) and returns matching items.
func selectItems(items []*pb.Item, byID map[string]*pb.Item, tokens []string) []*pb.Item {
	seen := make(map[string]bool)
	var out []*pb.Item
	for _, t := range tokens {
		t = strings.TrimSpace(strings.ToLower(t))
		if t == "" {
			continue
		}
		// Try as 1-based index
		var idx int
		if _, err := fmt.Sscanf(t, "%d", &idx); err == nil && idx >= 1 && idx <= len(items) {
			it := items[idx-1]
			if !seen[it.Id] {
				seen[it.Id] = true
				out = append(out, it)
			}
			continue
		}
		// Try as drink ID
		if it, ok := byID[t]; ok && !seen[it.Id] {
			seen[it.Id] = true
			out = append(out, it)
		}
	}
	return out
}

func isReady(status string) bool {
	s := strings.ToLower(strings.TrimSpace(status))
	return s == "ready" || s == "completed" || s == "ready for pickup" || s == "done"
}
