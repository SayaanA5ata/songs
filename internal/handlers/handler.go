package handlers

import (
	"dinushc/gorutines/internal/payload"
	"dinushc/gorutines/internal/service"
	"dinushc/gorutines/pkg/req"
	"dinushc/gorutines/pkg/res"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SongHandler struct {
	Service *service.SongService
}

type SongHandlerDeps struct {
	Service *service.SongService
}

func NewSongHandler(router *chi.Mux, deps SongHandlerDeps) {
	handler := &SongHandler{
		Service: deps.Service,
	}
	router.Post("/song", handler.Create())
	router.Patch("/song/{id}", handler.Update())
	router.Delete("/song/{id}", handler.Delete())
	router.Get("/{alias}", handler.GoTo())
}

func (handler *SongHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.SongCreateRequest](&w, r)
		if err != nil {
			return
		}
		song, err := handler.Service.CreateSong(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, song, 201)
	}
}

func (handler *SongHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[payload.SongUpdateRequest](&w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		song, err := handler.Service.UpdateSong(uint(id), body)
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
		err = handler.Service.DeleteSong(uint(id))
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
		song, err := handler.Service.GetSongByHash(alias)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, song.Link, http.StatusTemporaryRedirect)
	}
}
