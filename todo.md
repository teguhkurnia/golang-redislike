# Project Todo List

## Core Features

- [x] Implement basic data structures:
  - [x] String
  - [x] List
  - [x] Hash
  - [x] Set
  - [x] Sorted Set
- [x] Implement commands for each data structure:
  - [x] String Commands (SET, GET, DEL, INCR, DECR)
  - [x] List Commands (LPUSH, RPUSH, LPOP, RPOP, LLEN, LRANGE)
  - [x] Hash Commands (HSET, HGET, HDEL, HGETALL)
  - [x] Set Commands (SADD, SREM, SMEMBERS, SISMEMBER)
  - [x] Sorted Set Commands (ZADD, ZRANGE, ZREM)
- [x] Handle concurrent client connections.
- [x] Parser for the communication protocol (RESP - Redis Serialization Protocol).

## Data Persistence

- [x] Append-Only File (AOF) to log every write operation.
- [ ] Snapshotting (like Redis RDB) to save state to disk.

## Memory Management

- [ ] Eviction Policies when memory is full:
  - [ ] LRU (Least Recently Used)
  - [ ] LFU (Least Frequently Used)
- [x] TTL (Time To Live) support for keys.

## Replication

- [ ] Implement master-slave replication.

## Testing

- [ ] Unit tests for all commands and data structures.
- [ ] Integration tests to simulate client-server interaction.
- [ ] Benchmark tests to measure performance.

## Documentation

- [x] Update `README.md` with instructions on how to run and use the project.
- [x] Provide usage examples.
- [ ] API documentation for each command.

## Additional Features

- [ ] Implement Pub/Sub (Publish/Subscribe).
- [ ] Transactions (MULTI/EXEC).