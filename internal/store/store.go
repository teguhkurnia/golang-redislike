package store

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Data struct {
	Value any
	TTL   int64
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Data
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Data),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Data{
		Value: value,
		TTL:   0,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.data[key]

	if !exists || (value.TTL > 0 && value.TTL < time.Now().Unix()) {
		return "", false
	}
	if value.TTL > 0 && value.TTL < time.Now().Unix() {
		s.Del(key) // remove expired key
		return "", false
	}

	if _, ok := value.Value.(string); !ok {
		return "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n", false // value is not a string
	}

	return value.Value.(string), exists
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; !exists {
		return false
	}

	delete(s.data, key)
	return true
}

func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists
}

func (s *Store) Incr(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[key]
	if !exists {
		value = Data{
			Value: "0",
			TTL:   0,
		}
	}
	// Increment the value
	intValue, err := strconv.Atoi(fmt.Sprintf("%v", value.Value))
	if err != nil {
		return 0, false // Value is not an integer
	}
	value.Value = intValue + 1
	s.data[key] = value

	return intValue + 1, true
}

func (s *Store) Decr(key string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, exists := s.data[key]
	if !exists {
		value = Data{
			Value: "0",
			TTL:   0,
		}
	}

	intValue, err := strconv.Atoi(fmt.Sprintf("%v", value.Value))
	if err != nil {
		return 0, false // Value is not an integer
	}
	value.Value = intValue - 1
	s.data[key] = value

	return intValue - 1, true
}

func (s *Store) Expire(key string, seconds int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; !exists {
		return 0
	} else {
		value := s.data[key]
		value.TTL = time.Now().Unix() + int64(seconds)
		s.data[key] = value
	}

	return 1
}

func (s *Store) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if data, exists := s.data[key]; exists {
		if data.TTL <= 0 {
			return -1 // No expiration set
		}
		ttl := data.TTL - time.Now().Unix()
		if ttl < 0 {
			return -2 // Key has expired
		}
		return int(ttl)
	}
	return -2 // Key does not exist
}

// LIST
func (s *Store) LPush(key string, values []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentValues := s.lrangeInternal(key, 0, -1)
	if currentValues == nil {
		currentValues = []string{}
	}
	// reverse the values to maintain order
	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	s.data[key] = Data{
		Value: append(values, currentValues...),
		TTL:   0,
	}

	return len(s.data[key].Value.([]string))
}

func (s *Store) RPush(key string, values []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentValues := s.lrangeInternal(key, 0, -1)
	if currentValues == nil {
		currentValues = []string{}
	}
	s.data[key] = Data{
		Value: append(currentValues, values...),
		TTL:   0,
	}
	return len(s.data[key].Value.([]string))
}

func (s *Store) lrangeInternal(key string, start, end int) []string {
	entry, exists := s.data[key]
	if !exists {
		return nil
	}

	list := entry.Value.([]string)
	listLen := len(list)

	// handle the negative indices
	if start < 0 {
		start = listLen + start
	}
	if end < 0 {
		end = listLen + end
	}

	// Konversi [][]byte ke []string untuk dikembalikan
	results := make([]string, 0, end-start+1)
	for _, v := range list[start : end+1] {
		results = append(results, string(v))
	}
	return results
}

func (s *Store) LRange(key string, start, end int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lrangeInternal(key, start, end)
}

func (s *Store) LPop(key string, count int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return nil
	}

	values := s.data[key].Value.([]string)
	if len(values) == 0 {
		return nil
	}

	if count > len(values) {
		count = len(values)
	}

	poppedValues := values[:count]
	newData := s.data[key]
	newData.Value = values[count:]
	if len(newData.Value.([]string)) == 0 {
		delete(s.data, key) // remove the key if no values left
	} else {
		s.data[key] = newData
	}

	return poppedValues
}

func (s *Store) RPop(key string, count int) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return nil
	}

	values := s.data[key].Value.([]string)
	if len(values) == 0 {
		return nil
	}
	if count > len(values) {
		count = len(values)
	}
	poppedValues := values[len(values)-count:]
	newData := s.data[key]
	newData.Value = values[:len(values)-count]
	if len(newData.Value.([]string)) == 0 {
		delete(s.data, key) // remove the key if no values left
	} else {
		s.data[key] = newData
	}
	return poppedValues
}

func (s *Store) LLen(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.data[key]; !exists {
		return 0, true
	}

	values, ok := s.data[key].Value.([]string)
	if !ok {
		return 0, false
	}
	return len(values), true
}

// HASH
func (s *Store) HSet(key, field, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash, exists := s.data[key]
	if !exists {
		hash = Data{
			Value: make(map[string]string),
			TTL:   0,
		}
	}
	hash.Value.(map[string]string)[field] = value
	s.data[key] = hash
	return 1
}

func (s *Store) HGet(key, field string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hash, exists := s.data[key]
	if !exists {
		return "", false
	}
	value, ok := hash.Value.(map[string]string)[field]
	if !ok {
		return "", false
	}
	return value, true
}

func (s *Store) HGetAll(key string) (map[string]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hash, exists := s.data[key]
	if !exists {
		return nil, false
	}
	return hash.Value.(map[string]string), true
}

func (s *Store) HDel(key, field string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash, exists := s.data[key]
	if !exists {
		return 0
	}
	if _, ok := hash.Value.(map[string]string)[field]; !ok {
		return 0
	}
	delete(hash.Value.(map[string]string), field)
	if len(hash.Value.(map[string]string)) == 0 {
		delete(s.data, key) // remove the key if no fields left
	} else {
		s.data[key] = hash
	}
	return 1
}

// SET
func (s *Store) SAdd(key string, members []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; !exists {
		s.data[key] = Data{
			Value: make(map[string]struct{}),
			TTL:   0,
		}
	}
	set := s.data[key].Value.(map[string]struct{})
	count := 0
	for _, member := range members {
		if _, ok := set[member]; !ok {
			set[member] = struct{}{}
			count++
		}
	}
	s.data[key] = Data{
		Value: set,
		TTL:   0,
	}
	return count
}

func (s *Store) SRem(key string, members []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	set, exists := s.data[key]
	if !exists {
		return 0
	}

	count := 0
	for _, member := range members {
		if _, ok := set.Value.(map[string]struct{})[member]; ok {
			delete(set.Value.(map[string]struct{}), member)
			count++
		}
	}

	if len(set.Value.(map[string]struct{})) == 0 {
		delete(s.data, key) // remove the key if no members left
	} else {
		s.data[key] = set
	}

	return count
}

func (s *Store) SMembers(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, exists := s.data[key]
	if !exists {
		return []string{}
	}

	members := make([]string, 0, len(set.Value.(map[string]struct{})))
	for member := range set.Value.(map[string]struct{}) {
		members = append(members, member)
	}

	return members
}

func (s *Store) SIsMember(key, member string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, exists := s.data[key]
	if !exists {
		return false
	}

	_, ok := set.Value.(map[string]struct{})[member]
	return ok
}

type SortedSet struct {
	Score  float64
	Member string
}

// Sorted Set
func (s *Store) ZAdd(key string, members []SortedSet) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, exists := s.data[key]
	if !exists {
		data = Data{
			Value: make([]SortedSet, 0),
			TTL:   0,
		}
	}

	appended := 0
	zset := data.Value.([]SortedSet)
	for _, member := range members {
		found := false
		for _, existing := range zset {
			if existing.Member == member.Member && existing.Score == member.Score {
				found = true // Member already exists with the same score
				break
			}
		}
		if !found {
			zset = append(zset, member) // Add new member
			appended++
		}
	}
	// Sort the zset by score
	sort.Slice(zset, func(i, j int) bool {
		return zset[i].Score < zset[j].Score
	})

	data.Value = zset
	s.data[key] = data

	return appended
}

func (s *Store) ZRange(key string, start, end int) []SortedSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, exists := s.data[key]; !exists {
		return nil
	}

	zset := s.data[key].Value.([]SortedSet)
	if start < 0 {
		start = len(zset) + start
	}
	if end < 0 {
		end = len(zset) + end
	}
	if start < 0 {
		start = 0
	}
	if end >= len(zset) {
		end = len(zset) - 1
	}
	return zset[start : end+1]
}

func (s *Store) ZRem(key string, members []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; !exists {
		return 0
	}

	zset := s.data[key].Value.([]SortedSet)
	count := 0
	for _, member := range members {
		for i := 0; i < len(zset); i++ {
			if zset[i].Member == member {
				zset = append(zset[:i], zset[i+1:]...)
				count++
				break
			}
		}
	}
	if len(zset) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = Data{
			Value: zset,
			TTL:   0,
		}
	}
	return count
}

// Clear the store of all expired keys
func (s *Store) ClearExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, data := range s.data {
		if data.TTL > 0 && data.TTL < time.Now().Unix() {
			delete(s.data, key)
		}
	}
}
