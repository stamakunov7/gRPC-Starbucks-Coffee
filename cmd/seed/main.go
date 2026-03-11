// Seed script populates Firestore (emulator or real) with categories and drinks.
// Run with: FIRESTORE_EMULATOR_HOST=localhost:8080 GOOGLE_CLOUD_PROJECT=grpc-starbucks-coffee go run ./cmd/seed
package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

type category struct {
	Active   bool   `firestore:"active"`
	Name     string `firestore:"name"`
	Slug     string `firestore:"slug"`
	SortOrder int   `firestore:"sortOrder"`
}

type drink struct {
	Name        string  `firestore:"name"`
	Description string  `firestore:"description"`
	Price       float64 `firestore:"price"`
	Category    string  `firestore:"category"`
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "grpc-starbucks-coffee"
		log.Printf("GOOGLE_CLOUD_PROJECT not set, using %q", projectID)
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	categories := []struct {
		ID string
		category
	}{
		{"protein_beverages", category{Active: true, Name: "Protein Beverages", Slug: "protein_beverages", SortOrder: 1}},
		{"hot_coffee", category{Active: true, Name: "Hot Coffee", Slug: "hot_coffee", SortOrder: 2}},
		{"cold_coffee", category{Active: true, Name: "Cold Coffee", Slug: "cold_coffee", SortOrder: 3}},
		{"matcha", category{Active: true, Name: "Matcha", Slug: "matcha", SortOrder: 4}},
		{"hot_tea", category{Active: true, Name: "Hot Tea", Slug: "hot_tea", SortOrder: 5}},
		{"cold_tea", category{Active: true, Name: "Cold Tea", Slug: "cold_tea", SortOrder: 6}},
		{"refreshers", category{Active: true, Name: "Refreshers", Slug: "refreshers", SortOrder: 7}},
		{"frappuccino_blended", category{Active: true, Name: "Frappuccino Blended Beverages", Slug: "frappuccino_blended", SortOrder: 8}},
		{"hot_chocolate_lemonade_more", category{Active: true, Name: "Hot Chocolate Lemonade & More", Slug: "hot_chocolate_lemonade_more", SortOrder: 9}},
		{"bottled_beverages", category{Active: true, Name: "Bottled Beverages", Slug: "bottled_beverages", SortOrder: 10}},
	}

	drinks := []struct {
		ID   string
		drink drink
	}{
		// Protein Beverages
		{"protein_cold_brew", drink{"Protein Cold Brew", "Cold brew with plant-based protein", 6.25, "protein_beverages"}},
		{"protein_chocolate_smoothie", drink{"Protein Chocolate Smoothie", "Chocolate smoothie with protein", 6.95, "protein_beverages"}},
		// Hot Coffee
		{"latte", drink{"Latte", "Espresso with steamed milk", 4.50, "hot_coffee"}},
		{"americano", drink{"Americano", "Espresso with hot water", 3.50, "hot_coffee"}},
		{"cappuccino", drink{"Cappuccino", "Espresso with foamed milk", 4.25, "hot_coffee"}},
		{"caramel_macchiato", drink{"Caramel Macchiato", "Espresso with vanilla and caramel", 5.25, "hot_coffee"}},
		{"flat_white", drink{"Flat White", "Ristretto shots with steamed milk", 4.75, "hot_coffee"}},
		// Cold Coffee
		{"cold_brew", drink{"Cold Brew", "Slow-steeped cold coffee", 4.00, "cold_coffee"}},
		{"iced_latte", drink{"Iced Latte", "Espresso with cold milk over ice", 4.50, "cold_coffee"}},
		{"iced_americano", drink{"Iced Americano", "Espresso with cold water over ice", 3.75, "cold_coffee"}},
		{"nitro_cold_brew", drink{"Nitro Cold Brew", "Cold brew infused with nitrogen", 4.95, "cold_coffee"}},
		// Matcha
		{"matcha_latte", drink{"Matcha Latte", "Green tea with steamed milk", 5.70, "matcha"}},
		{"matcha_espresso", drink{"Matcha Espresso Fusion", "Matcha and espresso", 5.95, "matcha"}},
		{"iced_matcha_latte", drink{"Iced Matcha Latte", "Matcha with cold milk over ice", 5.70, "matcha"}},
		// Hot Tea
		{"chai_latte", drink{"Chai Latte", "Spiced tea with steamed milk", 4.75, "hot_tea"}},
		{"green_tea", drink{"Green Tea", "Classic green tea", 3.00, "hot_tea"}},
		{"english_breakfast", drink{"English Breakfast Tea", "Black tea blend", 3.25, "hot_tea"}},
		{"honey_citrus_mint", drink{"Honey Citrus Mint Tea", "Green tea with honey and citrus", 4.50, "hot_tea"}},
		// Cold Tea
		{"iced_chai", drink{"Iced Chai Latte", "Chai with cold milk over ice", 4.75, "cold_tea"}},
		{"iced_green_tea", drink{"Iced Green Tea", "Green tea over ice", 3.25, "cold_tea"}},
		{"iced_black_tea", drink{"Iced Black Tea", "Black tea over ice", 3.00, "cold_tea"}},
		// Refreshers
		{"pink_drink", drink{"Pink Drink", "Strawberry acai with coconut milk", 5.45, "refreshers"}},
		{"mango_dragon", drink{"Mango Dragon Lemonade", "Mango and dragonfruit lemonade", 4.95, "refreshers"}},
		{"strawberry_acai", drink{"Strawberry Acai Lemonade", "Strawberry acai with lemonade", 4.95, "refreshers"}},
		{"paradise_drink", drink{"Paradise Drink", "Mango and passion fruit", 5.25, "refreshers"}},
		// Frappuccino Blended Beverages
		{"caramel_frappuccino", drink{"Caramel Frappuccino", "Coffee frappuccino with caramel", 5.50, "frappuccino_blended"}},
		{"vanilla_cream_frappuccino", drink{"Vanilla Cream Frappuccino", "Creamy vanilla frappuccino", 4.95, "frappuccino_blended"}},
		{"mocha_frappuccino", drink{"Mocha Frappuccino", "Coffee and chocolate frappuccino", 5.25, "frappuccino_blended"}},
		{"strawberry_creme_frappuccino", drink{"Strawberry Crème Frappuccino", "Strawberry and cream blended", 4.95, "frappuccino_blended"}},
		// Hot Chocolate Lemonade & More
		{"hot_chocolate", drink{"Hot Chocolate", "Steamed milk with chocolate", 4.25, "hot_chocolate_lemonade_more"}},
		{"steamed_milk", drink{"Steamed Milk", "Hot steamed milk", 3.00, "hot_chocolate_lemonade_more"}},
		{"steamed_apple_juice", drink{"Steamed Apple Juice", "Hot apple juice with cinnamon", 3.50, "hot_chocolate_lemonade_more"}},
		{"lemonade", drink{"Lemonade", "Fresh squeezed lemonade", 3.25, "hot_chocolate_lemonade_more"}},
		{"blended_strawberry_lemonade", drink{"Blended Strawberry Lemonade", "Strawberry and lemonade blended", 4.50, "hot_chocolate_lemonade_more"}},
		// Bottled Beverages
		{"evian", drink{"Evian Water", "Still bottled water", 2.50, "bottled_beverages"}},
		{"sparkling_water", drink{"Sparkling Water", "Carbonated water", 2.75, "bottled_beverages"}},
		{"cold_press_black", drink{"Starbucks Cold Press Black", "Cold-pressed black coffee", 3.95, "bottled_beverages"}},
		{"cold_press_vanilla", drink{"Starbucks Cold Press Vanilla", "Cold-pressed coffee with vanilla", 4.25, "bottled_beverages"}},
	}

	collCategories := client.Collection("categories")
	for _, c := range categories {
		_, err := collCategories.Doc(c.ID).Set(ctx, c.category)
		if err != nil {
			log.Fatalf("categories: set %q: %v", c.ID, err)
		}
		log.Printf("category: %s", c.ID)
	}

	collDrinks := client.Collection("drinks")
	for _, d := range drinks {
		_, err := collDrinks.Doc(d.ID).Set(ctx, d.drink)
		if err != nil {
			log.Fatalf("drinks: set %q: %v", d.ID, err)
		}
		log.Printf("drink: %s", d.ID)
	}

	log.Printf("seed done: %d categories, %d drinks", len(categories), len(drinks))
}
