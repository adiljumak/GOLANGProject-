package main

import (
	"Lecture6/internal/http"
	"Lecture6/internal/store/postgres"
	"context"
	lru "github.com/hashicorp/golang-lru"
	"github.com/jackc/pgx/v4"
)

func main() {
	urlExample := "postgres://postgres:1111@localhost:5432/goods"
	store := postgres.NewDB()
	s := pgx.Conn{}
	s = s
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}

	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	srv := http.NewServer(context.Background(), http.WithAddress(":8080"), http.WithStore(store), http.WithCache(cache))
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTerminator()
}
