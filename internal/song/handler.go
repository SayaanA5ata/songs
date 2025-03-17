package song

import (
	"dinushc/gorutines/pkg/req"
	"dinushc/gorutines/pkg/res"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type SongHandler struct {
	SongRepo *SongRepository
}

type SongHandlerDeps struct {
	SongRepo *SongRepository
}

func NewSongHandler(router *chi.Mux, deps SongHandlerDeps) {
	handler := &SongHandler{
		SongRepo: deps.SongRepo,
	}
	/*
		router.HandleFunc("POST /song", handler.Create())
		router.HandleFunc("PATCH /song/{id}", handler.Update())
		router.HandleFunc("DELETE /song/{id}", handler.Delete())
		router.HandleFunc("GET /{alias}", handler.GoTo())*/
	// Роуты с использованием Chi
	router.Post("/song", handler.Create())        // Создание песни
	router.Patch("/song/{id}", handler.Update())  // Обновление песни
	router.Delete("/song/{id}", handler.Delete()) // Удаление песни
	router.Get("/{alias}", handler.GoTo())        // Переход по алиасу
}

func (handler *SongHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SongCreateRequest](&w, r)
		if err != nil {
			return
		}

		song := NewSong(body)
		for {
			existedSong, _ := handler.SongRepo.GetByHash(song.Hash)
			if existedSong == nil {
				break
			}
			song.GenerateHash()
		}
		createdSong, err := handler.SongRepo.Create(song)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdSong, 201)
	}
}

func (handler *SongHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SongUpdateRequest](&w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		song, err := handler.SongRepo.Update(&SongModel{
			Model: gorm.Model{ID: uint(id)},
			Group: body.Group,
			Name:  body.Name,
			Date:  body.Date,
			Text:  body.Text,
			Link:  body.Link,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, song, 201)
	}
}

func (handler *SongHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = handler.SongRepo.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.SongRepo.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 200)
	}
}

func (handler *SongHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		song, err := handler.SongRepo.GetByHash(alias)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, song.Link, http.StatusTemporaryRedirect)
	}
}
