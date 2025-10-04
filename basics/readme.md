# Redis Basics

This directory contains fundamental Redis operations and commands.

## Files

### set_get_expire.go
Demonstrates basic Redis string operations:
- **SET**: Store key-value pairs
- **GET**: Retrieve values by key
- **EXPIRE**: Set expiration time for existing keys
- **EXPIREAT**: Set expiration at specific timestamp
- **PERSIST**: Remove expiration from keys
- **TTL**: Check remaining time to live

**Key Commands:**
```redis
SET key value [EX seconds|PX milliseconds] [NX|XX]
GET key
EXPIRE key seconds
EXPIREAT key timestamp
PERSIST key
TTL key
```

### delete_keys.go
Covers key deletion and existence checking:
- **DEL**: Delete one or more keys
- **UNLINK**: Asynchronously delete keys (better for large keys)
- **EXISTS**: Check if keys exist
- **RENAME**: Rename a key
- **RENAMENX**: Rename only if new name doesn't exist

**Key Commands:**
```redis
DEL key [key ...]
UNLINK key [key ...]
EXISTS key [key ...]
RENAME oldkey newkey
RENAMENX oldkey newkey
```

### ttl_check.go
Focuses on Time To Live (TTL) operations:
- **TTL**: Get remaining seconds until expiration
- **PTTL**: Get remaining milliseconds until expiration
- **EXPIRE**: Set expiration on existing keys
- **EXPIREAT**: Set expiration at specific timestamp
- **PERSIST**: Remove expiration

**Key Commands:**
```redis
TTL key
PTTL key
EXPIRE key seconds
EXPIREAT key timestamp
PERSIST key
```

## TTL Return Values

| Return Value | Meaning |
|--------------|---------|
| Positive number | Key exists and has expiration set (seconds/milliseconds remaining) |
| -1 | Key exists but has no expiration set |
| -2 | Key does not exist |

## Running the Examples

1. Make sure Redis is running on localhost:6379
2. Run any of the Go files:
   ```bash
   go run set_get_expire.go
   go run delete_keys.go
   go run ttl_check.go
   ```

## Common Use Cases

- **Session Management**: Store session data with automatic expiration
- **Caching**: Cache frequently accessed data with TTL
- **Rate Limiting**: Track request counts with expiration
- **Temporary Data**: Store data that should automatically expire
- **Cleanup**: Use UNLINK for better performance when deleting large keys

## Best Practices

1. **Use appropriate expiration times** based on your use case
2. **Use UNLINK instead of DEL** for large keys to avoid blocking
3. **Check TTL before operations** to ensure keys haven't expired
4. **Use EXPIREAT for precise timing** when you need exact expiration times
5. **Monitor TTL values** to understand key lifecycle
