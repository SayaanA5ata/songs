package main

import (
	"dinushc/gorutines/configs"
	"dinushc/gorutines/internal/song"
	"dinushc/gorutines/pkg/db"
	"dinushc/gorutines/pkg/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	conf := configs.LoadConfig()
	dbase := db.NewDb(conf)
	// Создаем Chi-роутер
	router := chi.NewRouter()
	// Middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)
	//repositories
	songRepository := song.NewSongRepository(dbase)
	//handler
	song.NewSongHandler(router, song.SongHandlerDeps{
		SongRepo: songRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	server.ListenAndServe()
}
