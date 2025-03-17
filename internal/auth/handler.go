package auth

import (
	"dinushc/gorutines/configs"
	"dinushc/gorutines/pkg/req"
	"dinushc/gorutines/pkg/res"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *chi.Mux, deps AuthHandlerDeps) {
	ah := &AuthHandler{
		Config: deps.Config,
	}
	router.Post("/auth/login", ah.Login())
	router.Post("/auth/register", ah.Register())
}

func (ah *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		fmt.Println(body)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (ah *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		fmt.Println(body)
	}
}
