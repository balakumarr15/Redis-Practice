package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis Sorted Set operations - ZADD, ZRANGE, ZRANK, etc.
func SortedSetOperations() {
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

	// 1. ZADD - Add members to sorted set
	fmt.Println("\n=== ZADD Command ===")

	// Add single member with score
	added, err := rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "player1"}).Result()
	if err != nil {
		log.Fatalf("Error adding to sorted set: %v", err)
	}
	fmt.Printf("ZADD single member result: %d\n", added)

	// Add multiple members
	players := []redis.Z{
		{Score: 150, Member: "player2"},
		{Score: 75, Member: "player3"},
		{Score: 200, Member: "player4"},
		{Score: 125, Member: "player5"},
	}
	added, err = rdb.ZAdd(ctx, "leaderboard", players...).Result()
	if err != nil {
		log.Fatalf("Error adding multiple members: %v", err)
	}
	fmt.Printf("ZADD multiple members result: %d\n", added)

	// 2. ZRANGE - Get range of members by rank
	fmt.Println("\n=== ZRANGE Command ===")

	// Get all members (ascending order)
	allMembers, err := rdb.ZRange(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting all members: %v", err)
	}
	fmt.Printf("All members (ascending): %v\n", allMembers)

	// Get top 3 members
	top3, err := rdb.ZRange(ctx, "leaderboard", 0, 2).Result()
	if err != nil {
		log.Fatalf("Error getting top 3: %v", err)
	}
	fmt.Printf("Top 3 members: %v\n", top3)

	// Get members with scores
	top3WithScores, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, 2).Result()
	if err != nil {
		log.Fatalf("Error getting top 3 with scores: %v", err)
	}
	fmt.Println("Top 3 with scores:")
	for _, z := range top3WithScores {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 3. ZREVRANGE - Get range in descending order
	fmt.Println("\n=== ZREVRANGE Command ===")

	// Get all members in descending order
	allMembersDesc, err := rdb.ZRevRange(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting all members desc: %v", err)
	}
	fmt.Printf("All members (descending): %v\n", allMembersDesc)

	// Get top 3 in descending order
	top3Desc, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
	if err != nil {
		log.Fatalf("Error getting top 3 desc: %v", err)
	}
	fmt.Println("Top 3 (descending) with scores:")
	for _, z := range top3Desc {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 4. ZRANK - Get rank of member
	fmt.Println("\n=== ZRANK Command ===")

	// Get rank of player4 (ascending)
	rank, err := rdb.ZRank(ctx, "leaderboard", "player4").Result()
	if err != nil {
		log.Fatalf("Error getting rank: %v", err)
	}
	fmt.Printf("Rank of player4 (ascending): %d\n", rank)

	// Get rank in descending order
	revRank, err := rdb.ZRevRank(ctx, "leaderboard", "player4").Result()
	if err != nil {
		log.Fatalf("Error getting reverse rank: %v", err)
	}
	fmt.Printf("Rank of player4 (descending): %d\n", revRank)

	// 5. ZSCORE - Get score of member
	fmt.Println("\n=== ZSCORE Command ===")

	score, err := rdb.ZScore(ctx, "leaderboard", "player4").Result()
	if err != nil {
		log.Fatalf("Error getting score: %v", err)
	}
	fmt.Printf("Score of player4: %.0f\n", score)

	// 6. ZCARD - Get cardinality of sorted set
	fmt.Println("\n=== ZCARD Command ===")

	cardinality, err := rdb.ZCard(ctx, "leaderboard").Result()
	if err != nil {
		log.Fatalf("Error getting cardinality: %v", err)
	}
	fmt.Printf("Number of players: %d\n", cardinality)

	// 7. ZCOUNT - Count members within score range
	fmt.Println("\n=== ZCOUNT Command ===")

	// Count players with score >= 100
	count, err := rdb.ZCount(ctx, "leaderboard", "100", "+inf").Result()
	if err != nil {
		log.Fatalf("Error counting members: %v", err)
	}
	fmt.Printf("Players with score >= 100: %d\n", count)

	// Count players with score between 100 and 150
	count, err = rdb.ZCount(ctx, "leaderboard", "100", "150").Result()
	if err != nil {
		log.Fatalf("Error counting members in range: %v", err)
	}
	fmt.Printf("Players with score 100-150: %d\n", count)

	// 8. ZRANGEBYSCORE - Get members by score range
	fmt.Println("\n=== ZRANGEBYSCORE Command ===")

	// Get players with score >= 100
	highScorers, err := rdb.ZRangeByScore(ctx, "leaderboard", &redis.ZRangeBy{
		Min: "100",
		Max: "+inf",
	}).Result()
	if err != nil {
		log.Fatalf("Error getting high scorers: %v", err)
	}
	fmt.Printf("High scorers (>=100): %v\n", highScorers)

	// Get players with score between 100 and 150, with scores
	midScorers, err := rdb.ZRangeByScoreWithScores(ctx, "leaderboard", &redis.ZRangeBy{
		Min: "100",
		Max: "150",
	}).Result()
	if err != nil {
		log.Fatalf("Error getting mid scorers: %v", err)
	}
	fmt.Println("Mid scorers (100-150) with scores:")
	for _, z := range midScorers {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 9. ZREM - Remove members from sorted set
	fmt.Println("\n=== ZREM Command ===")

	// Remove single member
	removed, err := rdb.ZRem(ctx, "leaderboard", "player3").Result()
	if err != nil {
		log.Fatalf("Error removing member: %v", err)
	}
	fmt.Printf("ZREM single member result: %d\n", removed)

	// Remove multiple members
	removed, err = rdb.ZRem(ctx, "leaderboard", "player1", "player5").Result()
	if err != nil {
		log.Fatalf("Error removing multiple members: %v", err)
	}
	fmt.Printf("ZREM multiple members result: %d\n", removed)

	// Show updated leaderboard
	updated, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting updated leaderboard: %v", err)
	}
	fmt.Println("Updated leaderboard:")
	for _, z := range updated {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 10. ZINCRBY - Increment score of member
	fmt.Println("\n=== ZINCRBY Command ===")

	// Increment player2's score by 50
	newScore, err := rdb.ZIncrBy(ctx, "leaderboard", 50, "player2").Result()
	if err != nil {
		log.Fatalf("Error incrementing score: %v", err)
	}
	fmt.Printf("Player2's new score: %.0f\n", newScore)

	// Show updated leaderboard
	final, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting final leaderboard: %v", err)
	}
	fmt.Println("Final leaderboard:")
	for _, z := range final {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 11. ZREMRANGEBYRANK - Remove members by rank range
	fmt.Println("\n=== ZREMRANGEBYRANK Command ===")

	// Add more players for demo
	morePlayers := []redis.Z{
		{Score: 50, Member: "player6"},
		{Score: 300, Member: "player7"},
		{Score: 25, Member: "player8"},
	}
	_, err = rdb.ZAdd(ctx, "leaderboard", morePlayers...).Result()
	if err != nil {
		log.Fatalf("Error adding more players: %v", err)
	}

	// Remove bottom 2 players (ranks 0 and 1)
	removedByRank, err := rdb.ZRemRangeByRank(ctx, "leaderboard", 0, 1).Result()
	if err != nil {
		log.Fatalf("Error removing by rank: %v", err)
	}
	fmt.Printf("Removed %d players by rank\n", removedByRank)

	// 12. ZREMRANGEBYSCORE - Remove members by score range
	fmt.Println("\n=== ZREMRANGEBYSCORE Command ===")

	// Remove players with score < 100
	removedByScore, err := rdb.ZRemRangeByScore(ctx, "leaderboard", "-inf", "(100").Result()
	if err != nil {
		log.Fatalf("Error removing by score: %v", err)
	}
	fmt.Printf("Removed %d players by score\n", removedByScore)

	// Show final leaderboard
	finalLeaderboard, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting final leaderboard: %v", err)
	}
	fmt.Println("Final leaderboard after cleanup:")
	for _, z := range finalLeaderboard {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 13. Practical example - Product ratings
	fmt.Println("\n=== Practical Example: Product Ratings ===")

	// Add product ratings
	ratings := []redis.Z{
		{Score: 4.5, Member: "product:laptop"},
		{Score: 3.8, Member: "product:phone"},
		{Score: 4.2, Member: "product:tablet"},
		{Score: 4.7, Member: "product:headphones"},
		{Score: 3.5, Member: "product:mouse"},
	}
	_, err = rdb.ZAdd(ctx, "product_ratings", ratings...).Result()
	if err != nil {
		log.Fatalf("Error adding product ratings: %v", err)
	}

	// Get top rated products
	topRated, err := rdb.ZRevRangeWithScores(ctx, "product_ratings", 0, 2).Result()
	if err != nil {
		log.Fatalf("Error getting top rated products: %v", err)
	}
	fmt.Println("Top rated products:")
	for _, z := range topRated {
		fmt.Printf("  %s: %.1f stars\n", z.Member, z.Score)
	}

	// Get products with rating >= 4.0
	highRated, err := rdb.ZRangeByScoreWithScores(ctx, "product_ratings", &redis.ZRangeBy{
		Min: "4.0",
		Max: "+inf",
	}).Result()
	if err != nil {
		log.Fatalf("Error getting high rated products: %v", err)
	}
	fmt.Println("Products with rating >= 4.0:")
	for _, z := range highRated {
		fmt.Printf("  %s: %.1f stars\n", z.Member, z.Score)
	}
}
