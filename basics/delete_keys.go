package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis DELETE and EXISTS operations
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

	// 1. Set up some test keys
	fmt.Println("\n=== Setting up test keys ===")
	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	values := []string{"value1", "value2", "value3", "value4", "value5"}

	for i, key := range keys {
		err := rdb.Set(ctx, key, values[i], 0).Err()
		if err != nil {
			log.Fatalf("Error setting key %s: %v", key, err)
		}
		fmt.Printf("Set %s = %s\n", key, values[i])
	}

	// 2. Check if keys exist using EXISTS
	fmt.Println("\n=== EXISTS Command ===")
	for _, key := range keys {
		exists, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			log.Fatalf("Error checking existence of %s: %v", key, err)
		}
		fmt.Printf("Key '%s' exists: %v\n", key, exists > 0)
	}

	// Check for non-existent key
	exists, err := rdb.Exists(ctx, "nonexistent").Result()
	if err != nil {
		log.Fatalf("Error checking nonexistent key: %v", err)
	}
	fmt.Printf("Key 'nonexistent' exists: %v\n", exists > 0)

	// 3. Delete single key using DEL
	fmt.Println("\n=== DEL Command (single key) ===")
	deleted, err := rdb.Del(ctx, "key1").Result()
	if err != nil {
		log.Fatalf("Error deleting key1: %v", err)
	}
	fmt.Printf("Deleted %d key(s)\n", deleted)

	// Verify deletion
	exists, err = rdb.Exists(ctx, "key1").Result()
	if err != nil {
		log.Fatalf("Error checking key1 after deletion: %v", err)
	}
	fmt.Printf("Key 'key1' exists after deletion: %v\n", exists > 0)

	// 4. Delete multiple keys at once
	fmt.Println("\n=== DEL Command (multiple keys) ===")
	keysToDelete := []string{"key2", "key3", "nonexistent_key"}
	deleted, err = rdb.Del(ctx, keysToDelete...).Result()
	if err != nil {
		log.Fatalf("Error deleting multiple keys: %v", err)
	}
	fmt.Printf("Deleted %d key(s) out of %d requested\n", deleted, len(keysToDelete))

	// Verify deletions
	for _, key := range keysToDelete {
		exists, err = rdb.Exists(ctx, key).Result()
		if err != nil {
			log.Fatalf("Error checking %s after deletion: %v", key, err)
		}
		fmt.Printf("Key '%s' exists after deletion: %v\n", key, exists > 0)
	}

	// 5. UNLINK command (asynchronous deletion)
	fmt.Println("\n=== UNLINK Command ===")
	// Set a key for UNLINK demo
	err = rdb.Set(ctx, "unlink_key", "unlink_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting unlink_key: %v", err)
	}
	fmt.Println("Set 'unlink_key' for UNLINK demo")

	// Use UNLINK instead of DEL
	unlinked, err := rdb.Unlink(ctx, "unlink_key").Result()
	if err != nil {
		log.Fatalf("Error unlinking key: %v", err)
	}
	fmt.Printf("Unlinked %d key(s)\n", unlinked)

	// Verify unlink
	exists, err = rdb.Exists(ctx, "unlink_key").Result()
	if err != nil {
		log.Fatalf("Error checking unlink_key after unlink: %v", err)
	}
	fmt.Printf("Key 'unlink_key' exists after unlink: %v\n", exists > 0)

	// 6. RENAME command
	fmt.Println("\n=== RENAME Command ===")
	// Set a key to rename
	err = rdb.Set(ctx, "old_name", "old_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting old_name: %v", err)
	}
	fmt.Println("Set 'old_name' for rename demo")

	// Rename the key
	err = rdb.Rename(ctx, "old_name", "new_name").Err()
	if err != nil {
		log.Fatalf("Error renaming key: %v", err)
	}
	fmt.Println("Renamed 'old_name' to 'new_name'")

	// Verify rename
	val, err := rdb.Get(ctx, "new_name").Result()
	if err != nil {
		log.Fatalf("Error getting renamed key: %v", err)
	}
	fmt.Printf("New key value: %s\n", val)

	// Check if old key exists
	exists, err = rdb.Exists(ctx, "old_name").Result()
	if err != nil {
		log.Fatalf("Error checking old_name: %v", err)
	}
	fmt.Printf("Old key 'old_name' exists: %v\n", exists > 0)

	// 7. RENAMENX command (rename only if new name doesn't exist)
	fmt.Println("\n=== RENAMENX Command ===")
	// Set another key
	err = rdb.Set(ctx, "source", "source_value", 0).Err()
	if err != nil {
		log.Fatalf("Error setting source: %v", err)
	}

	// Try to rename to existing key (should fail)
	renamed, err := rdb.RenameNX(ctx, "source", "new_name").Result()
	if err != nil {
		log.Fatalf("Error with RENAMENX: %v", err)
	}
	fmt.Printf("RENAMENX to existing key result: %v\n", renamed)

	// Try to rename to non-existing key (should succeed)
	renamed, err = rdb.RenameNX(ctx, "source", "unique_name").Result()
	if err != nil {
		log.Fatalf("Error with RENAMENX: %v", err)
	}
	fmt.Printf("RENAMENX to unique key result: %v\n", renamed)

	// Verify the rename
	val, err = rdb.Get(ctx, "unique_name").Result()
	if err != nil {
		log.Fatalf("Error getting unique_name: %v", err)
	}
	fmt.Printf("Renamed key value: %s\n", val)

	// 8. Show remaining keys
	fmt.Println("\n=== Remaining keys ===")
	remainingKeys := []string{"key4", "key5", "new_name", "unique_name"}
	for _, key := range remainingKeys {
		exists, err = rdb.Exists(ctx, key).Result()
		if err != nil {
			log.Fatalf("Error checking %s: %v", key, err)
		}
		if exists > 0 {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				log.Fatalf("Error getting %s: %v", key, err)
			}
			fmt.Printf("Key '%s' = %s\n", key, val)
		}
	}
}
