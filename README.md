# Go In-Memory Redis-Like Store

This project is a simple, in-memory, thread-safe key-value store implemented in Go. It mimics some of the basic functionalities of Redis, including support for string and list data types, as well as Time-To-Live (TTL) on keys.

## Features

- **Thread-Safe**: All operations are safe for concurrent use thanks to `sync.RWMutex`.
- **Key-Value Operations**: `SET`, `GET`, `DEL`, `EXISTS`.
- **Time-To-Live (TTL)**: Keys can be set with an expiration time (Unix timestamp).
- **List Data Type**: Supports `LPUSH`, `RPUSH`, `LRANGE`, `LPOP`, `RPOP`.
- **Expired Key Eviction**: A method is available to clear out all expired keys from the store.

## API and Usage

### Initialization

First, create a new store instance.

```go
import "your-repo/internal/store"

func main() {
    s := store.NewStore()
    // ... use the store
}
```

### String Commands

These commands operate on simple string values.

```go
// Set a simple key-value pair
s.Set("mykey", "hello")

// Get a value
// The second return value is false if the key doesn't exist or has expired.
val, exists := s.Get("mykey")
if exists {
    fmt.Printf("Got value: %s\n", val) // Output: Got value: hello
}

// Set a key with a 5-second TTL
// The TTL is a Unix timestamp.
ttl := time.Now().Add(5 * time.Second).Unix()
s.SetWithTTL("tempkey", "this will expire", ttl)
```

### List Commands

These commands operate on lists of strings.

```go
// Push values to the left (head) of a list.
// The list will be created if it doesn't exist.
s.LPush("mylist", []string{"world", "hello"}) // list is now ["hello", "world"]

// Push a value to the right (tail) of the list
s.RPush("mylist", []string{"!"}) // list is now ["hello", "world", "!"]

// Get a range of items from the list.
// 0 is the first element, -1 is the last.
items := s.LRange("mylist", 0, -1) // gets all items
fmt.Println(items) // Output: [hello world !]

// Pop one item from the left (head)
poppedLeft := s.LPop("mylist", 1)
fmt.Println(poppedLeft) // Output: [hello]
// The list is now ["world", "!"]

// Pop one item from the right (tail)
poppedRight := s.RPop("mylist", 1)
fmt.Println(poppedRight) // Output: [!]
// The list is now ["world"]
```

### Generic Commands

These commands work on any key.

```go
// Check if a key exists
keyExists := s.Exists("mylist") // true

// Delete a key
wasDeleted := s.Del("mylist") // true

// Check again
keyExists = s.Exists("mylist") // false
```

### Expired Key Management

The store checks for key expiration on `GET` operations. However, expired keys still consume memory until they are explicitly removed.

You should periodically call `ClearExpired` to purge all expired keys from the store. This is often done in a separate goroutine.

```go
func startCleanupRoutine(s *store.Store) {
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        for range ticker.C {
            s.ClearExpired()
            fmt.Println("Cleared expired keys.")
        }
    }()
}

// In your main function:
s := store.NewStore()
startCleanupRoutine(s)
```

## Upcoming Features

### Data Types
- **Hash Data Type**: Support for `HSET`, `HGET`, `HDEL`, `HKEYS`, `HVALS`, `HGETALL`
- **Set Data Type**: Support for `SADD`, `SREM`, `SMEMBERS`, `SISMEMBER`, `SCARD`
- **Sorted Set Data Type**: Support for `ZADD`, `ZREM`, `ZRANGE`, `ZRANK`, `ZSCORE`

### Time & Connection Management
- **TTL Commands**: `TTL`, `EXPIRE`, `EXPIREAT`, `PERSIST` for managing key expiration
- **Connection Commands**: `PING`, `ECHO`, `TIME`, `CLIENT` for connection management

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

## Roadmap

### Phase 1 (Current)
- ✅ Basic key-value operations (`SET`, `GET`, `DEL`, `EXISTS`)
- ✅ List operations (`LPUSH`, `RPUSH`, `LRANGE`, `LPOP`, `RPOP`)
- ✅ TTL support with expiration handling
- ✅ Thread-safe operations

### Phase 2 (Near Term)
- 🔄 Hash data type implementation
- 🔄 Time and connection management commands
- 📋 Enhanced TTL commands (`TTL`, `EXPIRE`, `EXPIREAT`, `PERSIST`)
- 📋 Basic metrics and monitoring

### Phase 3 (Medium Term)
- 📋 Set and Sorted Set data types
- 📋 Transaction support (`MULTI`, `EXEC`, `DISCARD`)
- 📋 Pub/Sub messaging system
- 📋 Memory management and eviction policies

### Phase 4 (Long Term)
- 📋 Lua scripting support
- 📋 Persistence layer with configurable backends
- 📋 Master-slave replication
- 📋 Advanced monitoring and performance tools

**Legend:**
- ✅ Complete
- 🔄 In Progress
- 📋 Planned
