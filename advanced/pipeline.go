package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis Pipeline operations for batch processing
func PipelineOperations() {
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

	// 1. Basic Pipeline
	fmt.Println("\n=== Basic Pipeline ===")

	// Create pipeline
	pipe := rdb.Pipeline()

	// Add commands to pipeline
	pipe.Set(ctx, "key1", "value1", 0)
	pipe.Set(ctx, "key2", "value2", 0)
	pipe.Set(ctx, "key3", "value3", 0)
	pipe.Get(ctx, "key1")
	pipe.Get(ctx, "key2")
	pipe.Get(ctx, "key3")

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing pipeline: %v", err)
	}

	// Process results
	fmt.Printf("Executed %d commands\n", len(cmds))
	for i, cmd := range cmds {
		if i < 3 {
			// SET commands
			fmt.Printf("Command %d (SET): %v\n", i+1, cmd.Err())
		} else {
			// GET commands
			val, _ := cmd.(*redis.StringCmd).Result()
			fmt.Printf("Command %d (GET): %s\n", i+1, val)
		}
	}

	// 2. Pipeline with different data types
	fmt.Println("\n=== Pipeline with Different Data Types ===")

	pipe = rdb.Pipeline()

	// String operations
	pipe.Set(ctx, "user:name", "Alice", 0)
	pipe.Set(ctx, "user:age", "30", 0)

	// Hash operations
	pipe.HSet(ctx, "user:profile", "email", "alice@example.com")
	pipe.HSet(ctx, "user:profile", "city", "New York")

	// List operations
	pipe.LPush(ctx, "user:hobbies", "reading", "swimming", "coding")

	// Set operations
	pipe.SAdd(ctx, "user:tags", "developer", "golang", "redis")

	// Sorted set operations
	pipe.ZAdd(ctx, "user:scores", redis.Z{Score: 100, Member: "math"})
	pipe.ZAdd(ctx, "user:scores", redis.Z{Score: 95, Member: "science"})

	// Execute pipeline
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing mixed pipeline: %v", err)
	}

	fmt.Printf("Executed %d mixed commands\n", len(cmds))

	// 3. Pipeline with error handling
	fmt.Println("\n=== Pipeline with Error Handling ===")

	pipe = rdb.Pipeline()

	// Valid commands
	pipe.Set(ctx, "valid_key", "valid_value", 0)
	pipe.Get(ctx, "valid_key")

	// Invalid command (will cause error)
	pipe.HGet(ctx, "nonexistent_hash", "field")

	// Execute pipeline
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		fmt.Printf("Pipeline execution error: %v\n", err)
	}

	// Check individual command results
	for i, cmd := range cmds {
		if cmd.Err() != nil {
			fmt.Printf("Command %d error: %v\n", i+1, cmd.Err())
		} else {
			fmt.Printf("Command %d success\n", i+1)
		}
	}

	// 4. Batch operations with Pipeline
	fmt.Println("\n=== Batch Operations ===")

	// Batch insert users
	pipe = rdb.Pipeline()

	users := []struct {
		ID    string
		Name  string
		Email string
		Age   int
	}{
		{"user:1", "Alice", "alice@example.com", 30},
		{"user:2", "Bob", "bob@example.com", 25},
		{"user:3", "Charlie", "charlie@example.com", 35},
		{"user:4", "Diana", "diana@example.com", 28},
		{"user:5", "Eve", "eve@example.com", 32},
	}

	for _, user := range users {
		pipe.HSet(ctx, user.ID, map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		})
	}

	// Execute batch insert
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing batch insert: %v", err)
	}

	fmt.Printf("Batch inserted %d users\n", len(cmds))

	// 5. Batch read operations
	fmt.Println("\n=== Batch Read Operations ===")

	pipe = rdb.Pipeline()

	// Read all users
	for _, user := range users {
		pipe.HGetAll(ctx, user.ID)
	}

	// Execute batch read
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing batch read: %v", err)
	}

	// Process read results
	fmt.Println("Batch read results:")
	for i, cmd := range cmds {
		userData, err := cmd.(*redis.MapStringStringCmd).Result()
		if err != nil {
			fmt.Printf("User %d read error: %v\n", i+1, err)
		} else {
			fmt.Printf("User %d: %v\n", i+1, userData)
		}
	}

	// 6. Performance comparison
	fmt.Println("\n=== Performance Comparison ===")

	// Individual commands
	start := time.Now()
	for i := 0; i < 100; i++ {
		_, err = rdb.Set(ctx, fmt.Sprintf("individual:%d", i), fmt.Sprintf("value%d", i), 0).Result()
		if err != nil {
			log.Fatalf("Error in individual command: %v", err)
		}
	}
	individualTime := time.Since(start)
	fmt.Printf("100 individual commands took: %v\n", individualTime)

	// Pipeline commands
	start = time.Now()
	pipe = rdb.Pipeline()
	for i := 0; i < 100; i++ {
		pipe.Set(ctx, fmt.Sprintf("pipeline:%d", i), fmt.Sprintf("value%d", i), 0)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error in pipeline commands: %v", err)
	}
	pipelineTime := time.Since(start)
	fmt.Printf("100 pipeline commands took: %v\n", pipelineTime)

	// Calculate speedup
	speedup := float64(individualTime) / float64(pipelineTime)
	fmt.Printf("Pipeline speedup: %.2fx\n", speedup)

	// 7. Pipeline with transactions
	fmt.Println("\n=== Pipeline with Transactions ===")

	// Start transaction
	tx := rdb.TxPipeline()

	// Add commands to transaction
	tx.Set(ctx, "tx:key1", "tx:value1", 0)
	tx.Set(ctx, "tx:key2", "tx:value2", 0)
	tx.HSet(ctx, "tx:hash", "field1", "value1")
	tx.LPush(ctx, "tx:list", "item1", "item2")

	// Execute transaction
	cmds, err = tx.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing transaction: %v", err)
	}

	fmt.Printf("Transaction executed %d commands\n", len(cmds))

	// 8. Pipeline with conditional operations
	fmt.Println("\n=== Pipeline with Conditional Operations ===")

	pipe = rdb.Pipeline()

	// Conditional operations
	pipe.SetNX(ctx, "conditional:key1", "value1", 0)
	pipe.SetNX(ctx, "conditional:key1", "value2", 0) // Should fail
	pipe.HSetNX(ctx, "conditional:hash", "field1", "value1")
	pipe.HSetNX(ctx, "conditional:hash", "field1", "value2") // Should fail

	// Execute pipeline
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing conditional pipeline: %v", err)
	}

	// Check results
	for i, cmd := range cmds {
		result, _ := cmd.(*redis.BoolCmd).Result()
		fmt.Printf("Conditional command %d result: %v\n", i+1, result)
	}

	// 9. Pipeline with expiration
	fmt.Println("\n=== Pipeline with Expiration ===")

	pipe = rdb.Pipeline()

	// Set keys with different expiration times
	pipe.Set(ctx, "expire:1s", "value1", 1*time.Second)
	pipe.Set(ctx, "expire:5s", "value2", 5*time.Second)
	pipe.Set(ctx, "expire:10s", "value3", 10*time.Second)

	// Check TTL
	pipe.TTL(ctx, "expire:1s")
	pipe.TTL(ctx, "expire:5s")
	pipe.TTL(ctx, "expire:10s")

	// Execute pipeline
	cmds, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing expiration pipeline: %v", err)
	}

	// Process TTL results
	for i := 3; i < len(cmds); i++ {
		ttl, _ := cmds[i].(*redis.DurationCmd).Result()
		fmt.Printf("TTL for expire key %d: %v\n", i-2, ttl)
	}

	// 10. Cleanup
	fmt.Println("\n=== Cleanup ===")

	// Clean up all test keys
	keys := []string{
		"key1", "key2", "key3", "valid_key",
		"user:name", "user:age", "user:profile", "user:hobbies", "user:tags", "user:scores",
		"user:1", "user:2", "user:3", "user:4", "user:5",
		"tx:key1", "tx:key2", "tx:hash", "tx:list",
		"conditional:key1", "conditional:hash",
		"expire:1s", "expire:5s", "expire:10s",
	}

	// Add individual and pipeline keys
	for i := 0; i < 100; i++ {
		keys = append(keys, fmt.Sprintf("individual:%d", i))
		keys = append(keys, fmt.Sprintf("pipeline:%d", i))
	}

	// Delete all keys
	deleted, err := rdb.Del(ctx, keys...).Result()
	if err != nil {
		log.Fatalf("Error cleaning up: %v", err)
	}
	fmt.Printf("Cleaned up %d keys\n", deleted)
}
