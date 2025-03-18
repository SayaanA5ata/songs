package handlers

import (
	"dinushc/gorutines/internal/payload"
	"dinushc/gorutines/internal/service"
	"dinushc/gorutines/pkg/req"
	"dinushc/gorutines/pkg/res"
	"log"
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

	log.Println("Registering song handlers...")
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
		log.Println("Handling GET /songs request...")

		// Получаем параметры пагинации
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if pageSize == 0 {
			pageSize = 10
		}
		log.Printf("Parsed pagination parameters: page=%d, pageSize=%d", page, pageSize)

		// Формируем фильтр
		filter := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if key != "page" && key != "pageSize" && len(values) > 0 {
				filter[key] = values[0]
				log.Printf("Applied filter: %s=%s", key, values[0])
			}
		}

		// Вызываем сервисный слой
		songs, total, err := handler.Service.GetSongs(filter, page, pageSize)
		if err != nil {
			log.Printf("Error fetching songs from service: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully fetched %d songs from service", len(songs))

		res.Json(w, map[string]interface{}{
			"data":       songs,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		}, http.StatusOK)
		log.Printf("Response sent with %d songs", len(songs))
	}
}

// Получение текста песни с пагинацией по куплетам
func (handler *SongHandler) GetSongVerses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling GET /songs/{id}/verses request...")

		// Получаем ID песни
		idString := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Printf("Invalid song ID: %s", idString)
			http.Error(w, "Invalid song ID", http.StatusBadRequest)
			return
		}
		log.Printf("Parsed song ID: %d", id)

		// Получаем параметры пагинации
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if pageSize == 0 {
			pageSize = 5
		}
		log.Printf("Parsed pagination parameters: page=%d, pageSize=%d", page, pageSize)

		// Вызываем сервисный слой
		verses, totalVerses, err := handler.Service.GetSongVerses(uint(id), page, pageSize)
		if err != nil {
			log.Printf("Error fetching verses from service: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully fetched %d verses from service", len(verses))

		res.Json(w, map[string]interface{}{
			"data":       verses,
			"total":      totalVerses,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (totalVerses + pageSize - 1) / pageSize,
		}, http.StatusOK)
		log.Printf("Response sent with %d verses", len(verses))
	}
}

func (handler *SongHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling POST /song request...")

		body, err := req.HandleBody[payload.SongCreateRequest](&w, r)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return
		}
		log.Printf("Parsed request body: %+v", body)

		song, err := handler.Service.CreateSong(body)
		if err != nil {
			log.Printf("Error creating song in service: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Song created successfully with ID: %d", song.ID)

		res.Json(w, song, 201)
		log.Println("Response sent with created song")
	}
}

func (handler *SongHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling PATCH /song/{id} request...")

		body, err := req.HandleBody[payload.SongUpdateRequest](&w, r)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return
		}
		log.Printf("Parsed request body: %+v", body)

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			log.Printf("Invalid song ID: %s", idString)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Parsed song ID: %d", id)

		song, err := handler.Service.UpdateSong(uint(id), body)
		if err != nil {
			log.Printf("Error updating song in service: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Song updated successfully with ID: %d", song.ID)

		res.Json(w, song, 201)
		log.Println("Response sent with updated song")
	}
}

func (handler *SongHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling DELETE /song/{id} request...")

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			log.Printf("Invalid song ID: %s", idString)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Parsed song ID: %d", id)

		err = handler.Service.DeleteSong(uint(id))
		if err != nil {
			log.Printf("Error deleting song in service: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Song deleted successfully with ID: %d", id)

		res.Json(w, nil, 200)
		log.Println("Response sent for successful deletion")
	}
}

func (handler *SongHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling GET /{alias} request...")

		alias := r.PathValue("alias")
		log.Printf("Parsed alias: %s", alias)

		song, err := handler.Service.GetSongByHash(alias)
		if err != nil {
			log.Printf("Error fetching song by hash: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Redirecting to song link: %s", song.Link)

		http.Redirect(w, r, song.Link, http.StatusTemporaryRedirect)
		log.Println("Redirect completed successfully")
	}
}
