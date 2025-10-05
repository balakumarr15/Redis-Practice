package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis Hash operations - HSET, HGET, HMSET, etc.
func HashOperations() {
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

	// 1. HSET - Set hash fields
	fmt.Println("\n=== HSET Command ===")

	// Set single field
	added, err := rdb.HSet(ctx, "user:1001", "name", "Alice").Result()
	if err != nil {
		log.Fatalf("Error setting hash field: %v", err)
	}
	fmt.Printf("HSET single field result: %d\n", added)

	// Set multiple fields at once
	fields := map[string]interface{}{
		"age":     30,
		"email":   "alice@example.com",
		"city":    "New York",
		"country": "USA",
	}
	added, err = rdb.HSet(ctx, "user:1001", fields).Result()
	if err != nil {
		log.Fatalf("Error setting multiple hash fields: %v", err)
	}
	fmt.Printf("HSET multiple fields result: %d\n", added)

	// 2. HGET - Get hash field value
	fmt.Println("\n=== HGET Command ===")

	name, err := rdb.HGet(ctx, "user:1001", "name").Result()
	if err != nil {
		log.Fatalf("Error getting hash field: %v", err)
	}
	fmt.Printf("User name: %s\n", name)

	age, err := rdb.HGet(ctx, "user:1001", "age").Result()
	if err != nil {
		log.Fatalf("Error getting age: %v", err)
	}
	fmt.Printf("User age: %s\n", age)

	// 3. HGETALL - Get all hash fields and values
	fmt.Println("\n=== HGETALL Command ===")

	allFields, err := rdb.HGetAll(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error getting all hash fields: %v", err)
	}
	fmt.Println("All user fields:")
	for field, value := range allFields {
		fmt.Printf("  %s: %s\n", field, value)
	}

	// 4. HMGET - Get multiple hash fields
	fmt.Println("\n=== HMGET Command ===")

	values, err := rdb.HMGet(ctx, "user:1001", "name", "email", "city").Result()
	if err != nil {
		log.Fatalf("Error getting multiple hash fields: %v", err)
	}
	fmt.Printf("Name: %v, Email: %v, City: %v\n", values[0], values[1], values[2])

	// 5. HKEYS - Get all hash field names
	fmt.Println("\n=== HKEYS Command ===")

	keys, err := rdb.HKeys(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error getting hash keys: %v", err)
	}
	fmt.Printf("Hash field names: %v\n", keys)

	// 6. HVALS - Get all hash field values
	fmt.Println("\n=== HVALS Command ===")

	vals, err := rdb.HVals(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error getting hash values: %v", err)
	}
	fmt.Printf("Hash field values: %v\n", vals)

	// 7. HEXISTS - Check if hash field exists
	fmt.Println("\n=== HEXISTS Command ===")

	exists, err := rdb.HExists(ctx, "user:1001", "name").Result()
	if err != nil {
		log.Fatalf("Error checking field existence: %v", err)
	}
	fmt.Printf("Field 'name' exists: %v\n", exists)

	exists, err = rdb.HExists(ctx, "user:1001", "phone").Result()
	if err != nil {
		log.Fatalf("Error checking field existence: %v", err)
	}
	fmt.Printf("Field 'phone' exists: %v\n", exists)

	// 8. HDEL - Delete hash fields
	fmt.Println("\n=== HDEL Command ===")

	deleted, err := rdb.HDel(ctx, "user:1001", "country").Result()
	if err != nil {
		log.Fatalf("Error deleting hash field: %v", err)
	}
	fmt.Printf("Deleted %d field(s)\n", deleted)

	// Verify deletion
	exists, err = rdb.HExists(ctx, "user:1001", "country").Result()
	if err != nil {
		log.Fatalf("Error checking field existence after deletion: %v", err)
	}
	fmt.Printf("Field 'country' exists after deletion: %v\n", exists)

	// 9. HINCRBY - Increment hash field by integer
	fmt.Println("\n=== HINCRBY Command ===")

	// Increment age by 1
	newAge, err := rdb.HIncrBy(ctx, "user:1001", "age", 1).Result()
	if err != nil {
		log.Fatalf("Error incrementing age: %v", err)
	}
	fmt.Printf("New age after increment: %d\n", newAge)

	// Increment by 5
	newAge, err = rdb.HIncrBy(ctx, "user:1001", "age", 5).Result()
	if err != nil {
		log.Fatalf("Error incrementing age by 5: %v", err)
	}
	fmt.Printf("New age after increment by 5: %d\n", newAge)

	// 10. HINCRBYFLOAT - Increment hash field by float
	fmt.Println("\n=== HINCRBYFLOAT Command ===")

	// Add a score field
	_, err = rdb.HSet(ctx, "user:1001", "score", "100.5").Result()
	if err != nil {
		log.Fatalf("Error setting score: %v", err)
	}

	newScore, err := rdb.HIncrByFloat(ctx, "user:1001", "score", 15.3).Result()
	if err != nil {
		log.Fatalf("Error incrementing score: %v", err)
	}
	fmt.Printf("New score after increment: %.2f\n", newScore)

	// 11. HLEN - Get number of fields in hash
	fmt.Println("\n=== HLEN Command ===")

	length, err := rdb.HLen(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error getting hash length: %v", err)
	}
	fmt.Printf("Number of fields in hash: %d\n", length)

	// 12. HSETNX - Set field only if it doesn't exist
	fmt.Println("\n=== HSETNX Command ===")

	// Try to set existing field (should fail)
	set, err := rdb.HSetNX(ctx, "user:1001", "name", "Bob").Result()
	if err != nil {
		log.Fatalf("Error with HSETNX on existing field: %v", err)
	}
	fmt.Printf("HSETNX on existing field 'name': %v\n", set)

	// Try to set new field (should succeed)
	set, err = rdb.HSetNX(ctx, "user:1001", "phone", "123-456-7890").Result()
	if err != nil {
		log.Fatalf("Error with HSETNX on new field: %v", err)
	}
	fmt.Printf("HSETNX on new field 'phone': %v\n", set)

	// 13. HMSET - Set multiple fields (deprecated but still works)
	fmt.Println("\n=== HMSET Command (deprecated) ===")

	// Note: HMSET is deprecated, use HSET instead
	fields2 := map[string]interface{}{
		"department": "Engineering",
		"role":       "Developer",
		"salary":     "75000",
	}
	_, err = rdb.HSet(ctx, "user:1001", fields2).Result()
	if err != nil {
		log.Fatalf("Error setting additional fields: %v", err)
	}
	fmt.Println("Set additional user fields")

	// 14. Final hash state
	fmt.Println("\n=== Final Hash State ===")

	allFields, err = rdb.HGetAll(ctx, "user:1001").Result()
	if err != nil {
		log.Fatalf("Error getting final hash state: %v", err)
	}
	fmt.Println("Final user profile:")
	for field, value := range allFields {
		fmt.Printf("  %s: %s\n", field, value)
	}

	// 15. Working with multiple users
	fmt.Println("\n=== Multiple Users Example ===")

	// Create another user
	user2Fields := map[string]interface{}{
		"name":   "Bob",
		"age":    25,
		"email":  "bob@example.com",
		"city":   "San Francisco",
		"role":   "Designer",
		"salary": "65000",
	}
	_, err = rdb.HSet(ctx, "user:1002", user2Fields).Result()
	if err != nil {
		log.Fatalf("Error creating user 1002: %v", err)
	}

	// Get all users
	userKeys, err := rdb.Keys(ctx, "user:*").Result()
	if err != nil {
		log.Fatalf("Error getting user keys: %v", err)
	}

	fmt.Printf("Found %d users:\n", len(userKeys))
	for _, key := range userKeys {
		userData, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			log.Fatalf("Error getting user data for %s: %v", key, err)
		}
		fmt.Printf("\n%s:\n", key)
		for field, value := range userData {
			fmt.Printf("  %s: %s\n", field, value)
		}
	}
}
