package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vmpyr/afterlight/internal/api"
	"github.com/vmpyr/afterlight/internal/store"
)

//go:embed all:web/dist
var dist embed.FS

func main() {
	storage, err := store.NewStorage("afterlight.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	repo := store.NewStore(storage.DB())

	authHandler := api.NewAuthHandler(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Afterlight Systems: ONLINE"))
		})

		r.Mount("/auth", authHandler.Routes())
	})

	contentStatic, _ := fs.Sub(dist, "web/dist")
	r.Handle("/*", http.FileServer(http.FS(contentStatic)))

	log.Println("Afterlight running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
