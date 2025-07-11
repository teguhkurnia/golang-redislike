# Go In-Memory Redis-Like Store

A high-performance, thread-safe, in-memory key-value store implemented in Go that mimics Redis functionality. This project provides a complete Redis-compatible server with support for multiple data types, TTL management, and network protocol handling.

## Features

- **Thread-Safe**: All operations are safe for concurrent use with `sync.RWMutex`.
- **TCP Server**: Full Redis-compatible network protocol implementation (RESP).
- **Data Structures**: Support for Strings, Lists, Hashes, Sets, and Sorted Sets.
- **Data Persistence**: Append-Only File (AOF) to log all write operations for durability.
- **TTL Management**: Automatic key expiration with background cleanup.
- **Concurrent Connections**: Handles multiple clients concurrently.
- **Protocol Compatible**: Implements Redis Serialization Protocol (RESP).
- **Documentation**: Includes this `README.md` with setup and usage examples.

## Implemented Commands

#### String Commands
- `SET key value` - Set the value of a key
- `GET key` - Get the value of a key
- `DEL key [key ...]` - Delete one or more keys
- `INCR key` - Increment the integer value of a key by one
- `DECR key` - Decrement the integer value of a key by one

#### List Commands
- `LPUSH key value [value ...]` - Add elements to the beginning of a list
- `RPUSH key value [value ...]` - Add elements to the end of a list
- `LPOP key [count]` - Remove and return elements from the beginning of a list
- `RPOP key [count]` - Remove and return elements from the end of a list
- `LRANGE key start stop` - Get a range of elements from a list
- `LLEN key` - Get the length of a list

#### Hash Commands
- `HSET key field value [field value ...]` - Set hash field values
- `HGET key field` - Get the value of a hash field
- `HGETALL key` - Get all fields and values in a hash
- `HDEL key field [field ...]` - Delete hash fields

#### Set Commands
- `SADD key member [member ...]` - Add one or more members to a set
- `SREM key member [member ...]` - Remove one or more members from a set
- `SMEMBERS key` - Get all the members in a set
- `SISMEMBER key member` - Determine if a given value is a member of a set

#### Sorted Set Commands
- `ZADD key score member [score member ...]` - Add one or more members to a sorted set, or update its score if it already exists
- `ZRANGE key start stop [WITHSCORES]` - Returns a range of members in a sorted set, by index
- `ZREM key member [member ...]` - Remove one or more members from a sorted set

#### Time/TTL Commands
- `EXPIRE key seconds` - Set a key's time to live in seconds
- `TTL key` - Get the time to live for a key

#### Connection Commands
- `PING [message]` - Ping the server

## Roadmap

The following features are planned for future releases:

- **Data Persistence**:
  - [x] Append-Only File (AOF)
  - [ ] Snapshotting (RDB-style)
- **Memory Management**:
  - [ ] Eviction Policies (LRU, LFU)
- **Replication**:
  - [ ] Master-slave replication
- **Testing**:
  - [x] Comprehensive unit tests for all commands
  - [ ] Integration tests
  - [ ] Performance benchmarks
- **Additional Features**:
  - [ ] Pub/Sub
  - [ ] Transactions (MULTI/EXEC)

## Quick Start

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start on port `8080` by default.

### Connecting to the Server

You can connect using any Redis client or `telnet`:

```bash
# Using redis-cli
redis-cli -p 8080

# Using telnet
telnet localhost 8080
```

### Example Usage

```redis
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

```
cmd/server/          # Server entry point
internal/
├── server/          # TCP server implementation
├── store/           # In-memory data store
└── protocol/        # Redis protocol handling
    ├── parser/      # RESP protocol parser
    └── commands/    # Command implementations
```

## Development

### Building from Source

```bash
go build -o redis-like cmd/server/main.go
```

### Running with Custom Configuration

You can configure the server programmatically:

```go
// Example of running on a custom port
store := store.NewStore()
server := server.NewServer(":9000", store)
server.Start()
```

## Contributing

1.  Fork the repository.
2.  Create a feature branch.
3.  Commit your changes.
4.  Push to the branch.
5.  Create a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
