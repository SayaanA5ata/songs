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

func NewSongHandler(router *chi.Mux, service *service.SongService) {
	handler := &SongHandler{
		Service: service,
	}

	router.Get("/songs", handler.GetSongs())
	router.Get("/songs/{id}/verses", handler.GetSongVerses())
	router.Post("/song", handler.Create())
	router.Patch("/song/{id}", handler.Update())
	router.Delete("/song/{id}", handler.Delete())
	router.Get("/{alias}", handler.GoTo())
}

// Получение данных библиотеки с фильтрацией и пагинацией
func (handler *SongHandler) GetSongs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметры пагинации
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if pageSize == 0 {
			pageSize = 10
		}

		// Формируем фильтр
		filter := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if key != "page" && key != "pageSize" && len(values) > 0 {
				filter[key] = values[0]
			}
		}

		// Вызываем сервисный слой
		songs, total, err := handler.Service.GetSongs(filter, page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]interface{}{
			"data":       songs,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		}, http.StatusOK)
	}
}

// Получение текста песни с пагинацией по куплетам
func (handler *SongHandler) GetSongVerses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID песни
		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}

		// Получаем параметры пагинации
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if pageSize == 0 {
			pageSize = 5
		}

		// Вызываем сервисный слой
		verses, totalVerses, err := handler.Service.GetSongVerses(uint(id), page, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, map[string]interface{}{
			"data":       verses,
			"total":      totalVerses,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (totalVerses + pageSize - 1) / pageSize,
		}, http.StatusOK)
	}
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
