package main

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestRedisConnection tests basic Redis connection
func TestRedisConnection(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Could not connect to Redis: %v", err)
	}

	if pong != "PONG" {
		t.Errorf("Expected PONG, got %s", pong)
	}
}

// TestBasicSetGet tests basic SET and GET operations
func TestBasicSetGet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Test SET
	err := rdb.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		t.Fatalf("Error setting key: %v", err)
	}

	// Test GET
	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}

	if val != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", val)
	}

	// Cleanup
	rdb.Del(ctx, "test_key")
}

// TestExpiration tests key expiration
func TestExpiration(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set key with 1 second expiration
	err := rdb.Set(ctx, "expire_key", "expire_value", 1*time.Second).Err()
	if err != nil {
		t.Fatalf("Error setting key with expiration: %v", err)
	}

	// Check key exists
	exists, err := rdb.Exists(ctx, "expire_key").Result()
	if err != nil {
		t.Fatalf("Error checking key existence: %v", err)
	}

	if exists == 0 {
		t.Error("Key should exist immediately after setting")
	}

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Check key no longer exists
	exists, err = rdb.Exists(ctx, "expire_key").Result()
	if err != nil {
		t.Fatalf("Error checking key existence after expiration: %v", err)
	}

	if exists > 0 {
		t.Error("Key should not exist after expiration")
	}
}

// TestDeleteKey tests key deletion
func TestDeleteKey(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set key
	err := rdb.Set(ctx, "delete_key", "delete_value", 0).Err()
	if err != nil {
		t.Fatalf("Error setting key: %v", err)
	}

	// Delete key
	deleted, err := rdb.Del(ctx, "delete_key").Result()
	if err != nil {
		t.Fatalf("Error deleting key: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected 1 deleted key, got %d", deleted)
	}

	// Verify key is deleted
	exists, err := rdb.Exists(ctx, "delete_key").Result()
	if err != nil {
		t.Fatalf("Error checking key existence: %v", err)
	}

	if exists > 0 {
		t.Error("Key should not exist after deletion")
	}
}

// TestTTL tests TTL operations
func TestTTL(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set key with expiration
	err := rdb.Set(ctx, "ttl_key", "ttl_value", 10*time.Second).Err()
	if err != nil {
		t.Fatalf("Error setting key: %v", err)
	}

	// Check TTL
	ttl, err := rdb.TTL(ctx, "ttl_key").Result()
	if err != nil {
		t.Fatalf("Error getting TTL: %v", err)
	}

	if ttl <= 0 || ttl > 10*time.Second {
		t.Errorf("Expected TTL between 0 and 10 seconds, got %v", ttl)
	}

	// Test EXPIRE command
	expired, err := rdb.Expire(ctx, "ttl_key", 5*time.Second).Result()
	if err != nil {
		t.Fatalf("Error setting expiration: %v", err)
	}

	if !expired {
		t.Error("EXPIRE command should return true")
	}

	// Check new TTL
	newTTL, err := rdb.TTL(ctx, "ttl_key").Result()
	if err != nil {
		t.Fatalf("Error getting new TTL: %v", err)
	}

	if newTTL <= 0 || newTTL > 5*time.Second {
		t.Errorf("Expected new TTL between 0 and 5 seconds, got %v", newTTL)
	}

	// Cleanup
	rdb.Del(ctx, "ttl_key")
}

// TestRenameKey tests key renaming
func TestRenameKey(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set original key
	err := rdb.Set(ctx, "old_key", "old_value", 0).Err()
	if err != nil {
		t.Fatalf("Error setting original key: %v", err)
	}

	// Rename key
	err = rdb.Rename(ctx, "old_key", "new_key").Err()
	if err != nil {
		t.Fatalf("Error renaming key: %v", err)
	}

	// Check old key doesn't exist
	exists, err := rdb.Exists(ctx, "old_key").Result()
	if err != nil {
		t.Fatalf("Error checking old key existence: %v", err)
	}

	if exists > 0 {
		t.Error("Old key should not exist after rename")
	}

	// Check new key exists with correct value
	val, err := rdb.Get(ctx, "new_key").Result()
	if err != nil {
		t.Fatalf("Error getting renamed key: %v", err)
	}

	if val != "old_value" {
		t.Errorf("Expected 'old_value', got '%s'", val)
	}

	// Cleanup
	rdb.Del(ctx, "new_key")
}

// TestMultipleKeys tests operations with multiple keys
func TestMultipleKeys(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set multiple keys
	keys := []string{"key1", "key2", "key3"}
	values := []string{"value1", "value2", "value3"}

	for i, key := range keys {
		err := rdb.Set(ctx, key, values[i], 0).Err()
		if err != nil {
			t.Fatalf("Error setting key %s: %v", key, err)
		}
	}

	// Check all keys exist
	exists, err := rdb.Exists(ctx, keys...).Result()
	if err != nil {
		t.Fatalf("Error checking multiple keys existence: %v", err)
	}

	if exists != int64(len(keys)) {
		t.Errorf("Expected %d keys to exist, got %d", len(keys), exists)
	}

	// Delete multiple keys
	deleted, err := rdb.Del(ctx, keys...).Result()
	if err != nil {
		t.Fatalf("Error deleting multiple keys: %v", err)
	}

	if deleted != int64(len(keys)) {
		t.Errorf("Expected %d keys to be deleted, got %d", len(keys), deleted)
	}

	// Verify all keys are deleted
	exists, err = rdb.Exists(ctx, keys...).Result()
	if err != nil {
		t.Fatalf("Error checking keys existence after deletion: %v", err)
	}

	if exists > 0 {
		t.Error("All keys should be deleted")
	}
}

// TestNonExistentKey tests operations on non-existent keys
func TestNonExistentKey(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Try to get non-existent key
	_, err := rdb.Get(ctx, "nonexistent_key").Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil error, got %v", err)
	}

	// Check if non-existent key exists
	exists, err := rdb.Exists(ctx, "nonexistent_key").Result()
	if err != nil {
		t.Fatalf("Error checking non-existent key: %v", err)
	}

	if exists > 0 {
		t.Error("Non-existent key should not exist")
	}

	// Try to delete non-existent key
	deleted, err := rdb.Del(ctx, "nonexistent_key").Result()
	if err != nil {
		t.Fatalf("Error deleting non-existent key: %v", err)
	}

	if deleted != 0 {
		t.Errorf("Expected 0 deleted keys, got %d", deleted)
	}
}

// TestPersistKey tests PERSIST operation
func TestPersistKey(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set key with expiration
	err := rdb.Set(ctx, "persist_key", "persist_value", 10*time.Second).Err()
	if err != nil {
		t.Fatalf("Error setting key with expiration: %v", err)
	}

	// Check TTL
	ttl, err := rdb.TTL(ctx, "persist_key").Result()
	if err != nil {
		t.Fatalf("Error getting TTL: %v", err)
	}

	if ttl <= 0 {
		t.Error("Key should have TTL")
	}

	// Persist key (remove expiration)
	persisted, err := rdb.Persist(ctx, "persist_key").Result()
	if err != nil {
		t.Fatalf("Error persisting key: %v", err)
	}

	if !persisted {
		t.Error("PERSIST should return true")
	}

	// Check TTL after persist
	ttl, err = rdb.TTL(ctx, "persist_key").Result()
	if err != nil {
		t.Fatalf("Error getting TTL after persist: %v", err)
	}

	if ttl != -1 {
		t.Errorf("Expected TTL -1 (no expiration), got %v", ttl)
	}

	// Cleanup
	rdb.Del(ctx, "persist_key")
}

// TestRenameNX tests RENAMENX operation
func TestRenameNX(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	// Set source key
	err := rdb.Set(ctx, "source_key", "source_value", 0).Err()
	if err != nil {
		t.Fatalf("Error setting source key: %v", err)
	}

	// Set destination key
	err = rdb.Set(ctx, "dest_key", "dest_value", 0).Err()
	if err != nil {
		t.Fatalf("Error setting destination key: %v", err)
	}

	// Try to rename to existing key (should fail)
	renamed, err := rdb.RenameNX(ctx, "source_key", "dest_key").Result()
	if err != nil {
		t.Fatalf("Error with RENAMENX: %v", err)
	}

	if renamed {
		t.Error("RENAMENX should fail when destination exists")
	}

	// Try to rename to non-existing key (should succeed)
	renamed, err = rdb.RenameNX(ctx, "source_key", "new_key").Result()
	if err != nil {
		t.Fatalf("Error with RENAMENX: %v", err)
	}

	if !renamed {
		t.Error("RENAMENX should succeed when destination doesn't exist")
	}

	// Verify source key is gone
	exists, err := rdb.Exists(ctx, "source_key").Result()
	if err != nil {
		t.Fatalf("Error checking source key: %v", err)
	}

	if exists > 0 {
		t.Error("Source key should not exist after rename")
	}

	// Verify new key exists
	val, err := rdb.Get(ctx, "new_key").Result()
	if err != nil {
		t.Fatalf("Error getting new key: %v", err)
	}

	if val != "source_value" {
		t.Errorf("Expected 'source_value', got '%s'", val)
	}

	// Cleanup
	rdb.Del(ctx, "dest_key", "new_key")
}
