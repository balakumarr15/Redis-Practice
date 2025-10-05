package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis Set operations - SADD, SREM, SISMEMBER, etc.
func SetOperations() {
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

	// 1. SADD - Add members to set
	fmt.Println("\n=== SADD Command ===")

	// Add single member
	added, err := rdb.SAdd(ctx, "fruits", "apple").Result()
	if err != nil {
		log.Fatalf("Error adding to set: %v", err)
	}
	fmt.Printf("SADD single member result: %d\n", added)

	// Add multiple members
	added, err = rdb.SAdd(ctx, "fruits", "banana", "orange", "grape", "apple").Result()
	if err != nil {
		log.Fatalf("Error adding multiple members: %v", err)
	}
	fmt.Printf("SADD multiple members result: %d (duplicates not added)\n", added)

	// 2. SMEMBERS - Get all members of set
	fmt.Println("\n=== SMEMBERS Command ===")

	members, err := rdb.SMembers(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error getting set members: %v", err)
	}
	fmt.Printf("Fruits set members: %v\n", members)

	// 3. SISMEMBER - Check if member exists in set
	fmt.Println("\n=== SISMEMBER Command ===")

	// Check for existing member
	exists, err := rdb.SIsMember(ctx, "fruits", "apple").Result()
	if err != nil {
		log.Fatalf("Error checking membership: %v", err)
	}
	fmt.Printf("'apple' is member of fruits: %v\n", exists)

	// Check for non-existing member
	exists, err = rdb.SIsMember(ctx, "fruits", "mango").Result()
	if err != nil {
		log.Fatalf("Error checking membership: %v", err)
	}
	fmt.Printf("'mango' is member of fruits: %v\n", exists)

	// 4. SCARD - Get cardinality (number of members) of set
	fmt.Println("\n=== SCARD Command ===")

	cardinality, err := rdb.SCard(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error getting set cardinality: %v", err)
	}
	fmt.Printf("Number of fruits: %d\n", cardinality)

	// 5. SREM - Remove members from set
	fmt.Println("\n=== SREM Command ===")

	// Remove single member
	removed, err := rdb.SRem(ctx, "fruits", "grape").Result()
	if err != nil {
		log.Fatalf("Error removing member: %v", err)
	}
	fmt.Printf("SREM single member result: %d\n", removed)

	// Remove multiple members
	removed, err = rdb.SRem(ctx, "fruits", "banana", "nonexistent").Result()
	if err != nil {
		log.Fatalf("Error removing multiple members: %v", err)
	}
	fmt.Printf("SREM multiple members result: %d\n", removed)

	// Show updated set
	updatedMembers, err := rdb.SMembers(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error getting updated members: %v", err)
	}
	fmt.Printf("Updated fruits set: %v\n", updatedMembers)

	// 6. SPOP - Remove and return random member
	fmt.Println("\n=== SPOP Command ===")

	// Add more fruits for SPOP demo
	_, err = rdb.SAdd(ctx, "fruits", "strawberry", "blueberry", "raspberry").Result()
	if err != nil {
		log.Fatalf("Error adding fruits for SPOP demo: %v", err)
	}

	// Pop random member
	popped, err := rdb.SPop(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error popping member: %v", err)
	}
	fmt.Printf("Popped random member: %s\n", popped)

	// Pop multiple random members
	poppedMultiple, err := rdb.SPopN(ctx, "fruits", 2).Result()
	if err != nil {
		log.Fatalf("Error popping multiple members: %v", err)
	}
	fmt.Printf("Popped multiple members: %v\n", poppedMultiple)

	// Show remaining set
	remaining, err := rdb.SMembers(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error getting remaining members: %v", err)
	}
	fmt.Printf("Remaining fruits: %v\n", remaining)

	// 7. SRANDMEMBER - Get random member without removing
	fmt.Println("\n=== SRANDMEMBER Command ===")

	// Add more fruits for SRANDMEMBER demo
	_, err = rdb.SAdd(ctx, "fruits", "kiwi", "pineapple", "mango", "peach").Result()
	if err != nil {
		log.Fatalf("Error adding fruits for SRANDMEMBER demo: %v", err)
	}

	// Get single random member
	random, err := rdb.SRandMember(ctx, "fruits").Result()
	if err != nil {
		log.Fatalf("Error getting random member: %v", err)
	}
	fmt.Printf("Random member: %s\n", random)

	// Get multiple random members
	randomMultiple, err := rdb.SRandMemberN(ctx, "fruits", 3).Result()
	if err != nil {
		log.Fatalf("Error getting multiple random members: %v", err)
	}
	fmt.Printf("Random members: %v\n", randomMultiple)

	// 8. Set operations - Union, Intersection, Difference
	fmt.Println("\n=== Set Operations ===")

	// Create two sets for operations
	_, err = rdb.SAdd(ctx, "set1", "a", "b", "c", "d").Result()
	if err != nil {
		log.Fatalf("Error creating set1: %v", err)
	}
	_, err = rdb.SAdd(ctx, "set2", "c", "d", "e", "f").Result()
	if err != nil {
		log.Fatalf("Error creating set2: %v", err)
	}

	// SUNION - Union of sets
	union, err := rdb.SUnion(ctx, "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error getting union: %v", err)
	}
	fmt.Printf("Union of set1 and set2: %v\n", union)

	// SINTER - Intersection of sets
	intersection, err := rdb.SInter(ctx, "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error getting intersection: %v", err)
	}
	fmt.Printf("Intersection of set1 and set2: %v\n", intersection)

	// SDIFF - Difference of sets (elements in set1 but not in set2)
	difference, err := rdb.SDiff(ctx, "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error getting difference: %v", err)
	}
	fmt.Printf("Difference (set1 - set2): %v\n", difference)

	// SDIFF - Difference in reverse (elements in set2 but not in set1)
	difference2, err := rdb.SDiff(ctx, "set2", "set1").Result()
	if err != nil {
		log.Fatalf("Error getting reverse difference: %v", err)
	}
	fmt.Printf("Difference (set2 - set1): %v\n", difference2)

	// 9. Store operations - SUNIONSTORE, SINTERSTORE, SDIFFSTORE
	fmt.Println("\n=== Store Operations ===")

	// SUNIONSTORE - Store union in new set
	unionCount, err := rdb.SUnionStore(ctx, "union_result", "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error storing union: %v", err)
	}
	fmt.Printf("Stored union with %d members\n", unionCount)

	// SINTERSTORE - Store intersection in new set
	intersectionCount, err := rdb.SInterStore(ctx, "intersection_result", "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error storing intersection: %v", err)
	}
	fmt.Printf("Stored intersection with %d members\n", intersectionCount)

	// SDIFFSTORE - Store difference in new set
	differenceCount, err := rdb.SDiffStore(ctx, "difference_result", "set1", "set2").Result()
	if err != nil {
		log.Fatalf("Error storing difference: %v", err)
	}
	fmt.Printf("Stored difference with %d members\n", differenceCount)

	// Show stored results
	unionResult, err := rdb.SMembers(ctx, "union_result").Result()
	if err != nil {
		log.Fatalf("Error getting union result: %v", err)
	}
	fmt.Printf("Union result: %v\n", unionResult)

	intersectionResult, err := rdb.SMembers(ctx, "intersection_result").Result()
	if err != nil {
		log.Fatalf("Error getting intersection result: %v", err)
	}
	fmt.Printf("Intersection result: %v\n", intersectionResult)

	differenceResult, err := rdb.SMembers(ctx, "difference_result").Result()
	if err != nil {
		log.Fatalf("Error getting difference result: %v", err)
	}
	fmt.Printf("Difference result: %v\n", differenceResult)

	// 10. SMOVE - Move member from one set to another
	fmt.Println("\n=== SMOVE Command ===")

	// Move member from set1 to set2
	moved, err := rdb.SMove(ctx, "set1", "set2", "a").Result()
	if err != nil {
		log.Fatalf("Error moving member: %v", err)
	}
	fmt.Printf("SMOVE result: %v\n", moved)

	// Show updated sets
	set1Members, err := rdb.SMembers(ctx, "set1").Result()
	if err != nil {
		log.Fatalf("Error getting set1 members: %v", err)
	}
	set2Members, err := rdb.SMembers(ctx, "set2").Result()
	if err != nil {
		log.Fatalf("Error getting set2 members: %v", err)
	}
	fmt.Printf("Set1 after move: %v\n", set1Members)
	fmt.Printf("Set2 after move: %v\n", set2Members)

	// 11. Practical example - User tags and interests
	fmt.Println("\n=== Practical Example: User Tags ===")

	// User interests
	_, err = rdb.SAdd(ctx, "user:alice:interests", "programming", "music", "travel", "cooking").Result()
	if err != nil {
		log.Fatalf("Error setting user interests: %v", err)
	}
	_, err = rdb.SAdd(ctx, "user:bob:interests", "programming", "sports", "gaming", "travel").Result()
	if err != nil {
		log.Fatalf("Error setting user interests: %v", err)
	}
	_, err = rdb.SAdd(ctx, "user:charlie:interests", "music", "art", "cooking", "photography").Result()
	if err != nil {
		log.Fatalf("Error setting user interests: %v", err)
	}

	// Find common interests between Alice and Bob
	commonInterests, err := rdb.SInter(ctx, "user:alice:interests", "user:bob:interests").Result()
	if err != nil {
		log.Fatalf("Error finding common interests: %v", err)
	}
	fmt.Printf("Common interests (Alice & Bob): %v\n", commonInterests)

	// Find all unique interests across all users
	allInterests, err := rdb.SUnion(ctx, "user:alice:interests", "user:bob:interests", "user:charlie:interests").Result()
	if err != nil {
		log.Fatalf("Error finding all interests: %v", err)
	}
	fmt.Printf("All unique interests: %v\n", allInterests)

	// Find interests unique to Alice
	aliceUnique, err := rdb.SDiff(ctx, "user:alice:interests", "user:bob:interests", "user:charlie:interests").Result()
	if err != nil {
		log.Fatalf("Error finding Alice's unique interests: %v", err)
	}
	fmt.Printf("Alice's unique interests: %v\n", aliceUnique)

	// 12. Clean up
	fmt.Println("\n=== Cleanup ===")
	keysToDelete := []string{"fruits", "set1", "set2", "union_result", "intersection_result", "difference_result", "user:alice:interests", "user:bob:interests", "user:charlie:interests"}
	deleted, err := rdb.Del(ctx, keysToDelete...).Result()
	if err != nil {
		log.Fatalf("Error cleaning up: %v", err)
	}
	fmt.Printf("Cleaned up %d keys\n", deleted)
}
