package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestHashOperations tests Redis hash operations
func TestHashOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test HSET
	added, err := rdb.HSet(ctx, "test_hash", "field1", "value1").Result()
	if err != nil {
		t.Fatalf("Error setting hash field: %v", err)
	}

	if added != 1 {
		t.Errorf("Expected 1 added field, got %d", added)
	}

	// Test HGET
	val, err := rdb.HGet(ctx, "test_hash", "field1").Result()
	if err != nil {
		t.Fatalf("Error getting hash field: %v", err)
	}

	if val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}

	// Test HGETALL
	allFields, err := rdb.HGetAll(ctx, "test_hash").Result()
	if err != nil {
		t.Fatalf("Error getting all hash fields: %v", err)
	}

	if len(allFields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(allFields))
	}

	// Test HDEL
	deleted, err := rdb.HDel(ctx, "test_hash", "field1").Result()
	if err != nil {
		t.Fatalf("Error deleting hash field: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected 1 deleted field, got %d", deleted)
	}

	// Cleanup
	rdb.Del(ctx, "test_hash")
}

// TestListOperations tests Redis list operations
func TestListOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test LPUSH
	length, err := rdb.LPush(ctx, "test_list", "item1", "item2").Result()
	if err != nil {
		t.Fatalf("Error pushing to list: %v", err)
	}

	if length != 2 {
		t.Errorf("Expected list length 2, got %d", length)
	}

	// Test LRANGE
	items, err := rdb.LRange(ctx, "test_list", 0, -1).Result()
	if err != nil {
		t.Fatalf("Error getting list range: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Test LPOP
	popped, err := rdb.LPop(ctx, "test_list").Result()
	if err != nil {
		t.Fatalf("Error popping from list: %v", err)
	}

	if popped != "item2" {
		t.Errorf("Expected 'item2', got '%s'", popped)
	}

	// Cleanup
	rdb.Del(ctx, "test_list")
}

// TestSetOperations tests Redis set operations
func TestSetOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test SADD
	added, err := rdb.SAdd(ctx, "test_set", "member1", "member2").Result()
	if err != nil {
		t.Fatalf("Error adding to set: %v", err)
	}

	if added != 2 {
		t.Errorf("Expected 2 added members, got %d", added)
	}

	// Test SMEMBERS
	members, err := rdb.SMembers(ctx, "test_set").Result()
	if err != nil {
		t.Fatalf("Error getting set members: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}

	// Test SISMEMBER
	isMember, err := rdb.SIsMember(ctx, "test_set", "member1").Result()
	if err != nil {
		t.Fatalf("Error checking membership: %v", err)
	}

	if !isMember {
		t.Error("member1 should be a member")
	}

	// Cleanup
	rdb.Del(ctx, "test_set")
}

// TestSortedSetOperations tests Redis sorted set operations
func TestSortedSetOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test ZADD
	added, err := rdb.ZAdd(ctx, "test_zset", redis.Z{Score: 100, Member: "member1"}).Result()
	if err != nil {
		t.Fatalf("Error adding to sorted set: %v", err)
	}

	if added != 1 {
		t.Errorf("Expected 1 added member, got %d", added)
	}

	// Test ZRANGE
	members, err := rdb.ZRange(ctx, "test_zset", 0, -1).Result()
	if err != nil {
		t.Fatalf("Error getting sorted set range: %v", err)
	}

	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}

	// Test ZSCORE
	score, err := rdb.ZScore(ctx, "test_zset", "member1").Result()
	if err != nil {
		t.Fatalf("Error getting member score: %v", err)
	}

	if score != 100 {
		t.Errorf("Expected score 100, got %f", score)
	}

	// Cleanup
	rdb.Del(ctx, "test_zset")
}

// TestPipelineOperations tests Redis pipeline operations
func TestPipelineOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Create pipeline
	pipe := rdb.Pipeline()

	// Add commands to pipeline
	pipe.Set(ctx, "pipeline_key1", "value1", 0)
	pipe.Set(ctx, "pipeline_key2", "value2", 0)
	pipe.Get(ctx, "pipeline_key1")

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		t.Fatalf("Error executing pipeline: %v", err)
	}

	if len(cmds) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(cmds))
	}

	// Check results
	val, err := cmds[2].(*redis.StringCmd).Result()
	if err != nil {
		t.Fatalf("Error getting pipeline result: %v", err)
	}

	if val != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val)
	}

	// Cleanup
	rdb.Del(ctx, "pipeline_key1", "pipeline_key2")
}

// TestTransactionOperations tests Redis transaction operations
func TestTransactionOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Start transaction
	pipe := rdb.TxPipeline()

	// Add commands to transaction
	pipe.Set(ctx, "tx_key1", "value1", 0)
	pipe.Set(ctx, "tx_key2", "value2", 0)

	// Execute transaction
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		t.Fatalf("Error executing transaction: %v", err)
	}

	if len(cmds) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(cmds))
	}

	// Verify results
	val1, err := rdb.Get(ctx, "tx_key1").Result()
	if err != nil {
		t.Fatalf("Error getting tx_key1: %v", err)
	}

	if val1 != "value1" {
		t.Errorf("Expected 'value1', got '%s'", val1)
	}

	// Cleanup
	rdb.Del(ctx, "tx_key1", "tx_key2")
}

// TestPubSubOperations tests Redis Pub/Sub operations
func TestPubSubOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Subscribe to channel
	subscriber := rdb.Subscribe(ctx, "test_channel")
	defer subscriber.Close()

	// Wait for subscription
	_, err := subscriber.Receive(ctx)
	if err != nil {
		t.Fatalf("Error subscribing: %v", err)
	}

	// Publish message
	err = rdb.Publish(ctx, "test_channel", "test_message").Err()
	if err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}

	// Receive message
	msg, err := subscriber.ReceiveMessage(ctx)
	if err != nil {
		t.Fatalf("Error receiving message: %v", err)
	}

	if msg.Channel != "test_channel" {
		t.Errorf("Expected channel 'test_channel', got '%s'", msg.Channel)
	}

	if msg.Payload != "test_message" {
		t.Errorf("Expected payload 'test_message', got '%s'", msg.Payload)
	}
}

// TestStreamOperations tests Redis stream operations
func TestStreamOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Add entry to stream
	entryID, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "test_stream",
		Values: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}).Result()
	if err != nil {
		t.Fatalf("Error adding to stream: %v", err)
	}

	if entryID == "" {
		t.Error("Expected non-empty entry ID")
	}

	// Read from stream
	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"test_stream", "0"},
		Count:   1,
	}).Result()
	if err != nil {
		t.Fatalf("Error reading from stream: %v", err)
	}

	if len(streams) != 1 {
		t.Errorf("Expected 1 stream, got %d", len(streams))
	}

	if len(streams[0].Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(streams[0].Messages))
	}

	// Cleanup
	rdb.Del(ctx, "test_stream")
}

// TestErrorHandling tests error handling
func TestErrorHandling(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test getting non-existent key
	_, err := rdb.Get(ctx, "nonexistent_key").Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil error, got %v", err)
	}

	// Test getting non-existent hash field
	_, err = rdb.HGet(ctx, "nonexistent_hash", "field").Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil error, got %v", err)
	}

	// Test getting non-existent list element
	_, err = rdb.LIndex(ctx, "nonexistent_list", 0).Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil error, got %v", err)
	}
}

// TestConcurrentOperations tests concurrent operations
func TestConcurrentOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test concurrent SET operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_key_%d", i)
			err := rdb.Set(ctx, key, fmt.Sprintf("value_%d", i), 0).Err()
			if err != nil {
				t.Errorf("Error setting key %s: %v", key, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all keys were set
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_key_%d", i)
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			t.Errorf("Error getting key %s: %v", key, err)
		}
		expected := fmt.Sprintf("value_%d", i)
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	}

	// Cleanup
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("concurrent_key_%d", i)
	}
	rdb.Del(ctx, keys...)
}

// TestPerformance tests basic performance
func TestPerformance(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test SET performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_key_%d", i)
		err := rdb.Set(ctx, key, fmt.Sprintf("value_%d", i), 0).Err()
		if err != nil {
			t.Fatalf("Error setting key %s: %v", key, err)
		}
	}
	duration := time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("SET operations took too long: %v", duration)
	}

	// Test GET performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_key_%d", i)
		_, err := rdb.Get(ctx, key).Result()
		if err != nil {
			t.Fatalf("Error getting key %s: %v", key, err)
		}
	}
	duration = time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("GET operations took too long: %v", duration)
	}

	// Cleanup
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = fmt.Sprintf("perf_key_%d", i)
	}
	rdb.Del(ctx, keys...)
}
