package store

import (
	"fmt"
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

func (s *Store) SetWithTTL(key string, value interface{}, ttl int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = Data{
		Value: value,
		TTL:   ttl,
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
	fmt.Println("List for key:", key, "is:", list)
	listLen := len(list)
	fmt.Println("List length for key:", key, "is:", listLen, "start:", start, "end:", end)

	// handle the negative indices
	if start < 0 {
		start = listLen + start
	}
	if end < 0 {
		end = listLen + end
	}

	fmt.Println("LRange command received for key:", key, "start:", start, "end:", end)

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
