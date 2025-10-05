package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis List operations - LPUSH, RPUSH, LPOP, RPOP, etc.
func ListOperations() {
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

	// 1. LPUSH - Push elements to the left (head) of list
	fmt.Println("\n=== LPUSH Command ===")

	// Push single element
	length, err := rdb.LPush(ctx, "tasks", "task1").Result()
	if err != nil {
		log.Fatalf("Error pushing to list: %v", err)
	}
	fmt.Printf("LPUSH single element, list length: %d\n", length)

	// Push multiple elements
	length, err = rdb.LPush(ctx, "tasks", "task2", "task3", "task4").Result()
	if err != nil {
		log.Fatalf("Error pushing multiple elements: %v", err)
	}
	fmt.Printf("LPUSH multiple elements, list length: %d\n", length)

	// 2. RPUSH - Push elements to the right (tail) of list
	fmt.Println("\n=== RPUSH Command ===")

	length, err = rdb.RPush(ctx, "tasks", "task5", "task6").Result()
	if err != nil {
		log.Fatalf("Error pushing to right: %v", err)
	}
	fmt.Printf("RPUSH elements, list length: %d\n", length)

	// 3. LRANGE - Get range of elements from list
	fmt.Println("\n=== LRANGE Command ===")

	// Get all elements
	allTasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting all elements: %v", err)
	}
	fmt.Printf("All tasks: %v\n", allTasks)

	// Get first 3 elements
	firstThree, err := rdb.LRange(ctx, "tasks", 0, 2).Result()
	if err != nil {
		log.Fatalf("Error getting first 3 elements: %v", err)
	}
	fmt.Printf("First 3 tasks: %v\n", firstThree)

	// Get last 2 elements
	lastTwo, err := rdb.LRange(ctx, "tasks", -2, -1).Result()
	if err != nil {
		log.Fatalf("Error getting last 2 elements: %v", err)
	}
	fmt.Printf("Last 2 tasks: %v\n", lastTwo)

	// 4. LLEN - Get list length
	fmt.Println("\n=== LLEN Command ===")

	length, err = rdb.LLen(ctx, "tasks").Result()
	if err != nil {
		log.Fatalf("Error getting list length: %v", err)
	}
	fmt.Printf("List length: %d\n", length)

	// 5. LPOP - Pop element from left (head)
	fmt.Println("\n=== LPOP Command ===")

	popped, err := rdb.LPop(ctx, "tasks").Result()
	if err != nil {
		log.Fatalf("Error popping from left: %v", err)
	}
	fmt.Printf("Popped from left: %s\n", popped)

	// Show remaining list
	remaining, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting remaining elements: %v", err)
	}
	fmt.Printf("Remaining tasks: %v\n", remaining)

	// 6. RPOP - Pop element from right (tail)
	fmt.Println("\n=== RPOP Command ===")

	popped, err = rdb.RPop(ctx, "tasks").Result()
	if err != nil {
		log.Fatalf("Error popping from right: %v", err)
	}
	fmt.Printf("Popped from right: %s\n", popped)

	// Show remaining list
	remaining, err = rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting remaining elements: %v", err)
	}
	fmt.Printf("Remaining tasks: %v\n", remaining)

	// 7. LINDEX - Get element at specific index
	fmt.Println("\n=== LINDEX Command ===")

	// Get first element (index 0)
	first, err := rdb.LIndex(ctx, "tasks", 0).Result()
	if err != nil {
		log.Fatalf("Error getting first element: %v", err)
	}
	fmt.Printf("First element (index 0): %s\n", first)

	// Get last element (index -1)
	last, err := rdb.LIndex(ctx, "tasks", -1).Result()
	if err != nil {
		log.Fatalf("Error getting last element: %v", err)
	}
	fmt.Printf("Last element (index -1): %s\n", last)

	// 8. LINSERT - Insert element before or after pivot
	fmt.Println("\n=== LINSERT Command ===")

	// Insert before "task3"
	length, err = rdb.LInsertBefore(ctx, "tasks", "task3", "urgent_task").Result()
	if err != nil {
		log.Fatalf("Error inserting before: %v", err)
	}
	fmt.Printf("Inserted before 'task3', new length: %d\n", length)

	// Insert after "task4"
	length, err = rdb.LInsertAfter(ctx, "tasks", "task4", "follow_up_task").Result()
	if err != nil {
		log.Fatalf("Error inserting after: %v", err)
	}
	fmt.Printf("Inserted after 'task4', new length: %d\n", length)

	// Show updated list
	updated, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting updated list: %v", err)
	}
	fmt.Printf("Updated tasks: %v\n", updated)

	// 9. LREM - Remove elements from list
	fmt.Println("\n=== LREM Command ===")

	// Add some duplicate elements
	_, err = rdb.RPush(ctx, "tasks", "duplicate", "duplicate", "duplicate").Result()
	if err != nil {
		log.Fatalf("Error adding duplicates: %v", err)
	}

	// Remove 2 occurrences of "duplicate" (from left)
	removed, err := rdb.LRem(ctx, "tasks", 2, "duplicate").Result()
	if err != nil {
		log.Fatalf("Error removing duplicates: %v", err)
	}
	fmt.Printf("Removed %d occurrences of 'duplicate'\n", removed)

	// Show list after removal
	afterRemoval, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting list after removal: %v", err)
	}
	fmt.Printf("Tasks after removal: %v\n", afterRemoval)

	// 10. LSET - Set element at specific index
	fmt.Println("\n=== LSET Command ===")

	// Set element at index 1
	err = rdb.LSet(ctx, "tasks", 1, "updated_task").Err()
	if err != nil {
		log.Fatalf("Error setting element at index: %v", err)
	}
	fmt.Println("Set element at index 1 to 'updated_task'")

	// Show updated list
	afterSet, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting list after set: %v", err)
	}
	fmt.Printf("Tasks after LSET: %v\n", afterSet)

	// 11. LTRIM - Trim list to specified range
	fmt.Println("\n=== LTRIM Command ===")

	// Keep only elements from index 1 to 3
	err = rdb.LTrim(ctx, "tasks", 1, 3).Err()
	if err != nil {
		log.Fatalf("Error trimming list: %v", err)
	}
	fmt.Println("Trimmed list to keep elements from index 1 to 3")

	// Show trimmed list
	trimmed, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting trimmed list: %v", err)
	}
	fmt.Printf("Trimmed tasks: %v\n", trimmed)

	// 12. RPOPLPUSH - Pop from right of source, push to left of destination
	fmt.Println("\n=== RPOPLPUSH Command ===")

	// Create a completed tasks list
	_, err = rdb.RPush(ctx, "completed_tasks", "old_completed").Result()
	if err != nil {
		log.Fatalf("Error creating completed tasks list: %v", err)
	}

	// Move task from tasks to completed_tasks
	moved, err := rdb.RPopLPush(ctx, "tasks", "completed_tasks").Result()
	if err != nil {
		log.Fatalf("Error moving task: %v", err)
	}
	fmt.Printf("Moved task '%s' from tasks to completed_tasks\n", moved)

	// Show both lists
	tasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting tasks: %v", err)
	}
	completed, err := rdb.LRange(ctx, "completed_tasks", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting completed tasks: %v", err)
	}
	fmt.Printf("Remaining tasks: %v\n", tasks)
	fmt.Printf("Completed tasks: %v\n", completed)

	// 13. Blocking operations - BLPOP and BRPOP
	fmt.Println("\n=== Blocking Operations (BLPOP/BRPOP) ===")

	// Note: In a real application, you'd run this in a goroutine
	// to avoid blocking the main thread
	fmt.Println("Note: BLPOP and BRPOP are blocking operations")
	fmt.Println("They wait for elements to become available in the list")
	fmt.Println("In production, use them in separate goroutines")

	// 14. Working with multiple lists
	fmt.Println("\n=== Multiple Lists Example ===")

	// Create different priority queues
	highPriority := []string{"critical_bug", "security_fix", "urgent_feature"}
	mediumPriority := []string{"new_feature", "improvement", "refactoring"}
	lowPriority := []string{"documentation", "cleanup", "optimization"}

	// Add to high priority queue
	for _, task := range highPriority {
		_, err = rdb.LPush(ctx, "high_priority", task).Result()
		if err != nil {
			log.Fatalf("Error adding to high priority: %v", err)
		}
	}

	// Add to medium priority queue
	for _, task := range mediumPriority {
		_, err = rdb.LPush(ctx, "medium_priority", task).Result()
		if err != nil {
			log.Fatalf("Error adding to medium priority: %v", err)
		}
	}

	// Add to low priority queue
	for _, task := range lowPriority {
		_, err = rdb.LPush(ctx, "low_priority", task).Result()
		if err != nil {
			log.Fatalf("Error adding to low priority: %v", err)
		}
	}

	// Process tasks by priority
	fmt.Println("Processing tasks by priority:")

	// Process high priority first
	highTasks, err := rdb.LRange(ctx, "high_priority", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting high priority tasks: %v", err)
	}
	fmt.Printf("High priority tasks: %v\n", highTasks)

	// Process medium priority
	mediumTasks, err := rdb.LRange(ctx, "medium_priority", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting medium priority tasks: %v", err)
	}
	fmt.Printf("Medium priority tasks: %v\n", mediumTasks)

	// Process low priority
	lowTasks, err := rdb.LRange(ctx, "low_priority", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting low priority tasks: %v", err)
	}
	fmt.Printf("Low priority tasks: %v\n", lowTasks)
}
