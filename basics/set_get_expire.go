package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Basic Redis operations: SET, GET, EXPIRE
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

	// 1. Basic SET and GET
	fmt.Println("\n=== Basic SET and GET ===")
	err = rdb.Set(ctx, "name", "Alice", 0).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}
	fmt.Println("Set 'name' to 'Alice'")

	val, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}
	fmt.Printf("Retrieved 'name': %s\n", val)

	// 2. SET with expiration
	fmt.Println("\n=== SET with Expiration ===")
	err = rdb.Set(ctx, "temp_key", "temporary_value", 5*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting key with expiry: %v", err)
	}
	fmt.Println("Set 'temp_key' with 5-second expiration")

	// Check TTL
	ttl, err := rdb.TTL(ctx, "temp_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL: %v", err)
	}
	fmt.Printf("TTL remaining: %v\n", ttl)

	// 3. EXPIRE command on existing key
	fmt.Println("\n=== EXPIRE Command ===")
	err = rdb.Set(ctx, "persistent_key", "persistent_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting persistent key: %v", err)
	}
	fmt.Println("Set 'persistent_key' without expiration")

	// Set expiration using EXPIRE
	expired, err := rdb.Expire(ctx, "persistent_key", 3*time.Second).Result()
	if err != nil {
		log.Fatalf("Error setting expiration: %v", err)
	}
	fmt.Printf("EXPIRE command result: %v\n", expired)

	// 4. EXPIREAT - Set expiration at specific time
	fmt.Println("\n=== EXPIREAT Command ===")
	futureTime := time.Now().Add(2 * time.Second)
	err = rdb.Set(ctx, "expireat_key", "expireat_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting expireat key: %v", err)
	}

	expired, err = rdb.ExpireAt(ctx, "expireat_key", futureTime).Result()
	if err != nil {
		log.Fatalf("Error setting expireat: %v", err)
	}
	fmt.Printf("EXPIREAT command result: %v\n", expired)

	// 5. PERSIST - Remove expiration
	fmt.Println("\n=== PERSIST Command ===")
	err = rdb.Set(ctx, "persist_key", "persist_value", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting persist key: %v", err)
	}

	persisted, err := rdb.Persist(ctx, "persist_key").Result()
	if err != nil {
		log.Fatalf("Error persisting key: %v", err)
	}
	fmt.Printf("PERSIST command result: %v\n", persisted)

	// 6. Wait and check expiration
	fmt.Println("\n=== Waiting for expiration ===")
	fmt.Println("Waiting 6 seconds for keys to expire...")
	time.Sleep(6 * time.Second)

	// Check if temp_key expired
	val, err = rdb.Get(ctx, "temp_key").Result()
	if err == redis.Nil {
		fmt.Println("'temp_key' has expired")
	} else if err != nil {
		log.Fatalf("Error checking temp_key: %v", err)
	} else {
		fmt.Printf("'temp_key' still exists: %s\n", val)
	}

	// Check if persistent_key expired
	val, err = rdb.Get(ctx, "persistent_key").Result()
	if err == redis.Nil {
		fmt.Println("'persistent_key' has expired")
	} else if err != nil {
		log.Fatalf("Error checking persistent_key: %v", err)
	} else {
		fmt.Printf("'persistent_key' still exists: %s\n", val)
	}

	// Check if expireat_key expired
	val, err = rdb.Get(ctx, "expireat_key").Result()
	if err == redis.Nil {
		fmt.Println("'expireat_key' has expired")
	} else if err != nil {
		log.Fatalf("Error checking expireat_key: %v", err)
	} else {
		fmt.Printf("'expireat_key' still exists: %s\n", val)
	}

	// Check if persist_key still exists (should not have expired)
	val, err = rdb.Get(ctx, "persist_key").Result()
	if err == redis.Nil {
		fmt.Println("'persist_key' has expired")
	} else if err != nil {
		log.Fatalf("Error checking persist_key: %v", err)
	} else {
		fmt.Printf("'persist_key' still exists: %s\n", val)
	}
}
