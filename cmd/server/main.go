package main

import (
	"github.com/teguhkurnia/redis-like/internal/server"
	"github.com/teguhkurnia/redis-like/internal/store"
)

func main() {
	store := store.NewStore()
	server := server.NewServer(":8080", store)
	server.Start()
}
