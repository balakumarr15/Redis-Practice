package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis Pub/Sub operations
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

	// 1. Basic PUBLISH and SUBSCRIBE
	fmt.Println("\n=== Basic Pub/Sub ===")

	// Create a subscriber
	subscriber := rdb.Subscribe(ctx, "channel1")
	defer subscriber.Close()

	// Start a goroutine to handle messages
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			fmt.Printf("Received message on %s: %s\n", msg.Channel, msg.Payload)
		}
	}()

	// Wait a moment for subscriber to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish messages
	for i := 1; i <= 3; i++ {
		err = rdb.Publish(ctx, "channel1", fmt.Sprintf("Hello %d", i)).Err()
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
		fmt.Printf("Published: Hello %d\n", i)
		time.Sleep(100 * time.Millisecond)
	}

	// 2. Multiple channels subscription
	fmt.Println("\n=== Multiple Channels ===")

	// Subscribe to multiple channels
	multiSub := rdb.Subscribe(ctx, "channel2", "channel3", "channel4")
	defer multiSub.Close()

	// Handle messages from multiple channels
	go func() {
		for {
			msg, err := multiSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			fmt.Printf("Multi-channel message on %s: %s\n", msg.Channel, msg.Payload)
		}
	}()

	// Wait for subscriber to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish to different channels
	channels := []string{"channel2", "channel3", "channel4"}
	for i, channel := range channels {
		err = rdb.Publish(ctx, channel, fmt.Sprintf("Message %d from %s", i+1, channel)).Err()
		if err != nil {
			log.Fatalf("Error publishing to %s: %v", channel, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 3. Pattern-based subscription
	fmt.Println("\n=== Pattern Subscription ===")

	// Subscribe to pattern
	patternSub := rdb.PSubscribe(ctx, "news:*")
	defer patternSub.Close()

	// Handle pattern messages
	go func() {
		for {
			msg, err := patternSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving pattern message: %v", err)
				return
			}
			fmt.Printf("Pattern message on %s: %s\n", msg.Channel, msg.Payload)
		}
	}()

	// Wait for subscriber to be ready
	time.Sleep(100 * time.Millisecond)

	// Publish to channels matching pattern
	newsChannels := []string{"news:sports", "news:tech", "news:politics"}
	for i, channel := range newsChannels {
		err = rdb.Publish(ctx, channel, fmt.Sprintf("Breaking news %d from %s", i+1, channel)).Err()
		if err != nil {
			log.Fatalf("Error publishing to %s: %v", channel, err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 4. Unsubscribe from channels
	fmt.Println("\n=== Unsubscribe ===")

	// Unsubscribe from specific channel
	err = multiSub.Unsubscribe(ctx, "channel2")
	if err != nil {
		log.Fatalf("Error unsubscribing from channel2: %v", err)
	}
	fmt.Println("Unsubscribed from channel2")

	// Publish to unsubscribed channel (should not be received)
	err = rdb.Publish(ctx, "channel2", "This should not be received").Err()
	if err != nil {
		log.Fatalf("Error publishing to unsubscribed channel: %v", err)
	}

	// Publish to still subscribed channels
	err = rdb.Publish(ctx, "channel3", "This should be received").Err()
	if err != nil {
		log.Fatalf("Error publishing to subscribed channel: %v", err)
	}

	// 5. Channel information
	fmt.Println("\n=== Channel Information ===")

	// Get number of subscribers for a channel
	subCount, err := rdb.PubSubNumSub(ctx, "channel1").Result()
	if err != nil {
		log.Fatalf("Error getting subscriber count: %v", err)
	}
	fmt.Printf("Subscribers for channel1: %v\n", subCount)

	// Get channels with subscribers
	channelsWithSubs, err := rdb.PubSubChannels(ctx, "*").Result()
	if err != nil {
		log.Fatalf("Error getting channels with subscribers: %v", err)
	}
	fmt.Printf("Channels with subscribers: %v\n", channelsWithSubs)

	// Get pattern subscriptions
	patterns, err := rdb.PubSubNumPat(ctx).Result()
	if err != nil {
		log.Fatalf("Error getting pattern count: %v", err)
	}
	fmt.Printf("Number of pattern subscriptions: %d\n", patterns)

	// 6. Practical example - Chat application
	fmt.Println("\n=== Chat Application Example ===")

	// Create chat room subscriber
	chatSub := rdb.Subscribe(ctx, "chat:room1")
	defer chatSub.Close()

	// Handle chat messages
	go func() {
		for {
			msg, err := chatSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving chat message: %v", err)
				return
			}
			fmt.Printf("[CHAT] %s\n", msg.Payload)
		}
	}()

	// Wait for subscriber to be ready
	time.Sleep(100 * time.Millisecond)

	// Simulate chat messages
	chatMessages := []string{
		"Alice: Hello everyone!",
		"Bob: Hi Alice, how are you?",
		"Charlie: Good morning!",
		"Alice: I'm doing great, thanks!",
		"Bob: Anyone up for a game?",
	}

	for _, message := range chatMessages {
		err = rdb.Publish(ctx, "chat:room1", message).Err()
		if err != nil {
			log.Fatalf("Error publishing chat message: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// 7. Notification system example
	fmt.Println("\n=== Notification System Example ===")

	// Subscribe to user notifications
	userSub := rdb.Subscribe(ctx, "notifications:user123")
	defer userSub.Close()

	// Handle notifications
	go func() {
		for {
			msg, err := userSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving notification: %v", err)
				return
			}
			fmt.Printf("[NOTIFICATION] %s\n", msg.Payload)
		}
	}()

	// Wait for subscriber to be ready
	time.Sleep(100 * time.Millisecond)

	// Send notifications
	notifications := []string{
		"New message from John",
		"Your order has been shipped",
		"Reminder: Meeting at 3 PM",
		"System maintenance scheduled",
	}

	for _, notification := range notifications {
		err = rdb.Publish(ctx, "notifications:user123", notification).Err()
		if err != nil {
			log.Fatalf("Error publishing notification: %v", err)
		}
		time.Sleep(300 * time.Millisecond)
	}

	// 8. Cleanup and final messages
	fmt.Println("\n=== Cleanup ===")

	// Unsubscribe from all channels
	err = subscriber.Unsubscribe(ctx)
	if err != nil {
		log.Fatalf("Error unsubscribing: %v", err)
	}

	err = multiSub.Unsubscribe(ctx)
	if err != nil {
		log.Fatalf("Error unsubscribing from multi: %v", err)
	}

	err = patternSub.PUnsubscribe(ctx)
	if err != nil {
		log.Fatalf("Error unsubscribing from pattern: %v", err)
	}

	err = chatSub.Unsubscribe(ctx)
	if err != nil {
		log.Fatalf("Error unsubscribing from chat: %v", err)
	}

	err = userSub.Unsubscribe(ctx)
	if err != nil {
		log.Fatalf("Error unsubscribing from notifications: %v", err)
	}

	fmt.Println("All subscriptions closed")
}
