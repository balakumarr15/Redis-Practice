package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis Streams operations
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

	// 1. XADD - Add entries to stream
	fmt.Println("\n=== XADD Command ===")

	// Add single entry
	entryID, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "events",
		Values: map[string]interface{}{
			"event":     "user_login",
			"user_id":   "123",
			"timestamp": time.Now().Unix(),
		},
	}).Result()
	if err != nil {
		log.Fatalf("Error adding entry: %v", err)
	}
	fmt.Printf("Added entry with ID: %s\n", entryID)

	// Add multiple entries
	for i := 1; i <= 3; i++ {
		entryID, err = rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "events",
			Values: map[string]interface{}{
				"event":     "page_view",
				"page":      fmt.Sprintf("/page%d", i),
				"user_id":   "123",
				"timestamp": time.Now().Unix(),
			},
		}).Result()
		if err != nil {
			log.Fatalf("Error adding entry %d: %v", i, err)
		}
		fmt.Printf("Added entry %d with ID: %s\n", i, entryID)
		time.Sleep(10 * time.Millisecond)
	}

	// 2. XREAD - Read entries from stream
	fmt.Println("\n=== XREAD Command ===")

	// Read all entries from beginning
	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"events", "0"},
		Count:   10,
	}).Result()
	if err != nil {
		log.Fatalf("Error reading stream: %v", err)
	}

	for _, stream := range streams {
		fmt.Printf("Stream: %s\n", stream.Stream)
		for _, message := range stream.Messages {
			fmt.Printf("  ID: %s\n", message.ID)
			for field, value := range message.Values {
				fmt.Printf("    %s: %s\n", field, value)
			}
		}
	}

	// 3. XRANGE - Get entries by ID range
	fmt.Println("\n=== XRANGE Command ===")

	// Get all entries
	entries, err := rdb.XRange(ctx, "events", "-", "+").Result()
	if err != nil {
		log.Fatalf("Error getting range: %v", err)
	}
	fmt.Printf("Total entries: %d\n", len(entries))

	// Get entries with limit
	limitedEntries, err := rdb.XRangeN(ctx, "events", "-", "+", 2).Result()
	if err != nil {
		log.Fatalf("Error getting limited range: %v", err)
	}
	fmt.Printf("Limited entries: %d\n", len(limitedEntries))

	// 4. XREVRANGE - Get entries in reverse order
	fmt.Println("\n=== XREVRANGE Command ===")

	// Get last 2 entries
	revEntries, err := rdb.XRevRangeN(ctx, "events", "+", "-", 2).Result()
	if err != nil {
		log.Fatalf("Error getting reverse range: %v", err)
	}
	fmt.Println("Last 2 entries:")
	for _, entry := range revEntries {
		fmt.Printf("  ID: %s\n", entry.ID)
		for field, value := range entry.Values {
			fmt.Printf("    %s: %s\n", field, value)
		}
	}

	// 5. XGROUP - Create consumer group
	fmt.Println("\n=== XGROUP Command ===")

	// Create consumer group
	err = rdb.XGroupCreate(ctx, "events", "processors", "0").Err()
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	fmt.Println("Created consumer group 'processors'")

	// 6. XREADGROUP - Read from consumer group
	fmt.Println("\n=== XREADGROUP Command ===")

	// Read from consumer group
	groupStreams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "processors",
		Consumer: "consumer1",
		Streams:  []string{"events", ">"},
		Count:    2,
	}).Result()
	if err != nil {
		log.Fatalf("Error reading from group: %v", err)
	}

	for _, stream := range groupStreams {
		fmt.Printf("Group stream: %s\n", stream.Stream)
		for _, message := range stream.Messages {
			fmt.Printf("  ID: %s\n", message.ID)
			for field, value := range message.Values {
				fmt.Printf("    %s: %s\n", field, value)
			}
		}
	}

	// 7. XACK - Acknowledge processed messages
	fmt.Println("\n=== XACK Command ===")

	// Acknowledge processed messages
	ackCount, err := rdb.XAck(ctx, "events", "processors", groupStreams[0].Messages[0].ID).Result()
	if err != nil {
		log.Fatalf("Error acknowledging message: %v", err)
	}
	fmt.Printf("Acknowledged %d messages\n", ackCount)

	// 8. XPENDING - Check pending messages
	fmt.Println("\n=== XPENDING Command ===")

	// Get pending messages info
	pending, err := rdb.XPending(ctx, "events", "processors").Result()
	if err != nil {
		log.Fatalf("Error getting pending info: %v", err)
	}
	fmt.Printf("Pending messages: %d\n", pending.Count)
	fmt.Printf("Min ID: %s\n", pending.Lower)
	fmt.Printf("Max ID: %s\n", pending.Higher)

	// 9. XCLAIM - Claim pending messages
	fmt.Println("\n=== XCLAIM Command ===")

	// Add more entries for claiming demo
	for i := 1; i <= 2; i++ {
		_, err = rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "events",
			Values: map[string]interface{}{
				"event":    "order_created",
				"order_id": fmt.Sprintf("order_%d", i),
				"amount":   fmt.Sprintf("%d.99", i*100),
			},
		}).Result()
		if err != nil {
			log.Fatalf("Error adding order entry: %v", err)
		}
	}

	// Read with another consumer
	groupStreams2, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "processors",
		Consumer: "consumer2",
		Streams:  []string{"events", ">"},
		Count:    2,
	}).Result()
	if err != nil {
		log.Fatalf("Error reading with consumer2: %v", err)
	}

	// Claim messages from consumer2 to consumer1
	if len(groupStreams2) > 0 && len(groupStreams2[0].Messages) > 0 {
		claimed, err := rdb.XClaim(ctx, &redis.XClaimArgs{
			Stream:   "events",
			Group:    "processors",
			Consumer: "consumer1",
			MinIdle:  time.Millisecond,
			Messages: []string{groupStreams2[0].Messages[0].ID},
		}).Result()
		if err != nil {
			log.Fatalf("Error claiming messages: %v", err)
		}
		fmt.Printf("Claimed %d messages\n", len(claimed))
	}

	// 10. XDEL - Delete entries from stream
	fmt.Println("\n=== XDEL Command ===")

	// Get first entry ID for deletion
	if len(entries) > 0 {
		deleted, err := rdb.XDel(ctx, "events", entries[0].ID).Result()
		if err != nil {
			log.Fatalf("Error deleting entry: %v", err)
		}
		fmt.Printf("Deleted %d entries\n", deleted)
	}

	// 11. XTRIM - Trim stream to specified length
	fmt.Println("\n=== XTRIM Command ===")

	// Trim stream to keep only last 5 entries
	trimmed, err := rdb.XTrimMaxLen(ctx, "events", 5).Result()
	if err != nil {
		log.Fatalf("Error trimming stream: %v", err)
	}
	fmt.Printf("Trimmed %d entries\n", trimmed)

	// 12. XLEN - Get stream length
	fmt.Println("\n=== XLEN Command ===")

	length, err := rdb.XLen(ctx, "events").Result()
	if err != nil {
		log.Fatalf("Error getting stream length: %v", err)
	}
	fmt.Printf("Stream length: %d\n", length)

	// 13. Practical example - Event sourcing
	fmt.Println("\n=== Event Sourcing Example ===")

	// Create event store
	eventStore := "user_events"

	// Add user events
	userEvents := []map[string]interface{}{
		{
			"event_type": "user_created",
			"user_id":    "user123",
			"email":      "user@example.com",
			"timestamp":  time.Now().Unix(),
		},
		{
			"event_type": "profile_updated",
			"user_id":    "user123",
			"name":       "John Doe",
			"timestamp":  time.Now().Unix(),
		},
		{
			"event_type": "email_changed",
			"user_id":    "user123",
			"old_email":  "user@example.com",
			"new_email":  "john@example.com",
			"timestamp":  time.Now().Unix(),
		},
	}

	for i, event := range userEvents {
		entryID, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: eventStore,
			Values: event,
		}).Result()
		if err != nil {
			log.Fatalf("Error adding user event %d: %v", i, err)
		}
		fmt.Printf("Added user event %d with ID: %s\n", i+1, entryID)
		time.Sleep(10 * time.Millisecond)
	}

	// Replay events
	fmt.Println("\nReplaying user events:")
	replayEntries, err := rdb.XRange(ctx, eventStore, "-", "+").Result()
	if err != nil {
		log.Fatalf("Error replaying events: %v", err)
	}

	for _, entry := range replayEntries {
		fmt.Printf("Event ID: %s\n", entry.ID)
		for field, value := range entry.Values {
			fmt.Printf("  %s: %s\n", field, value)
		}
		fmt.Println()
	}

	// 14. Cleanup
	fmt.Println("\n=== Cleanup ===")

	// Delete streams
	_, err = rdb.Del(ctx, "events", eventStore).Result()
	if err != nil {
		log.Fatalf("Error cleaning up streams: %v", err)
	}
	fmt.Println("Streams deleted")
}
