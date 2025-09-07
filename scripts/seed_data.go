package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// SeedData populates Redis with sample data for practice
func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis Connected:", pong)

	// Clear existing data
	fmt.Println("\n=== Clearing Existing Data ===")
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Fatalf("Error getting keys: %v", err)
	}

	if len(keys) > 0 {
		_, err = rdb.Del(ctx, keys...).Result()
		if err != nil {
			log.Fatalf("Error clearing data: %v", err)
		}
	}
	fmt.Printf("Cleared %d existing keys\n", len(keys))

	// 1. Seed user data
	fmt.Println("\n=== Seeding User Data ===")
	users := []struct {
		ID       string
		Name     string
		Email    string
		Age      int
		City     string
		Country  string
		Role     string
		Salary   int
		JoinDate time.Time
	}{
		{"user:1", "Alice Johnson", "alice@example.com", 28, "New York", "USA", "Developer", 75000, time.Now().AddDate(-2, -3, -15)},
		{"user:2", "Bob Smith", "bob@example.com", 32, "San Francisco", "USA", "Designer", 65000, time.Now().AddDate(-1, -8, -22)},
		{"user:3", "Charlie Brown", "charlie@example.com", 25, "London", "UK", "Developer", 55000, time.Now().AddDate(-1, -2, -10)},
		{"user:4", "Diana Prince", "diana@example.com", 30, "Paris", "France", "Manager", 85000, time.Now().AddDate(-3, -1, -5)},
		{"user:5", "Eve Wilson", "eve@example.com", 27, "Berlin", "Germany", "Developer", 60000, time.Now().AddDate(-1, -6, -18)},
		{"user:6", "Frank Miller", "frank@example.com", 35, "Tokyo", "Japan", "Architect", 90000, time.Now().AddDate(-4, -2, -8)},
		{"user:7", "Grace Lee", "grace@example.com", 29, "Seoul", "South Korea", "Designer", 70000, time.Now().AddDate(-2, -9, -12)},
		{"user:8", "Henry Davis", "henry@example.com", 33, "Sydney", "Australia", "Manager", 80000, time.Now().AddDate(-2, -11, -20)},
	}

	for _, user := range users {
		// Store user as hash
		err := rdb.HSet(ctx, user.ID, map[string]interface{}{
			"name":      user.Name,
			"email":     user.Email,
			"age":       user.Age,
			"city":      user.City,
			"country":   user.Country,
			"role":      user.Role,
			"salary":    user.Salary,
			"join_date": user.JoinDate.Format(time.RFC3339),
		}).Err()
		if err != nil {
			log.Fatalf("Error storing user %s: %v", user.ID, err)
		}
	}
	fmt.Printf("Seeded %d users\n", len(users))

	// 2. Seed product data
	fmt.Println("\n=== Seeding Product Data ===")
	products := []struct {
		ID          string
		Name        string
		Category    string
		Price       float64
		Stock       int
		Rating      float64
		Description string
	}{
		{"product:1", "MacBook Pro", "Electronics", 1999.99, 50, 4.8, "High-performance laptop for professionals"},
		{"product:2", "iPhone 15", "Electronics", 999.99, 100, 4.7, "Latest smartphone with advanced features"},
		{"product:3", "AirPods Pro", "Electronics", 249.99, 200, 4.6, "Wireless earbuds with noise cancellation"},
		{"product:4", "Nike Air Max", "Clothing", 129.99, 75, 4.5, "Comfortable running shoes"},
		{"product:5", "Coffee Maker", "Appliances", 89.99, 30, 4.3, "Automatic coffee brewing machine"},
		{"product:6", "Desk Chair", "Furniture", 199.99, 25, 4.4, "Ergonomic office chair"},
		{"product:7", "Book: Clean Code", "Books", 29.99, 100, 4.9, "Programming best practices guide"},
		{"product:8", "Yoga Mat", "Sports", 39.99, 150, 4.2, "Non-slip exercise mat"},
	}

	for _, product := range products {
		err := rdb.HSet(ctx, product.ID, map[string]interface{}{
			"name":        product.Name,
			"category":    product.Category,
			"price":       product.Price,
			"stock":       product.Stock,
			"rating":      product.Rating,
			"description": product.Description,
		}).Err()
		if err != nil {
			log.Fatalf("Error storing product %s: %v", product.ID, err)
		}
	}
	fmt.Printf("Seeded %d products\n", len(products))

	// 3. Seed leaderboard data
	fmt.Println("\n=== Seeding Leaderboard Data ===")
	leaderboard := "game_leaderboard"

	// Add players to leaderboard
	for i := 1; i <= 20; i++ {
		score := rand.Float64() * 10000
		playerID := fmt.Sprintf("player_%d", i)
		playerName := fmt.Sprintf("Player%d", i)

		err := rdb.ZAdd(ctx, leaderboard, redis.Z{
			Score:  score,
			Member: playerID,
		}).Err()
		if err != nil {
			log.Fatalf("Error adding player to leaderboard: %v", err)
		}

		// Store player metadata
		err = rdb.HSet(ctx, fmt.Sprintf("player:%s", playerID), map[string]interface{}{
			"name":  playerName,
			"level": rand.Intn(100) + 1,
			"class": []string{"Warrior", "Mage", "Rogue", "Paladin"}[rand.Intn(4)],
		}).Err()
		if err != nil {
			log.Fatalf("Error storing player metadata: %v", err)
		}
	}
	fmt.Printf("Seeded leaderboard with 20 players\n")

	// 4. Seed session data
	fmt.Println("\n=== Seeding Session Data ===")
	sessions := []struct {
		ID        string
		UserID    string
		Username  string
		Email     string
		CreatedAt time.Time
		ExpiresAt time.Time
	}{
		{"session:1", "user:1", "Alice Johnson", "alice@example.com", time.Now().Add(-2 * time.Hour), time.Now().Add(22 * time.Hour)},
		{"session:2", "user:2", "Bob Smith", "bob@example.com", time.Now().Add(-1 * time.Hour), time.Now().Add(23 * time.Hour)},
		{"session:3", "user:3", "Charlie Brown", "charlie@example.com", time.Now().Add(-30 * time.Minute), time.Now().Add(23*time.Hour + 30*time.Minute)},
	}

	for _, session := range sessions {
		sessionData := fmt.Sprintf(`{"user_id":"%s","username":"%s","email":"%s","created_at":%d,"expires_at":%d,"data":{}}`,
			session.UserID, session.Username, session.Email, session.CreatedAt.Unix(), session.ExpiresAt.Unix())

		ttl := session.ExpiresAt.Sub(time.Now())
		err := rdb.Set(ctx, session.ID, sessionData, ttl).Err()
		if err != nil {
			log.Fatalf("Error storing session: %v", err)
		}
	}
	fmt.Printf("Seeded %d sessions\n", len(sessions))

	// 5. Seed chat data
	fmt.Println("\n=== Seeding Chat Data ===")
	rooms := []string{"general", "gaming", "tech", "random"}

	for _, room := range rooms {
		// Add users to room
		userKey := fmt.Sprintf("chat_users:%s", room)
		for i := 1; i <= 5; i++ {
			userID := fmt.Sprintf("user:%d", i)
			err := rdb.SAdd(ctx, userKey, userID).Err()
			if err != nil {
				log.Fatalf("Error adding user to room: %v", err)
			}
		}

		// Add some messages to room
		historyKey := fmt.Sprintf("chat_history:%s", room)
		messages := []string{
			fmt.Sprintf("Welcome to %s room!", room),
			"Hello everyone!",
			"How is everyone doing?",
			"Great to be here!",
			"Looking forward to the discussion",
		}

		for i, message := range messages {
			userID := fmt.Sprintf("user:%d", (i%5)+1)
			messageData := fmt.Sprintf("%d|%s|%s|%s|%d",
				time.Now().UnixNano()+int64(i),
				room,
				userID,
				message,
				time.Now().Unix()+int64(i))

			err := rdb.LPush(ctx, historyKey, messageData).Err()
			if err != nil {
				log.Fatalf("Error adding message to room: %v", err)
			}
		}
	}
	fmt.Printf("Seeded %d chat rooms\n", len(rooms))

	// 6. Seed event stream data
	fmt.Println("\n=== Seeding Event Stream Data ===")
	stream := "events"

	eventTypes := []string{"user_login", "user_logout", "page_view", "purchase", "search", "click"}

	for i := 0; i < 50; i++ {
		eventType := eventTypes[rand.Intn(len(eventTypes))]
		userID := fmt.Sprintf("user:%d", rand.Intn(8)+1)

		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{
				"event_type": eventType,
				"user_id":    userID,
				"timestamp":  time.Now().Unix(),
				"ip_address": fmt.Sprintf("192.168.1.%d", rand.Intn(255)+1),
				"user_agent": "Mozilla/5.0 (compatible; RedisPractice/1.0)",
			},
		}).Result()
		if err != nil {
			log.Fatalf("Error adding event to stream: %v", err)
		}

		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	}
	fmt.Printf("Seeded event stream with 50 events\n")

	// 7. Seed rate limiting data
	fmt.Println("\n=== Seeding Rate Limiting Data ===")
	rateLimitKeys := []string{"api:user:1", "api:user:2", "api:user:3", "api:global"}

	for _, key := range rateLimitKeys {
		// Add some requests to rate limiter
		for i := 0; i < 3; i++ {
			_, err := rdb.ZAdd(ctx, fmt.Sprintf("rate_limit:%s", key), redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: fmt.Sprintf("%d", time.Now().UnixNano()),
			}).Err()
			if err != nil {
				log.Fatalf("Error adding to rate limiter: %v", err)
			}
		}
	}
	fmt.Printf("Seeded rate limiting data for %d keys\n", len(rateLimitKeys))

	// 8. Seed tags and categories
	fmt.Println("\n=== Seeding Tags and Categories ===")

	// Product tags
	productTags := []string{"electronics", "clothing", "books", "sports", "furniture", "appliances", "gadgets", "accessories"}
	for _, tag := range productTags {
		err := rdb.SAdd(ctx, "product_tags", tag).Err()
		if err != nil {
			log.Fatalf("Error adding product tag: %v", err)
		}
	}

	// User interests
	interests := []string{"programming", "gaming", "music", "sports", "travel", "cooking", "photography", "reading"}
	for i := 1; i <= 8; i++ {
		userInterests := []string{}
		numInterests := rand.Intn(4) + 2 // 2-5 interests per user
		for j := 0; j < numInterests; j++ {
			interest := interests[rand.Intn(len(interests))]
			userInterests = append(userInterests, interest)
		}

		for _, interest := range userInterests {
			err := rdb.SAdd(ctx, fmt.Sprintf("user_interests:user:%d", i), interest).Err()
			if err != nil {
				log.Fatalf("Error adding user interest: %v", err)
			}
		}
	}
	fmt.Printf("Seeded tags and categories\n")

	// 9. Final statistics
	fmt.Println("\n=== Seeding Complete ===")

	// Get total key count
	totalKeys, err := rdb.DBSize(ctx).Result()
	if err != nil {
		log.Fatalf("Error getting database size: %v", err)
	}

	fmt.Printf("Total keys in database: %d\n", totalKeys)

	// Get memory usage
	info, err := rdb.Info(ctx, "memory").Result()
	if err != nil {
		log.Fatalf("Error getting memory info: %v", err)
	}

	fmt.Printf("Memory usage info available in Redis INFO\n")
	fmt.Println("\nSample data has been successfully seeded!")
	fmt.Println("You can now run the practice examples with realistic data.")
}
