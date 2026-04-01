# Redis Advanced Features

This directory contains advanced Redis features and operations for complex use cases.

## Files

### pubsub.go
Demonstrates Redis Pub/Sub (Publish/Subscribe) operations:
- **PUBLISH**: Publish messages to channels
- **SUBSCRIBE**: Subscribe to channels
- **UNSUBSCRIBE**: Unsubscribe from channels
- **PSUBSCRIBE**: Subscribe to channel patterns
- **PUNSUBSCRIBE**: Unsubscribe from patterns
- **PUBSUB**: Get information about channels and patterns

**Key Commands:**
```redis
PUBLISH channel message
SUBSCRIBE channel [channel ...]
UNSUBSCRIBE [channel [channel ...]]
PSUBSCRIBE pattern [pattern ...]
PUNSUBSCRIBE [pattern [pattern ...]]
PUBSUB CHANNELS [pattern]
PUBSUB NUMSUB [channel [channel ...]]
PUBSUB NUMPAT
```

### streams.go
Covers Redis Streams for event sourcing and messaging:
- **XADD**: Add entries to stream
- **XREAD**: Read entries from stream
- **XRANGE**: Get entries by ID range
- **XREVRANGE**: Get entries in reverse order
- **XGROUP**: Create consumer groups
- **XREADGROUP**: Read from consumer group
- **XACK**: Acknowledge processed messages
- **XPENDING**: Check pending messages
- **XCLAIM**: Claim pending messages
- **XDEL**: Delete entries from stream
- **XTRIM**: Trim stream to specified length
- **XLEN**: Get stream length

**Key Commands:**
```redis
XADD stream * field value [field value ...]
XREAD STREAMS stream [stream ...] id [id ...]
XRANGE stream start end [COUNT count]
XREVRANGE stream end start [COUNT count]
XGROUP CREATE stream groupname id
XREADGROUP GROUP group consumer STREAMS stream [stream ...] id [id ...]
XACK stream group id [id ...]
XPENDING stream group
XCLAIM stream group consumer min-idle-time id [id ...]
XDEL stream id [id ...]
XTRIM stream MAXLEN count
XLEN stream
```

### pipeline.go
Demonstrates Redis Pipeline for batch operations:
- **Pipeline**: Batch multiple commands
- **TxPipeline**: Transactional pipeline
- **Performance optimization**: Reduce round trips
- **Error handling**: Handle pipeline errors
- **Mixed operations**: Combine different command types

**Key Concepts:**
- Use `rdb.Pipeline()` for batch operations
- Use `rdb.TxPipeline()` for transactional operations
- Execute with `pipe.Exec(ctx)`
- Process results from `[]redis.Cmder`

### transactions.go
Covers Redis Transactions for atomic operations:
- **MULTI/EXEC**: Execute commands atomically
- **WATCH**: Monitor keys for changes
- **DISCARD**: Cancel transaction
- **Error handling**: Handle transaction failures
- **Conditional operations**: Execute based on conditions

**Key Concepts:**
- Use `rdb.Watch()` for conditional transactions
- Use `rdb.TxPipeline()` for transactional operations
- Transactions are atomic - all or nothing
- WATCH prevents race conditions

## Use Cases

### Pub/Sub
- **Real-time messaging**: Chat applications, notifications
- **Event broadcasting**: System events, updates
- **Microservices communication**: Inter-service messaging
- **Live updates**: Real-time data synchronization

### Streams
- **Event sourcing**: Store and replay events
- **Message queues**: Reliable message processing
- **Log aggregation**: Collect and process logs
- **Time series data**: Store time-ordered events

### Pipeline
- **Batch operations**: Reduce network round trips
- **Performance optimization**: Improve throughput
- **Bulk data loading**: Load large datasets
- **Atomic operations**: Group related commands

### Transactions
- **Atomic operations**: Ensure data consistency
- **Race condition prevention**: Use WATCH for safety
- **Complex business logic**: Multi-step operations
- **Data integrity**: Maintain referential integrity

## Running the Examples

1. Make sure Redis is running on localhost:6379
2. Run any of the Go files:
   ```bash
   go run pubsub.go
   go run streams.go
   go run pipeline.go
   go run transactions.go
   ```

## Best Practices

### Pub/Sub
1. **Handle disconnections** gracefully
2. **Use patterns** for dynamic channel subscriptions
3. **Monitor subscriber count** for debugging
4. **Clean up subscriptions** when done

### Streams
1. **Use consumer groups** for reliable processing
2. **Acknowledge messages** after processing
3. **Handle failures** with XCLAIM
4. **Trim streams** to prevent memory issues

### Pipeline
1. **Batch related operations** together
2. **Handle errors** individually
3. **Use for performance** when appropriate
4. **Consider memory usage** for large batches

### Transactions
1. **Use WATCH** for conditional operations
2. **Handle failures** gracefully
3. **Keep transactions short** to avoid blocking
4. **Test error scenarios** thoroughly

## Performance Considerations

- **Pipeline**: 2-10x faster than individual commands
- **Transactions**: Slightly slower than individual commands
- **Pub/Sub**: Very fast for real-time messaging
- **Streams**: Good for high-throughput event processing

## Error Handling

- **Pipeline errors**: Check individual command results
- **Transaction errors**: Handle WATCH failures
- **Pub/Sub errors**: Handle connection issues
- **Stream errors**: Handle consumer group failures
