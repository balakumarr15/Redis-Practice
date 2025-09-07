package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Create a global context for Redis commands
var ctx = context.Background()

func main() {
	// 1. Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	// 2. Check connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis Connected:", pong)

	// 3. Set a key with a 10-second expiration
	fmt.Println("\n--- Setting key with 10-second expiry ---")
	err = rdb.Set(ctx, "session:token", "abc123", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}
	fmt.Println("Key 'session:token' set successfully!")

	// 4. Immediately retrieve the key
	val, err := rdb.Get(ctx, "session:token").Result()
	if err != nil {
		log.Fatalf("Error retrieving value: %v", err)
	}
	fmt.Printf("Immediately retrieved value: %s\n", val)

	// 5. Check remaining TTL
	ttl, err := rdb.TTL(ctx, "session:token").Result()
	if err != nil {
		log.Fatalf("Error fetching TTL: %v", err)
	}
	fmt.Printf("Initial TTL remaining: %v seconds\n", ttl.Seconds())

	// 6. Wait for 10 seconds
	fmt.Println("\nWaiting 10 seconds...")
	time.Sleep(10 * time.Second)

	// 7. Try retrieving the key after 10 seconds
	val, err = rdb.Get(ctx, "session:token").Result()
	if err == redis.Nil {
		fmt.Println("After 10 seconds: Key has expired and no longer exists.")
	} else if err != nil {
		log.Fatalf("Error retrieving after delay: %v", err)
	} else {
		fmt.Printf("After 10 seconds: Key still exists, value: %s\n", val)
	}

	// 8. Check TTL again after 10 seconds
	ttl, err = rdb.TTL(ctx, "session:token").Result()
	if err != nil {
		log.Fatalf("Error fetching TTL after delay: %v", err)
	}
	if ttl < 0 {
		fmt.Println("After 10 seconds: TTL expired or key does not exist.")
	} else {
		fmt.Printf("After 10 seconds: TTL remaining: %v seconds\n", ttl.Seconds())
	}

	// 9. Set another key without expiry
	fmt.Println("\n--- Setting key without expiry ---")
	err = rdb.Set(ctx, "user:1001", "Batman", 0).Err() // 0 = no expiry
	if err != nil {
		log.Fatalf("Error setting key without expiry: %v", err)
	}
	fmt.Println("Key 'user:1001' set without expiry!")

	// 10. Delete the key manually
	delCount, err := rdb.Del(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error deleting key: %v", err)
	}
	fmt.Printf("Deleted keys count: %d\n", delCount)
}
