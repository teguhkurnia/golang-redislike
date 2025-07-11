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
