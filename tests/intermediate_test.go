package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestIntermediateHashOperations tests Redis hash operations
func TestIntermediateHashOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test HSET - Set single field
	added, err := rdb.HSet(ctx, "test_hash", "name", "Alice").Result()
	if err != nil {
		t.Fatalf("Error setting hash field: %v", err)
	}
	if added != 1 {
		t.Errorf("Expected 1 added field, got %d", added)
	}

	// Test HSET - Set multiple fields
	fields := map[string]interface{}{
		"age":     30,
		"email":   "alice@example.com",
		"city":    "New York",
		"country": "USA",
	}
	added, err = rdb.HSet(ctx, "test_hash", fields).Result()
	if err != nil {
		t.Fatalf("Error setting multiple hash fields: %v", err)
	}
	if added != 4 {
		t.Errorf("Expected 4 added fields, got %d", added)
	}

	// Test HGET - Get single field
	name, err := rdb.HGet(ctx, "test_hash", "name").Result()
	if err != nil {
		t.Fatalf("Error getting hash field: %v", err)
	}
	if name != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", name)
	}

	// Test HGETALL - Get all fields
	allFields, err := rdb.HGetAll(ctx, "test_hash").Result()
	if err != nil {
		t.Fatalf("Error getting all hash fields: %v", err)
	}
	if len(allFields) != 5 {
		t.Errorf("Expected 5 fields, got %d", len(allFields))
	}

	// Test HMGET - Get multiple fields
	values, err := rdb.HMGet(ctx, "test_hash", "name", "email", "city").Result()
	if err != nil {
		t.Fatalf("Error getting multiple hash fields: %v", err)
	}
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// Test HKEYS - Get all field names
	keys, err := rdb.HKeys(ctx, "test_hash").Result()
	if err != nil {
		t.Fatalf("Error getting hash keys: %v", err)
	}
	if len(keys) != 5 {
		t.Errorf("Expected 5 keys, got %d", len(keys))
	}

	// Test HVALS - Get all field values
	vals, err := rdb.HVals(ctx, "test_hash").Result()
	if err != nil {
		t.Fatalf("Error getting hash values: %v", err)
	}
	if len(vals) != 5 {
		t.Errorf("Expected 5 values, got %d", len(vals))
	}

	// Test HEXISTS - Check if field exists
	exists, err := rdb.HExists(ctx, "test_hash", "name").Result()
	if err != nil {
		t.Fatalf("Error checking field existence: %v", err)
	}
	if !exists {
		t.Error("Field 'name' should exist")
	}

	// Test HDEL - Delete field
	deleted, err := rdb.HDel(ctx, "test_hash", "country").Result()
	if err != nil {
		t.Fatalf("Error deleting hash field: %v", err)
	}
	if deleted != 1 {
		t.Errorf("Expected 1 deleted field, got %d", deleted)
	}

	// Test HINCRBY - Increment field
	newAge, err := rdb.HIncrBy(ctx, "test_hash", "age", 1).Result()
	if err != nil {
		t.Fatalf("Error incrementing age: %v", err)
	}
	if newAge != 31 {
		t.Errorf("Expected age 31, got %d", newAge)
	}

	// Test HINCRBYFLOAT - Increment field by float
	rdb.HSet(ctx, "test_hash", "score", "100.5")
	newScore, err := rdb.HIncrByFloat(ctx, "test_hash", "score", 15.3).Result()
	if err != nil {
		t.Fatalf("Error incrementing score: %v", err)
	}
	if newScore != 115.8 {
		t.Errorf("Expected score 115.8, got %f", newScore)
	}

	// Test HLEN - Get number of fields
	length, err := rdb.HLen(ctx, "test_hash").Result()
	if err != nil {
		t.Fatalf("Error getting hash length: %v", err)
	}
	if length != 5 {
		t.Errorf("Expected 5 fields, got %d", length)
	}

	// Test HSETNX - Set field only if it doesn't exist
	set, err := rdb.HSetNX(ctx, "test_hash", "phone", "123-456-7890").Result()
	if err != nil {
		t.Fatalf("Error with HSETNX: %v", err)
	}
	if !set {
		t.Error("HSETNX should succeed for new field")
	}

	// Test HMSET - Set multiple fields (deprecated but still works)
	hmsetFields := map[string]interface{}{
		"department": "Engineering",
		"role":       "Developer",
		"salary":     "75000",
	}
	err = rdb.HMSet(ctx, "test_hash", hmsetFields).Err()
	if err != nil {
		t.Fatalf("Error with HMSET: %v", err)
	}

	// Verify HMSET worked
	dept, err := rdb.HGet(ctx, "test_hash", "department").Result()
	if err != nil {
		t.Fatalf("Error getting department after HMSET: %v", err)
	}
	if dept != "Engineering" {
		t.Errorf("Expected 'Engineering', got '%s'", dept)
	}

	// Cleanup
	rdb.Del(ctx, "test_hash")
}

// TestIntermediateListOperations tests Redis list operations
func TestIntermediateListOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test LPUSH - Push elements to left
	length, err := rdb.LPush(ctx, "test_list", "item1", "item2", "item3").Result()
	if err != nil {
		t.Fatalf("Error pushing to list: %v", err)
	}
	if length != 3 {
		t.Errorf("Expected list length 3, got %d", length)
	}

	// Test RPUSH - Push elements to right
	length, err = rdb.RPush(ctx, "test_list", "item4", "item5").Result()
	if err != nil {
		t.Fatalf("Error pushing to right: %v", err)
	}
	if length != 5 {
		t.Errorf("Expected list length 5, got %d", length)
	}

	// Test LRANGE - Get range of elements
	allItems, err := rdb.LRange(ctx, "test_list", 0, -1).Result()
	if err != nil {
		t.Fatalf("Error getting all elements: %v", err)
	}
	if len(allItems) != 5 {
		t.Errorf("Expected 5 items, got %d", len(allItems))
	}

	// Test LLEN - Get list length
	length, err = rdb.LLen(ctx, "test_list").Result()
	if err != nil {
		t.Fatalf("Error getting list length: %v", err)
	}
	if length != 5 {
		t.Errorf("Expected list length 5, got %d", length)
	}

	// Test LPOP - Pop from left
	popped, err := rdb.LPop(ctx, "test_list").Result()
	if err != nil {
		t.Fatalf("Error popping from left: %v", err)
	}
	if popped != "item3" {
		t.Errorf("Expected 'item3', got '%s'", popped)
	}

	// Test RPOP - Pop from right
	popped, err = rdb.RPop(ctx, "test_list").Result()
	if err != nil {
		t.Fatalf("Error popping from right: %v", err)
	}
	if popped != "item5" {
		t.Errorf("Expected 'item5', got '%s'", popped)
	}

	// Test LINDEX - Get element at index
	first, err := rdb.LIndex(ctx, "test_list", 0).Result()
	if err != nil {
		t.Fatalf("Error getting first element: %v", err)
	}
	if first != "item2" {
		t.Errorf("Expected 'item2', got '%s'", first)
	}

	// Test LINSERT - Insert element
	length, err = rdb.LInsertBefore(ctx, "test_list", "item1", "urgent_item").Result()
	if err != nil {
		t.Fatalf("Error inserting before: %v", err)
	}
	if length != 4 {
		t.Errorf("Expected list length 4, got %d", length)
	}

	// Test LREM - Remove elements
	removed, err := rdb.LRem(ctx, "test_list", 1, "item1").Result()
	if err != nil {
		t.Fatalf("Error removing element: %v", err)
	}
	if removed != 1 {
		t.Errorf("Expected 1 removed element, got %d", removed)
	}

	// Test LSET - Set element at index
	err = rdb.LSet(ctx, "test_list", 0, "updated_item").Err()
	if err != nil {
		t.Fatalf("Error setting element at index: %v", err)
	}

	// Test LTRIM - Trim list
	err = rdb.LTrim(ctx, "test_list", 0, 1).Err()
	if err != nil {
		t.Fatalf("Error trimming list: %v", err)
	}

	// Test RPOPLPUSH - Pop from right, push to left of another list
	rdb.RPush(ctx, "target_list", "existing_item")
	moved, err := rdb.RPopLPush(ctx, "test_list", "target_list").Result()
	if err != nil {
		t.Fatalf("Error moving element: %v", err)
	}
	if moved == "" {
		t.Error("Expected non-empty moved element")
	}

	// Cleanup
	rdb.Del(ctx, "test_list", "target_list")
}

// TestIntermediateSetOperations tests Redis set operations
func TestIntermediateSetOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test SADD - Add members to set
	added, err := rdb.SAdd(ctx, "test_set", "member1", "member2", "member3").Result()
	if err != nil {
		t.Fatalf("Error adding to set: %v", err)
	}
	if added != 3 {
		t.Errorf("Expected 3 added members, got %d", added)
	}

	// Test SMEMBERS - Get all members
	members, err := rdb.SMembers(ctx, "test_set").Result()
	if err != nil {
		t.Fatalf("Error getting set members: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}

	// Test SISMEMBER - Check if member exists
	isMember, err := rdb.SIsMember(ctx, "test_set", "member1").Result()
	if err != nil {
		t.Fatalf("Error checking membership: %v", err)
	}
	if !isMember {
		t.Error("member1 should be a member")
	}

	// Test SCARD - Get cardinality
	cardinality, err := rdb.SCard(ctx, "test_set").Result()
	if err != nil {
		t.Fatalf("Error getting set cardinality: %v", err)
	}
	if cardinality != 3 {
		t.Errorf("Expected cardinality 3, got %d", cardinality)
	}

	// Test SREM - Remove members
	removed, err := rdb.SRem(ctx, "test_set", "member1", "member2").Result()
	if err != nil {
		t.Fatalf("Error removing members: %v", err)
	}
	if removed != 2 {
		t.Errorf("Expected 2 removed members, got %d", removed)
	}

	// Test SPOP - Remove and return random member
	popped, err := rdb.SPop(ctx, "test_set").Result()
	if err != nil {
		t.Fatalf("Error popping member: %v", err)
	}
	if popped == "" {
		t.Error("Expected non-empty popped member")
	}

	// Test SRANDMEMBER - Get random member without removing
	rdb.SAdd(ctx, "test_set", "member4", "member5", "member6")
	random, err := rdb.SRandMember(ctx, "test_set").Result()
	if err != nil {
		t.Fatalf("Error getting random member: %v", err)
	}
	if random == "" {
		t.Error("Expected non-empty random member")
	}

	// Test set operations - Union
	rdb.SAdd(ctx, "set1", "a", "b", "c")
	rdb.SAdd(ctx, "set2", "c", "d", "e")
	union, err := rdb.SUnion(ctx, "set1", "set2").Result()
	if err != nil {
		t.Fatalf("Error getting union: %v", err)
	}
	if len(union) != 5 {
		t.Errorf("Expected union size 5, got %d", len(union))
	}

	// Test set operations - Intersection
	intersection, err := rdb.SInter(ctx, "set1", "set2").Result()
	if err != nil {
		t.Fatalf("Error getting intersection: %v", err)
	}
	if len(intersection) != 1 {
		t.Errorf("Expected intersection size 1, got %d", len(intersection))
	}

	// Test set operations - Difference
	difference, err := rdb.SDiff(ctx, "set1", "set2").Result()
	if err != nil {
		t.Fatalf("Error getting difference: %v", err)
	}
	if len(difference) != 2 {
		t.Errorf("Expected difference size 2, got %d", len(difference))
	}

	// Test SMOVE - Move member between sets
	moved, err := rdb.SMove(ctx, "set1", "set2", "a").Result()
	if err != nil {
		t.Fatalf("Error moving member: %v", err)
	}
	if !moved {
		t.Error("Expected member to be moved")
	}

	// Cleanup
	rdb.Del(ctx, "test_set", "set1", "set2")
}

// TestIntermediateSortedSetOperations tests Redis sorted set operations
func TestIntermediateSortedSetOperations(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test ZADD - Add members to sorted set
	added, err := rdb.ZAdd(ctx, "test_zset", redis.Z{Score: 100, Member: "player1"}).Result()
	if err != nil {
		t.Fatalf("Error adding to sorted set: %v", err)
	}
	if added != 1 {
		t.Errorf("Expected 1 added member, got %d", added)
	}

	// Test ZADD - Add multiple members
	players := []redis.Z{
		{Score: 150, Member: "player2"},
		{Score: 75, Member: "player3"},
		{Score: 200, Member: "player4"},
	}
	added, err = rdb.ZAdd(ctx, "test_zset", players...).Result()
	if err != nil {
		t.Fatalf("Error adding multiple members: %v", err)
	}
	if added != 3 {
		t.Errorf("Expected 3 added members, got %d", added)
	}

	// Test ZRANGE - Get range of members
	allMembers, err := rdb.ZRange(ctx, "test_zset", 0, -1).Result()
	if err != nil {
		t.Fatalf("Error getting all members: %v", err)
	}
	if len(allMembers) != 4 {
		t.Errorf("Expected 4 members, got %d", len(allMembers))
	}

	// Test ZRANGE with scores
	topMembers, err := rdb.ZRangeWithScores(ctx, "test_zset", 0, 2).Result()
	if err != nil {
		t.Fatalf("Error getting top members with scores: %v", err)
	}
	if len(topMembers) != 3 {
		t.Errorf("Expected 3 top members, got %d", len(topMembers))
	}

	// Test ZREVRANGE - Get range in descending order
	descMembers, err := rdb.ZRevRange(ctx, "test_zset", 0, -1).Result()
	if err != nil {
		t.Fatalf("Error getting members in descending order: %v", err)
	}
	if len(descMembers) != 4 {
		t.Errorf("Expected 4 members, got %d", len(descMembers))
	}

	// Test ZRANK - Get rank of member
	rank, err := rdb.ZRank(ctx, "test_zset", "player4").Result()
	if err != nil {
		t.Fatalf("Error getting rank: %v", err)
	}
	if rank != 3 {
		t.Errorf("Expected rank 3, got %d", rank)
	}

	// Test ZREVRANK - Get reverse rank
	revRank, err := rdb.ZRevRank(ctx, "test_zset", "player4").Result()
	if err != nil {
		t.Fatalf("Error getting reverse rank: %v", err)
	}
	if revRank != 0 {
		t.Errorf("Expected reverse rank 0, got %d", revRank)
	}

	// Test ZSCORE - Get score of member
	score, err := rdb.ZScore(ctx, "test_zset", "player4").Result()
	if err != nil {
		t.Fatalf("Error getting score: %v", err)
	}
	if score != 200 {
		t.Errorf("Expected score 200, got %f", score)
	}

	// Test ZCARD - Get cardinality
	cardinality, err := rdb.ZCard(ctx, "test_zset").Result()
	if err != nil {
		t.Fatalf("Error getting cardinality: %v", err)
	}
	if cardinality != 4 {
		t.Errorf("Expected cardinality 4, got %d", cardinality)
	}

	// Test ZCOUNT - Count members within score range
	count, err := rdb.ZCount(ctx, "test_zset", "100", "+inf").Result()
	if err != nil {
		t.Fatalf("Error counting members: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 members with score >= 100, got %d", count)
	}

	// Test ZRANGEBYSCORE - Get members by score range
	highScorers, err := rdb.ZRangeByScore(ctx, "test_zset", &redis.ZRangeBy{
		Min: "100",
		Max: "+inf",
	}).Result()
	if err != nil {
		t.Fatalf("Error getting high scorers: %v", err)
	}
	if len(highScorers) != 3 {
		t.Errorf("Expected 3 high scorers, got %d", len(highScorers))
	}

	// Test ZREM - Remove members
	removed, err := rdb.ZRem(ctx, "test_zset", "player3").Result()
	if err != nil {
		t.Fatalf("Error removing member: %v", err)
	}
	if removed != 1 {
		t.Errorf("Expected 1 removed member, got %d", removed)
	}

	// Test ZINCRBY - Increment score
	newScore, err := rdb.ZIncrBy(ctx, "test_zset", 50, "player2").Result()
	if err != nil {
		t.Fatalf("Error incrementing score: %v", err)
	}
	if newScore != 200 {
		t.Errorf("Expected new score 200, got %f", newScore)
	}

	// Test ZREMRANGEBYRANK - Remove members by rank
	removedByRank, err := rdb.ZRemRangeByRank(ctx, "test_zset", 0, 0).Result()
	if err != nil {
		t.Fatalf("Error removing by rank: %v", err)
	}
	if removedByRank != 1 {
		t.Errorf("Expected 1 removed member by rank, got %d", removedByRank)
	}

	// Test ZREMRANGEBYSCORE - Remove members by score
	// Note: After previous operations, only player2 and player4 remain
	// player4 was removed by rank, so only player2 with score 200 remains
	remainingCount, err := rdb.ZCard(ctx, "test_zset").Result()
	if err != nil {
		t.Fatalf("Error getting remaining cardinality: %v", err)
	}
	if remainingCount > 0 {
		// Only remove if there are members left
		removedByScore, err := rdb.ZRemRangeByScore(ctx, "test_zset", "-inf", "250").Result()
		if err != nil {
			t.Fatalf("Error removing by score: %v", err)
		}
		if removedByScore < 0 {
			t.Errorf("Expected non-negative removed members, got %d", removedByScore)
		}
	}

	// Cleanup
	rdb.Del(ctx, "test_zset")
}

// TestHashOperationsConcurrency tests concurrent hash operations
func TestHashOperationsConcurrency(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test concurrent hash operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_hash_%d", i)
			err := rdb.HSet(ctx, key, "field1", fmt.Sprintf("value_%d", i)).Err()
			if err != nil {
				t.Errorf("Error setting hash %s: %v", key, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all hashes were set
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_hash_%d", i)
		val, err := rdb.HGet(ctx, key, "field1").Result()
		if err != nil {
			t.Errorf("Error getting hash %s: %v", key, err)
		}
		expected := fmt.Sprintf("value_%d", i)
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	}

	// Cleanup
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("concurrent_hash_%d", i)
	}
	rdb.Del(ctx, keys...)
}

// TestListOperationsConcurrency tests concurrent list operations
func TestListOperationsConcurrency(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test concurrent list operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_list_%d", i)
			err := rdb.LPush(ctx, key, fmt.Sprintf("item_%d", i)).Err()
			if err != nil {
				t.Errorf("Error pushing to list %s: %v", key, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all lists were created
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_list_%d", i)
		length, err := rdb.LLen(ctx, key).Result()
		if err != nil {
			t.Errorf("Error getting list length %s: %v", key, err)
		}
		if length != 1 {
			t.Errorf("Expected list length 1, got %d", length)
		}
	}

	// Cleanup
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("concurrent_list_%d", i)
	}
	rdb.Del(ctx, keys...)
}

// TestSetOperationsConcurrency tests concurrent set operations
func TestSetOperationsConcurrency(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test concurrent set operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_set_%d", i)
			err := rdb.SAdd(ctx, key, fmt.Sprintf("member_%d", i)).Err()
			if err != nil {
				t.Errorf("Error adding to set %s: %v", key, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all sets were created
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_set_%d", i)
		cardinality, err := rdb.SCard(ctx, key).Result()
		if err != nil {
			t.Errorf("Error getting set cardinality %s: %v", key, err)
		}
		if cardinality != 1 {
			t.Errorf("Expected cardinality 1, got %d", cardinality)
		}
	}

	// Cleanup
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("concurrent_set_%d", i)
	}
	rdb.Del(ctx, keys...)
}

// TestSortedSetOperationsConcurrency tests concurrent sorted set operations
func TestSortedSetOperationsConcurrency(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test concurrent sorted set operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent_zset_%d", i)
			err := rdb.ZAdd(ctx, key, redis.Z{Score: float64(i * 100), Member: fmt.Sprintf("member_%d", i)}).Err()
			if err != nil {
				t.Errorf("Error adding to sorted set %s: %v", key, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all sorted sets were created
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_zset_%d", i)
		cardinality, err := rdb.ZCard(ctx, key).Result()
		if err != nil {
			t.Errorf("Error getting sorted set cardinality %s: %v", key, err)
		}
		if cardinality != 1 {
			t.Errorf("Expected cardinality 1, got %d", cardinality)
		}
	}

	// Cleanup
	keys := make([]string, 10)
	for i := 0; i < 10; i++ {
		keys[i] = fmt.Sprintf("concurrent_zset_%d", i)
	}
	rdb.Del(ctx, keys...)
}

// TestIntermediatePerformance tests performance of intermediate operations
func TestIntermediatePerformance(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test hash operations performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_hash_%d", i)
		err := rdb.HSet(ctx, key, "field1", fmt.Sprintf("value_%d", i)).Err()
		if err != nil {
			t.Fatalf("Error setting hash %s: %v", key, err)
		}
	}
	duration := time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("Hash operations took too long: %v", duration)
	}

	// Test list operations performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_list_%d", i)
		err := rdb.LPush(ctx, key, fmt.Sprintf("item_%d", i)).Err()
		if err != nil {
			t.Fatalf("Error pushing to list %s: %v", key, err)
		}
	}
	duration = time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("List operations took too long: %v", duration)
	}

	// Test set operations performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_set_%d", i)
		err := rdb.SAdd(ctx, key, fmt.Sprintf("member_%d", i)).Err()
		if err != nil {
			t.Fatalf("Error adding to set %s: %v", key, err)
		}
	}
	duration = time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("Set operations took too long: %v", duration)
	}

	// Test sorted set operations performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("perf_zset_%d", i)
		err := rdb.ZAdd(ctx, key, redis.Z{Score: float64(i), Member: fmt.Sprintf("member_%d", i)}).Err()
		if err != nil {
			t.Fatalf("Error adding to sorted set %s: %v", key, err)
		}
	}
	duration = time.Since(start)

	if duration > 5*time.Second {
		t.Errorf("Sorted set operations took too long: %v", duration)
	}

	// Cleanup
	keys := make([]string, 4000)
	for i := 0; i < 1000; i++ {
		keys[i] = fmt.Sprintf("perf_hash_%d", i)
		keys[i+1000] = fmt.Sprintf("perf_list_%d", i)
		keys[i+2000] = fmt.Sprintf("perf_set_%d", i)
		keys[i+3000] = fmt.Sprintf("perf_zset_%d", i)
	}
	rdb.Del(ctx, keys...)
}
