package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vmpyr/afterlight/internal/store"
)

//go:embed web/dist/*
var dist embed.FS

func main() {
	db, err := store.NewStorage("afterlight.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Afterlight Systems: ONLINE"))
		})
	})

	contentStatic, _ := fs.Sub(dist, "web/dist")
	r.Handle("/*", http.FileServer(http.FS(contentStatic)))

	log.Println("Afterlight running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
