package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	pool           *pgxpool.Pool
	sessionManager *scs.SessionManager
}

func main() {
	dbURL := "postgres://streamly:streamly@localhost:5432/streamly"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)

	app := &application{
		pool:           pool,
		sessionManager: sessionManager,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/put", app.putHandler)
	mux.HandleFunc("/get", app.getHandler)
	mux.HandleFunc("/health", app.healthHandler)

	log.Printf("listening on :8080")
	if err := http.ListenAndServe(":8080", app.sessionManager.LoadAndSave(mux)); err != nil {
		log.Fatal(err)
	}
}

func (app *application) putHandler(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Put(r.Context(), "message", "Hello from a session!")
}

func (app *application) getHandler(w http.ResponseWriter, r *http.Request) {
	msg := app.sessionManager.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	err := app.pool.Ping(r.Context())
	if err != nil {
		http.Error(w, "db unreachable", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "ok")
}
