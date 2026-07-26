package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgres://streamly:streamly@localhost:5432/streamly"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := pool.Ping(r.Context())
		if err != nil {
			http.Error(w, "db unreachable", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "ok")
	})

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}