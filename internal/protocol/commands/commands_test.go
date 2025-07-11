package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teguhkurnia/redis-like/internal/store"
)

func TestSetAndGet(t *testing.T) {
	s := store.NewStore()
	cmdSet := &Command{Name: "SET", Args: [][]byte{[]byte("key"), []byte("value")}}
	cmdGet := &Command{Name: "GET", Args: [][]byte{[]byte("key")}}

	// Test SET
	result := handleSet(cmdSet, s)
	assert.Equal(t, "+OK\r\n", string(result))

	// Test GET
	result = handleGet(cmdGet, s)
	assert.Equal(t, "$5\r\nvalue\r\n", string(result))
}

func TestDel(t *testing.T) {
	s := store.NewStore()
	s.Set("key1", "value1")
	s.Set("key2", "value2")
	cmd := &Command{Name: "DEL", Args: [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}} // key3 does not exist

	result := handleDel(cmd, s)
	assert.Equal(t, ":2\r\n", string(result))

	_, exists := s.Get("key1")
	assert.False(t, exists)
	_, exists = s.Get("key2")
	assert.False(t, exists)
}

func TestIncrDecr(t *testing.T) {
	s := store.NewStore()
	s.Set("mykey", "10")
	cmdIncr := &Command{Name: "INCR", Args: [][]byte{[]byte("mykey")}}
	cmdDecr := &Command{Name: "DECR", Args: [][]byte{[]byte("mykey")}}

	// Test INCR
	result := handleIncr(cmdIncr, s)
	assert.Equal(t, ":11\r\n", string(result))

	// Test DECR
	result = handleDecr(cmdDecr, s)
	assert.Equal(t, ":10\r\n", string(result))
}

func TestPing(t *testing.T) {
	s := store.NewStore()
	cmd1 := &Command{Name: "PING", Args: [][]byte{}}
	cmd2 := &Command{Name: "PING", Args: [][]byte{[]byte("hello")}}

	// Test PING
	result := handlePing(cmd1, s)
	assert.Equal(t, "+PONG\r\n", string(result))

	// Test PING with message
	result = handlePing(cmd2, s)
	assert.Equal(t, "$5\r\nhello\r\n", string(result))
}

func TestListCommands(t *testing.T) {
	s := store.NewStore()
	cmdLPush := &Command{Name: "LPUSH", Args: [][]byte{[]byte("mylist"), []byte("world"), []byte("hello")}}
	cmdRPush := &Command{Name: "RPUSH", Args: [][]byte{[]byte("mylist"), []byte("!")}}
	cmdLRange := &Command{Name: "LRANGE", Args: [][]byte{[]byte("mylist"), []byte("0"), []byte("-1")}}
	cmdLLen := &Command{Name: "LLEN", Args: [][]byte{[]byte("mylist")}}
	cmdLPop := &Command{Name: "LPOP", Args: [][]byte{[]byte("mylist")}}
	cmdRPop := &Command{Name: "RPOP", Args: [][]byte{[]byte("mylist")}}

	// Test LPUSH
	result := handleLPush(cmdLPush, s)
	assert.Equal(t, ":2\r\n", string(result))

	// Test RPUSH
	result = handleRPush(cmdRPush, s)
	assert.Equal(t, ":3\r\n", string(result))

	// Test LRANGE
	result = handleLRange(cmdLRange, s)
	assert.Equal(t, "*3\r\n$5\r\nhello\r\n$5\r\nworld\r\n$1\r\n!\r\n", string(result))

	// Test LLEN
	result = handleLLen(cmdLLen, s)
	assert.Equal(t, ":3\r\n", string(result))

	// Test LPOP
	result = handleLPop(cmdLPop, s)
	assert.Equal(t, "*1\r\n$5\r\nhello\r\n", string(result))

	// Test RPOP
	result = handleRPop(cmdRPop, s)
	assert.Equal(t, "*1\r\n$1\r\n!\r\n", string(result))
}

func TestHashCommands(t *testing.T) {
	s := store.NewStore()
	cmdHSet := &Command{Name: "HSET", Args: [][]byte{[]byte("myhash"), []byte("field1"), []byte("Hello"), []byte("field2"), []byte("World")}}
	cmdHGet := &Command{Name: "HGET", Args: [][]byte{[]byte("myhash"), []byte("field1")}}
	cmdHGetAll := &Command{Name: "HGETALL", Args: [][]byte{[]byte("myhash")}}
	cmdHDel := &Command{Name: "HDEL", Args: [][]byte{[]byte("myhash"), []byte("field1"), []byte("field2")}}

	// Test HSET
	result := handleHSet(cmdHSet, s)
	assert.Equal(t, ":2\r\n", string(result))

	// Test HGET
	result = handleHGet(cmdHGet, s)
	assert.Equal(t, "$5\r\nHello\r\n", string(result))

	// Test HGETALL
	result = handleHGetAll(cmdHGetAll, s)
	assert.Contains(t, string(result), "*4\r\n")
	assert.Contains(t, string(result), "$6\r\nfield1\r\n$5\r\nHello\r\n")
	assert.Contains(t, string(result), "$6\r\nfield2\r\n$5\r\nWorld\r\n")

	// Test HDEL
	result = handleHDel(cmdHDel, s)
	assert.Equal(t, ":2\r\n", string(result))
}

func TestSetCommands(t *testing.T) {
	s := store.NewStore()
	cmdSAdd := &Command{Name: "SADD", Args: [][]byte{[]byte("myset"), []byte("member1"), []byte("member2")}}
	cmdSMembers := &Command{Name: "SMEMBERS", Args: [][]byte{[]byte("myset")}}
	cmdSIsMember := &Command{Name: "SISMEMBER", Args: [][]byte{[]byte("myset"), []byte("member1")}}
	cmdSRem := &Command{Name: "SREM", Args: [][]byte{[]byte("myset"), []byte("member1"), []byte("member3")}}

	// Test SADD
	result := handleSAdd(cmdSAdd, s)
	assert.Equal(t, ":2\r\n", string(result))

	// Test SMEMBERS
	result = handleSMembers(cmdSMembers, s)
	assert.Contains(t, string(result), "*2\r\n")
	assert.Contains(t, string(result), "$7\r\nmember1\r\n")
	assert.Contains(t, string(result), "$7\r\nmember2\r\n")

	// Test SISMEMBER
	result = handleSIsMember(cmdSIsMember, s)
	assert.Equal(t, ":1\r\n", string(result))

	// Test SREM
	result = handleSRem(cmdSRem, s)
	assert.Equal(t, ":1\r\n", string(result))
}

func TestSortedSetCommands(t *testing.T) {
	s := store.NewStore()
	cmdZAdd := &Command{Name: "ZADD", Args: [][]byte{[]byte("myzset"), []byte("1"), []byte("one"), []byte("2"), []byte("two")}}
	cmdZRange := &Command{Name: "ZRANGE", Args: [][]byte{[]byte("myzset"), []byte("0"), []byte("-1")}}
	cmdZRem := &Command{Name: "ZREM", Args: [][]byte{[]byte("myzset"), []byte("one"), []byte("three")}}

	// Test ZADD
	result := handleZAdd(cmdZAdd, s)
	assert.Equal(t, ":2\r\n", string(result))

	// Test ZRANGE
	result = handleZRange(cmdZRange, s)
	assert.Equal(t, "*2\r\n$3\r\none\r\n$3\r\ntwo\r\n", string(result))

	// Test ZREM
	result = handleZRem(cmdZRem, s)
	assert.Equal(t, ":1\r\n", string(result))
}

func TestTimeCommands(t *testing.T) {
	s := store.NewStore()
	s.Set("key", "value")
	cmdExpire := &Command{Name: "EXPIRE", Args: [][]byte{[]byte("key"), []byte("10")}}
	cmdTTL := &Command{Name: "TTL", Args: [][]byte{[]byte("key")}}

	// Test EXPIRE
	result := handleExpire(cmdExpire, s)
	assert.Equal(t, "+OK\r\n", string(result))

	// Test TTL
	result = handleTTL(cmdTTL, s)
	assert.Equal(t, ":10\r\n", string(result))
}
