package main

import (
	"time"

	"github.com/teguhkurnia/redis-like/internal/server"
	"github.com/teguhkurnia/redis-like/internal/store"
)

func main() {
	store := store.NewStore()
	server := server.NewServer(":8080", store)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			store.ClearExpired()
		}
	}()

	server.Start()
}
