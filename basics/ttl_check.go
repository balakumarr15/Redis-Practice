package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// TTL (Time To Live) operations in Redis
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

	// 1. TTL command - Check remaining time to live
	fmt.Println("\n=== TTL Command ===")

	// Set a key with expiration
	err = rdb.Set(ctx, "ttl_key", "ttl_value", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting ttl_key: %v", err)
	}
	fmt.Println("Set 'ttl_key' with 10-second expiration")

	// Check TTL immediately
	ttl, err := rdb.TTL(ctx, "ttl_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL: %v", err)
	}
	fmt.Printf("TTL remaining: %v (%.2f seconds)\n", ttl, ttl.Seconds())

	// 2. PTTL command - Check remaining time in milliseconds
	fmt.Println("\n=== PTTL Command ===")
	pttl, err := rdb.PTTL(ctx, "ttl_key").Result()
	if err != nil {
		log.Fatalf("Error getting PTTL: %v", err)
	}
	fmt.Printf("PTTL remaining: %v (%d milliseconds)\n", pttl, pttl.Milliseconds())

	// 3. EXPIRE command - Set expiration on existing key
	fmt.Println("\n=== EXPIRE Command ===")

	// Set a key without expiration
	err = rdb.Set(ctx, "no_expire_key", "no_expire_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting no_expire_key: %v", err)
	}
	fmt.Println("Set 'no_expire_key' without expiration")

	// Check TTL before setting expiration
	ttl, err = rdb.TTL(ctx, "no_expire_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL before expire: %v", err)
	}
	fmt.Printf("TTL before EXPIRE: %v\n", ttl)

	// Set expiration using EXPIRE
	expired, err := rdb.Expire(ctx, "no_expire_key", 5*time.Second).Result()
	if err != nil {
		log.Fatalf("Error setting expiration: %v", err)
	}
	fmt.Printf("EXPIRE command result: %v\n", expired)

	// Check TTL after setting expiration
	ttl, err = rdb.TTL(ctx, "no_expire_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL after expire: %v", err)
	}
	fmt.Printf("TTL after EXPIRE: %v (%.2f seconds)\n", ttl, ttl.Seconds())

	// 4. EXPIREAT command - Set expiration at specific timestamp
	fmt.Println("\n=== EXPIREAT Command ===")

	// Set a key for EXPIREAT demo
	err = rdb.Set(ctx, "expireat_key", "expireat_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting expireat_key: %v", err)
	}
	fmt.Println("Set 'expireat_key' for EXPIREAT demo")

	// Set expiration 3 seconds from now
	futureTime := time.Now().Add(3 * time.Second)
	expired, err = rdb.ExpireAt(ctx, "expireat_key", futureTime).Result()
	if err != nil {
		log.Fatalf("Error setting expireat: %v", err)
	}
	fmt.Printf("EXPIREAT command result: %v\n", expired)
	fmt.Printf("Key will expire at: %v\n", futureTime)

	// Check TTL
	ttl, err = rdb.TTL(ctx, "expireat_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL for expireat_key: %v", err)
	}
	fmt.Printf("TTL for expireat_key: %v (%.2f seconds)\n", ttl, ttl.Seconds())

	// 5. PERSIST command - Remove expiration
	fmt.Println("\n=== PERSIST Command ===")

	// Set a key with expiration
	err = rdb.Set(ctx, "persist_key", "persist_value", 15*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting persist_key: %v", err)
	}
	fmt.Println("Set 'persist_key' with 15-second expiration")

	// Check TTL before persist
	ttl, err = rdb.TTL(ctx, "persist_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL before persist: %v", err)
	}
	fmt.Printf("TTL before PERSIST: %v (%.2f seconds)\n", ttl, ttl.Seconds())

	// Remove expiration using PERSIST
	persisted, err := rdb.Persist(ctx, "persist_key").Result()
	if err != nil {
		log.Fatalf("Error persisting key: %v", err)
	}
	fmt.Printf("PERSIST command result: %v\n", persisted)

	// Check TTL after persist
	ttl, err = rdb.TTL(ctx, "persist_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL after persist: %v", err)
	}
	fmt.Printf("TTL after PERSIST: %v\n", ttl)

	// 6. Monitor expiration in real-time
	fmt.Println("\n=== Monitoring Expiration ===")

	// Set a key with short expiration for monitoring
	err = rdb.Set(ctx, "monitor_key", "monitor_value", 3*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting monitor_key: %v", err)
	}
	fmt.Println("Set 'monitor_key' with 3-second expiration for monitoring")

	// Monitor TTL every second
	for i := 0; i < 5; i++ {
		ttl, err = rdb.TTL(ctx, "monitor_key").Result()
		if err != nil {
			log.Fatalf("Error getting TTL during monitoring: %v", err)
		}

		if ttl < 0 {
			fmt.Printf("Time %d: Key has expired or doesn't exist\n", i+1)
			break
		} else {
			fmt.Printf("Time %d: TTL remaining: %.2f seconds\n", i+1, ttl.Seconds())
		}

		time.Sleep(1 * time.Second)
	}

	// 7. Check TTL for non-existent key
	fmt.Println("\n=== TTL for Non-existent Key ===")
	ttl, err = rdb.TTL(ctx, "nonexistent_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL for nonexistent key: %v", err)
	}
	fmt.Printf("TTL for nonexistent key: %v\n", ttl)

	// 8. Check TTL for key without expiration
	fmt.Println("\n=== TTL for Key Without Expiration ===")
	err = rdb.Set(ctx, "permanent_key", "permanent_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting permanent_key: %v", err)
	}

	ttl, err = rdb.TTL(ctx, "permanent_key").Result()
	if err != nil {
		log.Fatalf("Error getting TTL for permanent key: %v", err)
	}
	fmt.Printf("TTL for permanent key: %v\n", ttl)

	// 9. Summary of TTL return values
	fmt.Println("\n=== TTL Return Value Summary ===")
	fmt.Println("TTL return values:")
	fmt.Println("- Positive number: Key exists and has expiration set (seconds remaining)")
	fmt.Println("- -1: Key exists but has no expiration set")
	fmt.Println("- -2: Key does not exist")
	fmt.Println("- PTTL returns the same values but in milliseconds")
}
