# Redis Practice Project

A comprehensive Redis practice project with examples covering basic operations, intermediate data structures, advanced features, and real-world projects.

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- Redis 6.0 or higher
- Docker (optional, for containerized Redis)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Redis-Practice
```

2. Install dependencies:
```bash
go mod tidy
```

3. Start Redis (choose one option):

**Option A: Using Docker (Recommended)**
```bash
cd scripts
docker-compose up -d
```

**Option B: Local Redis installation**
```bash
# On macOS
brew install redis
brew services start redis

# On Ubuntu/Debian
sudo apt-get install redis-server
sudo systemctl start redis-server

# On Windows
# Download and install Redis from https://github.com/microsoftarchive/redis/releases
```

4. Seed sample data (optional):
```bash
go run scripts/seed_data.go
```

## 📁 Project Structure

```
Redis-Practice/
├── basics/                       # Fundamental Redis operations
│   ├── set_get_expire.go         # Basic SET, GET, TTL usage
│   ├── delete_keys.go            # DEL, EXISTS commands
│   ├── ttl_check.go              # Checking key expiration
│   └── readme.md                 # Notes about basic commands
│
├── intermediate/                 # Slightly advanced use cases
│   ├── hashes.go                 # HSET, HGET, HMSET for structured data
│   ├── lists.go                  # LPUSH, RPUSH, LPOP, etc.
│   ├── sets.go                   # SADD, SREM, etc.
│   ├── sorted_sets.go            # ZADD, ZRANGE
│   └── readme.md
│
├── advanced/                     # Complex Redis features
│   ├── pubsub.go                 # Redis Publish/Subscribe
│   ├── streams.go                # Redis Streams
│   ├── pipeline.go               # Command pipelining
│   ├── transactions.go           # MULTI/EXEC usage
│   └── readme.md
│
├── projects/                     # Small projects combining concepts
│   ├── session_manager.go        # Manage sessions with TTL
│   ├── rate_limiter.go          # Implement rate limiting
│   ├── leaderboard.go           # Leaderboard with sorted sets
│   └── chat_pubsub.go           # Simple chat app using Pub/Sub
│
├── tests/                        # Unit tests for practice
│   ├── basics_test.go
│   ├── advanced_test.go
│   └── projects_test.go
│
├── scripts/                      # Helper scripts
│   ├── docker-compose.yml        # Run Redis with Docker
│   ├── redis.conf                # Redis configuration
│   └── seed_data.go              # Script to seed sample Redis data
│
├── go.mod                        # Go module file
├── go.sum                        # Go dependencies lock file
├── .gitignore                    # Ignore binaries, logs, etc.
└── README.md                     # This file
```

## 🎯 Learning Path

### 1. Basics (Start Here)
Learn fundamental Redis operations:
```bash
go run basics/set_get_expire.go
go run basics/delete_keys.go
go run basics/ttl_check.go
```

**Key Concepts:**
- String operations (SET, GET)
- Key expiration (EXPIRE, TTL)
- Key deletion (DEL, EXISTS)
- Key renaming (RENAME, RENAMENX)

### 2. Intermediate
Explore Redis data structures:
```bash
go run intermediate/hashes.go
go run intermediate/lists.go
go run intermediate/sets.go
go run intermediate/sorted_sets.go
```

**Key Concepts:**
- Hashes for structured data
- Lists for ordered collections
- Sets for unique collections
- Sorted sets for ranked data

### 3. Advanced
Master complex Redis features:
```bash
go run advanced/pubsub.go
go run advanced/streams.go
go run advanced/pipeline.go
go run advanced/transactions.go
```

**Key Concepts:**
- Pub/Sub for real-time messaging
- Streams for event sourcing
- Pipelines for batch operations
- Transactions for atomic operations

### 4. Projects
Build real-world applications:
```bash
go run projects/session_manager.go
go run projects/rate_limiter.go
go run projects/leaderboard.go
go run projects/chat_pubsub.go
```

**Key Concepts:**
- Session management with TTL
- Rate limiting algorithms
- Leaderboards with sorted sets
- Real-time chat with Pub/Sub

## 🧪 Testing

Run the test suite:
```bash
# Run all tests
go test ./tests/...

# Run specific test files
go test ./tests/basics_test.go
go test ./tests/advanced_test.go
go test ./tests/projects_test.go

# Run tests with verbose output
go test -v ./tests/...

# Run tests with coverage
go test -cover ./tests/...
```

## 🐳 Docker Setup

The project includes Docker Compose configuration for easy Redis setup:

```bash
# Start Redis and management tools
cd scripts
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Clean up (removes data)
docker-compose down -v
```

**Included Services:**
- Redis 7 (port 6379)
- Redis Commander (port 8081) - Web UI for Redis
- Redis Insight (port 8001) - Advanced Redis management tool

## 📊 Sample Data

Seed the database with realistic sample data:
```bash
go run scripts/seed_data.go
```

This will populate Redis with:
- User profiles and metadata
- Product catalog with ratings
- Game leaderboard with scores
- Active user sessions
- Chat room history
- Event stream data
- Rate limiting counters

## 🔧 Configuration

### Redis Configuration
The project includes a custom Redis configuration (`scripts/redis.conf`) optimized for practice:
- Memory limit: 256MB
- Persistence enabled (RDB + AOF)
- Slow log monitoring
- Latency monitoring

### Go Configuration
- Module: `Redis`
- Go version: 1.25.0
- Dependencies: `github.com/redis/go-redis/v9`

## 📚 Learning Resources

### Redis Documentation
- [Redis Official Documentation](https://redis.io/docs/)
- [Redis Commands Reference](https://redis.io/commands/)
- [Redis Data Types](https://redis.io/docs/data-types/)

### Go Redis Client
- [go-redis Documentation](https://redis.uptrace.dev/)
- [go-redis GitHub](https://github.com/redis/go-redis)

### Best Practices
- [Redis Best Practices](https://redis.io/docs/manual/patterns/)
- [Redis Performance Tuning](https://redis.io/docs/manual/performance/)
- [Redis Security](https://redis.io/docs/manual/security/)

## 🚀 Advanced Usage

### Custom Redis Configuration
Modify `scripts/redis.conf` to customize Redis behavior:
```bash
# Edit configuration
vim scripts/redis.conf

# Restart Redis with new config
docker-compose restart redis
```

### Performance Testing
Run performance tests to measure Redis performance:
```bash
# Basic performance test
go test -run TestPerformance ./tests/advanced_test.go

# Concurrent operations test
go test -run TestConcurrent ./tests/projects_test.go
```

### Monitoring
Monitor Redis performance and usage:
```bash
# Connect to Redis CLI
redis-cli

# Monitor commands in real-time
MONITOR

# Get Redis info
INFO

# Check memory usage
INFO memory

# Check slow log
SLOWLOG GET 10
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Redis Connection Failed**
```bash
# Check if Redis is running
redis-cli ping

# Check Redis logs
docker-compose logs redis
```

**Port Already in Use**
```bash
# Find process using port 6379
lsof -i :6379

# Kill the process
kill -9 <PID>
```

**Memory Issues**
```bash
# Check Redis memory usage
redis-cli INFO memory

# Clear all data
redis-cli FLUSHALL
```

**Go Module Issues**
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod download
```

## 📞 Support

If you encounter any issues or have questions:
1. Check the troubleshooting section above
2. Review the Redis documentation
3. Open an issue on GitHub
4. Check existing issues and discussions

## 🎉 Next Steps

After completing this practice project:
1. Build your own Redis-based application
2. Explore Redis modules and extensions
3. Learn about Redis clustering and replication
4. Study Redis performance optimization
5. Contribute to open-source Redis projects

Happy coding! 🚀
