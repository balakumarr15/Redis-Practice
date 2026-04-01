package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis Transactions operations
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

	// 1. Basic Transaction with MULTI/EXEC
	fmt.Println("\n=== Basic Transaction ===")

	// Start transaction
	pipe := rdb.TxPipeline()

	// Add commands to transaction
	pipe.Set(ctx, "tx:key1", "value1", 0)
	pipe.Set(ctx, "tx:key2", "value2", 0)
	pipe.HSet(ctx, "tx:hash", "field1", "value1")
	pipe.LPush(ctx, "tx:list", "item1", "item2")

	// Execute transaction
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error executing transaction: %v", err)
	}

	fmt.Printf("Transaction executed %d commands\n", len(cmds))

	// Verify results
	val1, _ := rdb.Get(ctx, "tx:key1").Result()
	val2, _ := rdb.Get(ctx, "tx:key2").Result()
	hashVal, _ := rdb.HGet(ctx, "tx:hash", "field1").Result()
	listLen, _ := rdb.LLen(ctx, "tx:list").Result()

	fmt.Printf("tx:key1 = %s\n", val1)
	fmt.Printf("tx:key2 = %s\n", val2)
	fmt.Printf("tx:hash field1 = %s\n", hashVal)
	fmt.Printf("tx:list length = %d\n", listLen)

	// 2. Transaction with WATCH
	fmt.Println("\n=== Transaction with WATCH ===")

	// Set initial value
	rdb.Set(ctx, "counter", "10", 0)

	// Watch the key
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		// Get current value
		_, err := tx.Get(ctx, "counter").Result()
		if err != nil && err != redis.Nil {
			return err
		}

		// Increment in transaction
		pipe := tx.TxPipeline()
		pipe.Incr(ctx, "counter")
		pipe.Set(ctx, "last_updated", time.Now().Unix(), 0)

		_, err = pipe.Exec(ctx)
		return err
	}, "counter")

	if err != nil {
		log.Fatalf("Error in watched transaction: %v", err)
	}

	// Check result
	counter, _ := rdb.Get(ctx, "counter").Result()
	lastUpdated, _ := rdb.Get(ctx, "last_updated").Result()
	fmt.Printf("Counter after transaction: %s\n", counter)
	fmt.Printf("Last updated: %s\n", lastUpdated)

	// 3. Transaction with conditional operations
	fmt.Println("\n=== Transaction with Conditional Operations ===")

	// Set up initial state
	rdb.Set(ctx, "balance:user1", "100", 0)
	rdb.Set(ctx, "balance:user2", "50", 0)

	// Transfer money with conditions
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		// Get balances
		balance1, err := tx.Get(ctx, "balance:user1").Result()
		if err != nil {
			return err
		}

		_, err = tx.Get(ctx, "balance:user2").Result()
		if err != nil {
			return err
		}

		// Check if user1 has enough balance
		if balance1 < "30" {
			return fmt.Errorf("insufficient balance")
		}

		// Transfer 30 from user1 to user2
		pipe := tx.TxPipeline()
		pipe.DecrBy(ctx, "balance:user1", 30)
		pipe.IncrBy(ctx, "balance:user2", 30)
		pipe.Set(ctx, "transfer_log", fmt.Sprintf("Transferred 30 at %d", time.Now().Unix()), 0)

		_, err = pipe.Exec(ctx)
		return err
	}, "balance:user1", "balance:user2")

	if err != nil {
		log.Fatalf("Error in transfer transaction: %v", err)
	}

	// Check final balances
	finalBalance1, _ := rdb.Get(ctx, "balance:user1").Result()
	finalBalance2, _ := rdb.Get(ctx, "balance:user2").Result()
	transferLog, _ := rdb.Get(ctx, "transfer_log").Result()

	fmt.Printf("User1 balance: %s\n", finalBalance1)
	fmt.Printf("User2 balance: %s\n", finalBalance2)
	fmt.Printf("Transfer log: %s\n", transferLog)

	// 4. Transaction with DISCARD
	fmt.Println("\n=== Transaction with DISCARD ===")

	// Start transaction
	pipe = rdb.TxPipeline()

	// Add commands
	pipe.Set(ctx, "discard:key1", "value1", 0)
	pipe.Set(ctx, "discard:key2", "value2", 0)

	// Discard transaction (simulate by not executing)
	fmt.Println("Transaction prepared but not executed (simulated DISCARD)")

	// Check if keys exist (they shouldn't)
	exists1, _ := rdb.Exists(ctx, "discard:key1").Result()
	exists2, _ := rdb.Exists(ctx, "discard:key2").Result()
	fmt.Printf("discard:key1 exists: %v\n", exists1 > 0)
	fmt.Printf("discard:key2 exists: %v\n", exists2 > 0)

	// 5. Complex transaction example - Inventory management
	fmt.Println("\n=== Inventory Management Transaction ===")

	// Set up inventory
	rdb.HSet(ctx, "inventory:item1", map[string]interface{}{
		"name":     "Laptop",
		"quantity": "10",
		"price":    "999.99",
	})
	rdb.HSet(ctx, "inventory:item2", map[string]interface{}{
		"name":     "Mouse",
		"quantity": "50",
		"price":    "29.99",
	})

	// Process order
	orderItems := map[string]int{
		"item1": 2, // 2 laptops
		"item2": 5, // 5 mice
	}

	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		// Check inventory for all items
		for itemID, requestedQty := range orderItems {
			quantity, err := tx.HGet(ctx, fmt.Sprintf("inventory:%s", itemID), "quantity").Result()
			if err != nil {
				return err
			}

			if quantity < fmt.Sprintf("%d", requestedQty) {
				return fmt.Errorf("insufficient inventory for %s", itemID)
			}
		}

		// Process order
		pipe := tx.TxPipeline()
		for itemID, requestedQty := range orderItems {
			pipe.HIncrBy(ctx, fmt.Sprintf("inventory:%s", itemID), "quantity", -int64(requestedQty))
		}

		// Create order record
		orderID := fmt.Sprintf("order:%d", time.Now().Unix())
		pipe.HSet(ctx, orderID, map[string]interface{}{
			"status":    "processed",
			"timestamp": time.Now().Unix(),
			"items":     fmt.Sprintf("%v", orderItems),
		})

		_, err = pipe.Exec(ctx)
		return err
	}, "inventory:item1", "inventory:item2")

	if err != nil {
		log.Fatalf("Error in inventory transaction: %v", err)
	}

	// Check final inventory
	item1Qty, _ := rdb.HGet(ctx, "inventory:item1", "quantity").Result()
	item2Qty, _ := rdb.HGet(ctx, "inventory:item2", "quantity").Result()
	fmt.Printf("Item1 remaining quantity: %s\n", item1Qty)
	fmt.Printf("Item2 remaining quantity: %s\n", item2Qty)

	// 6. Transaction with error handling
	fmt.Println("\n=== Transaction with Error Handling ===")

	// Set up test data
	rdb.Set(ctx, "error:key1", "value1", 0)

	// Transaction that will fail
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		pipe := tx.TxPipeline()

		// Valid operation
		pipe.Set(ctx, "error:key2", "value2", 0)

		// Invalid operation (will cause error)
		pipe.HGet(ctx, "nonexistent_hash", "field")

		_, err = pipe.Exec(ctx)
		return err
	}, "error:key1")

	if err != nil {
		fmt.Printf("Transaction failed as expected: %v\n", err)
	}

	// Check if valid operation was executed (it shouldn't be)
	exists, _ := rdb.Exists(ctx, "error:key2").Result()
	fmt.Printf("error:key2 exists after failed transaction: %v\n", exists > 0)

	// 7. Performance comparison
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

	// Transaction commands
	start = time.Now()
	pipe = rdb.TxPipeline()
	for i := 0; i < 100; i++ {
		pipe.Set(ctx, fmt.Sprintf("transaction:%d", i), fmt.Sprintf("value%d", i), 0)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Error in transaction commands: %v", err)
	}
	transactionTime := time.Since(start)
	fmt.Printf("100 transaction commands took: %v\n", transactionTime)

	// 8. Cleanup
	fmt.Println("\n=== Cleanup ===")

	keys := []string{
		"tx:key1", "tx:key2", "tx:hash", "tx:list",
		"counter", "last_updated",
		"balance:user1", "balance:user2", "transfer_log",
		"inventory:item1", "inventory:item2",
		"error:key1", "error:key2",
	}

	// Add individual and transaction keys
	for i := 0; i < 100; i++ {
		keys = append(keys, fmt.Sprintf("individual:%d", i))
		keys = append(keys, fmt.Sprintf("transaction:%d", i))
	}

	deleted, err := rdb.Del(ctx, keys...).Result()
	if err != nil {
		log.Fatalf("Error cleaning up: %v", err)
	}
	fmt.Printf("Cleaned up %d keys\n", deleted)
}
