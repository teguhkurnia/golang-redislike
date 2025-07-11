# Go In-Memory Redis-Like Store

A high-performance, thread-safe, in-memory key-value store implemented in Go that mimics Redis functionality. This project provides a complete Redis-compatible server with support for multiple data types, TTL management, and network protocol handling.

## Features

### Core Functionality

- **Thread-Safe**: All operations are safe for concurrent use with `sync.RWMutex`
- **TCP Server**: Full Redis-compatible network protocol implementation
- **Multiple Data Types**: Support for strings, lists, and hash tables
- **TTL Management**: Automatic key expiration with background cleanup
- **Protocol Compatible**: Implements Redis Serialization Protocol (RESP)

### Implemented Commands

#### String Commands

- `SET key value` - Set the value of a key
- `GET key` - Get the value of a key
- `DEL key [key ...]` - Delete one or more keys

#### List Commands

- `LPUSH key value [value ...]` - Add elements to the beginning of a list
- `RPUSH key value [value ...]` - Add elements to the end of a list
- `LPOP key [count]` - Remove and return elements from the beginning of a list
- `RPOP key [count]` - Remove and return elements from the end of a list
- `LRANGE key start stop` - Get a range of elements from a list

#### Hash Commands

- `HSET key field value [field value ...]` - Set hash field values
- `HGET key field` - Get the value of a hash field
- `HGETALL key` - Get all fields and values in a hash
- `HDEL key field [field ...]` - Delete hash fields

#### Time/TTL Commands

- `EXPIRE key seconds` - Set a key's time to live in seconds
- `TTL key` - Get the time to live for a key

#### Connection Commands

- `PING [message]` - Ping the server

## Quick Start

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default and includes automatic cleanup of expired keys every second.

### Connecting to the Server

You can connect using any Redis client or telnet:

```bash
# Using redis-cli (if installed)
redis-cli -p 8080

# Using telnet
telnet localhost 8080
```

### Example Usage

```bash
# String operations
SET mykey "Hello World"
GET mykey
DEL mykey

# List operations
LPUSH mylist "item1" "item2"
RPUSH mylist "item3"
LRANGE mylist 0 -1
LPOP mylist

# Hash operations
HSET user:1 name "John" age "30"
HGET user:1 name
HGETALL user:1
HDEL user:1 age

# TTL operations
SET tempkey "temporary"
EXPIRE tempkey 60
TTL tempkey
```

## Architecture

### Project Structure

```
cmd/server/          # Server entry point
internal/
├── server/          # TCP server implementation
├── store/           # In-memory data store
└── protocol/        # Redis protocol handling
    ├── parser/      # RESP protocol parser
    └── commands/    # Command implementations
```

### Key Components

1. **Store**: Thread-safe in-memory storage with support for multiple data types
2. **Server**: TCP server handling client connections and protocol parsing
3. **Protocol**: Redis Serialization Protocol (RESP) implementation
4. **Commands**: Modular command system with proper error handling

## API Reference

### Using the Store Directly

```go
import "github.com/teguhkurnia/redis-like/internal/store"

func main() {
    store := store.NewStore()

    // String operations
    store.Set("key", "value")
    value, exists := store.Get("key")

    // List operations
    store.LPush("mylist", []string{"item1", "item2"})
    items := store.LRange("mylist", 0, -1)

    // Hash operations
    store.HSet("user:1", "name", "John")
    name, exists := store.HGet("user:1", "name")

    // TTL operations
    store.Expire("key", 60) // 60 seconds
    ttl := store.TTL("key")
}
```

### Error Handling

All network operations return proper Redis error responses:

- `-ERR wrong number of arguments for 'command' command`
- `-ERR invalid expire time in 'expire' command`
- `$-1` for non-existent keys (null bulk string)

## Performance Features

- **Memory Efficient**: Optimized data structures for minimal memory usage
- **Concurrent Safe**: Lock-free reads where possible, minimal lock contention
- **Background Cleanup**: Automatic removal of expired keys
- **Connection Pooling**: Efficient handling of multiple client connections

## Development

### Building

```bash
go build -o redis-like cmd/server/main.go
```

### Running with Custom Configuration

```go
store := store.NewStore()
server := server.NewServer(":9000", store) // Custom port
server.Start()
```

## Upcoming Features

### Data Types

- **Set Data Type**: Support for `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`
- **Sorted Set Data Type**: Support for `ZADD`, `ZREM`, `ZRANGE`, `ZRANK`, `ZSCORE`

### Enhanced Commands

- **Extended TTL**: `EXPIREAT`, `PERSIST`, `PEXPIRE`, `PTTL`
- **Key Operations**: `EXISTS`, `KEYS`, `RENAME`, `TYPE`
- **List Operations**: `LLEN`, `LINDEX`, `LSET`, `LTRIM`
- **Hash Operations**: `HEXISTS`, `HKEYS`, `HVALS`, `HLEN`

### Advanced Features

- **Transactions**: Support for `MULTI`, `EXEC`, `DISCARD`, `WATCH`
- **Pub/Sub**: Support for `PUBLISH`, `SUBSCRIBE`, `UNSUBSCRIBE`
- **Lua Scripting**: Support for `EVAL`, `EVALSHA` for custom operations
- **Persistence**: Optional disk persistence with configurable intervals
- **Memory Management**: Configurable memory limits and eviction policies
- **Replication**: Master-slave replication support

### Performance & Monitoring

- **Metrics Collection**: Built-in metrics for operations, memory usage, and performance
- **Connection Pooling**: Improved connection management for high-concurrency scenarios
- **Benchmarking Tools**: Performance testing utilities
- **Admin Commands**: `INFO`, `CONFIG`, `FLUSHALL`, `FLUSHDB`

## Roadmap

### Phase 1 (Current) ✅

- ✅ Basic string operations (`SET`, `GET`, `DEL`)
- ✅ List operations (`LPUSH`, `RPUSH`, `LRANGE`, `LPOP`, `RPOP`)
- ✅ Hash operations (`HSET`, `HGET`, `HGETALL`, `HDEL`)
- ✅ TTL support (`EXPIRE`, `TTL`)
- ✅ TCP server with RESP protocol
- ✅ Connection management (`PING`)
- ✅ Thread-safe operations
- ✅ Automatic expired key cleanup

### Phase 2 (Near Term)

- 📋 Set and Sorted Set data types
- 📋 Extended key and TTL commands
- 📋 Additional list and hash operations
- 📋 Basic metrics and monitoring
- 📋 Configuration management

### Phase 3 (Medium Term)

- 📋 Transaction support (`MULTI`, `EXEC`, `DISCARD`)
- 📋 Pub/Sub messaging system
- 📋 Memory management and eviction policies
- 📋 Admin and info commands
- 📋 Enhanced error handling

### Phase 4 (Long Term)

- 📋 Lua scripting support
- 📋 Persistence layer with configurable backends
- 📋 Master-slave replication
- 📋 Advanced monitoring and performance tools
- 📋 Clustering support

**Legend:**

- ✅ Complete
- 🔄 In Progress
- 📋 Planned

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is open source and available under the [MIT License](LICENSE).
