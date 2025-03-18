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

	router.Get("/songs", handler.GetSongs())
	router.Get("/songs/{id}/verses", handler.GetSongVerses())
	router.Post("/song", handler.Create())
	router.Patch("/song/{id}", handler.Update())
	router.Delete("/song/{id}", handler.Delete())
	router.Get("/{alias}", handler.GoTo())
}

// @Summary Получение данных библиотеки с фильтрацией и пагинацией
// @Tags songs
// @Param group query string false "Фильтр по названию группы"
// @Param name query string false "Фильтр по названию песни"
// @Param date query string false "Фильтр по дате выпуска"
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} map[string]interface{} "Успешный ответ с данными песен"
// @Failure 400 {object} map[string]string "Ошибка валидации параметров"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs [get]
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

// @Summary Получение текста песни с пагинацией по куплетам
// @Tags songs
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "Размер страницы"
// @Success 200 {object} map[string]interface{} "Успешный ответ с куплетами песни"
// @Failure 400 {object} map[string]string "Ошибка валидации параметров"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs/{id}/verses [get]
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

// @Summary Добавление новой песни
// @Tags songs
// @Accept json
// @Produce json
// @Param song body payload.SongCreateRequest true "Данные для создания песни"
// @Success 201 {object} domain.SongModel "Успешное создание песни"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /song [post]
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

		res.Json(w, song, 201)
		log.Println("Response sent with created song")
	}
}

// @Summary Обновление данных песни
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body payload.SongUpdateRequest true "Данные для обновления песни"
// @Success 200 {object} domain.SongModel "Успешное обновление песни"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /song/{id} [patch]
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

		res.Json(w, song, 201)
		log.Println("Response sent with updated song")
	}
}

// @Summary Удаление песни
// @Tags songs
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string "Песня успешно удалена"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /song/{id} [delete]
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
		res.Json(w, nil, 200)
	}
}

// @Summary Переход по ссылке песни
// @Tags songs
// @Param alias path string true "Хеш песни"
// @Success 302 {string} string "Редирект на ссылку песни"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /{alias} [get]
func (handler *SongHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling GET /{alias} request...")

		alias := r.PathValue("alias")
		log.Printf("Parsed alias: %s", alias)

		song, err := handler.Service.GetSongByHash(alias)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Redirecting to song link: %s", song.Link)

		http.Redirect(w, r, song.Link, http.StatusTemporaryRedirect)
		log.Println("Redirect completed successfully")
	}
}
