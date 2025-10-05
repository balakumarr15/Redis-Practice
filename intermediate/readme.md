# Redis Intermediate Operations

This directory contains intermediate Redis operations covering complex data structures and operations.

## Files

### hashes.go
Demonstrates Redis Hash operations for structured data:
- **HSET**: Set hash field values
- **HGET**: Get hash field value
- **HGETALL**: Get all hash fields and values
- **HMGET**: Get multiple hash field values
- **HKEYS**: Get all hash field names
- **HVALS**: Get all hash field values
- **HEXISTS**: Check if hash field exists
- **HDEL**: Delete hash fields
- **HINCRBY**: Increment hash field by integer
- **HINCRBYFLOAT**: Increment hash field by float
- **HLEN**: Get number of fields in hash
- **HSETNX**: Set field only if it doesn't exist

**Key Commands:**
```redis
HSET key field value [field value ...]
HGET key field
HGETALL key
HMGET key field [field ...]
HKEYS key
HVALS key
HEXISTS key field
HDEL key field [field ...]
HINCRBY key field increment
HINCRBYFLOAT key field increment
HLEN key
HSETNX key field value
```

### lists.go
Covers Redis List operations for ordered collections:
- **LPUSH**: Push elements to left (head) of list
- **RPUSH**: Push elements to right (tail) of list
- **LRANGE**: Get range of elements from list
- **LLEN**: Get list length
- **LPOP**: Pop element from left (head)
- **RPOP**: Pop element from right (tail)
- **LINDEX**: Get element at specific index
- **LINSERT**: Insert element before or after pivot
- **LREM**: Remove elements from list
- **LSET**: Set element at specific index
- **LTRIM**: Trim list to specified range
- **RPOPLPUSH**: Pop from right of source, push to left of destination

**Key Commands:**
```redis
LPUSH key element [element ...]
RPUSH key element [element ...]
LRANGE key start stop
LLEN key
LPOP key
RPOP key
LINDEX key index
LINSERT key BEFORE|AFTER pivot element
LREM key count element
LSET key index element
LTRIM key start stop
RPOPLPUSH source destination
```

### sets.go
Demonstrates Redis Set operations for unique collections:
- **SADD**: Add members to set
- **SMEMBERS**: Get all members of set
- **SISMEMBER**: Check if member exists in set
- **SCARD**: Get cardinality (number of members) of set
- **SREM**: Remove members from set
- **SPOP**: Remove and return random member
- **SRANDMEMBER**: Get random member without removing
- **SUNION**: Union of sets
- **SINTER**: Intersection of sets
- **SDIFF**: Difference of sets
- **SUNIONSTORE**: Store union in new set
- **SINTERSTORE**: Store intersection in new set
- **SDIFFSTORE**: Store difference in new set
- **SMOVE**: Move member from one set to another

**Key Commands:**
```redis
SADD key member [member ...]
SMEMBERS key
SISMEMBER key member
SCARD key
SREM key member [member ...]
SPOP key [count]
SRANDMEMBER key [count]
SUNION key [key ...]
SINTER key [key ...]
SDIFF key [key ...]
SUNIONSTORE destination key [key ...]
SINTERSTORE destination key [key ...]
SDIFFSTORE destination key [key ...]
SMOVE source destination member
```

### sorted_sets.go
Covers Redis Sorted Set operations for ranked collections:
- **ZADD**: Add members to sorted set
- **ZRANGE**: Get range of members by rank
- **ZREVRANGE**: Get range in descending order
- **ZRANK**: Get rank of member
- **ZREVRANK**: Get rank in descending order
- **ZSCORE**: Get score of member
- **ZCARD**: Get cardinality of sorted set
- **ZCOUNT**: Count members within score range
- **ZRANGEBYSCORE**: Get members by score range
- **ZREM**: Remove members from sorted set
- **ZINCRBY**: Increment score of member
- **ZREMRANGEBYRANK**: Remove members by rank range
- **ZREMRANGEBYSCORE**: Remove members by score range

**Key Commands:**
```redis
ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
ZRANGE key start stop [WITHSCORES]
ZREVRANGE key start stop [WITHSCORES]
ZRANK key member
ZREVRANK key member
ZSCORE key member
ZCARD key
ZCOUNT key min max
ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]
ZREM key member [member ...]
ZINCRBY key increment member
ZREMRANGEBYRANK key start stop
ZREMRANGEBYSCORE key min max
```

## Use Cases

### Hashes
- **User profiles**: Store user information as key-value pairs
- **Object storage**: Represent objects with multiple attributes
- **Configuration**: Store application settings
- **Counters**: Track multiple metrics per entity

### Lists
- **Queues**: Implement FIFO queues for task processing
- **Stacks**: Implement LIFO stacks
- **Timelines**: Store chronologically ordered events
- **Message queues**: Handle inter-service communication

### Sets
- **Tags**: Store unique tags for content
- **User interests**: Track user preferences
- **Friends lists**: Manage social connections
- **Deduplication**: Remove duplicate entries

### Sorted Sets
- **Leaderboards**: Rank users by scores
- **Time series**: Store time-ordered data
- **Priority queues**: Process items by priority
- **Product ratings**: Rank products by rating

## Running the Examples

1. Make sure Redis is running on localhost:6379
2. Run any of the Go files:
   ```bash
   go run hashes.go
   go run lists.go
   go run sets.go
   go run sorted_sets.go
   ```

## Best Practices

1. **Choose the right data structure** for your use case
2. **Use appropriate commands** for your operations
3. **Consider memory usage** when storing large collections
4. **Use atomic operations** when possible
5. **Monitor performance** for large datasets
6. **Use expiration** to prevent memory leaks
7. **Consider data serialization** for complex objects
