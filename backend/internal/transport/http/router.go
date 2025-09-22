package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	createHandler "github.com/heartmarshall/digital-forest/backend/internal/transport/http/handlers/plant/create"
	getRandomHandler "github.com/heartmarshall/digital-forest/backend/internal/transport/http/handlers/plant/get_random"
	createUseCase "github.com/heartmarshall/digital-forest/backend/internal/usecase/plant/create"
	getRandomUseCase "github.com/heartmarshall/digital-forest/backend/internal/usecase/plant/get_random"
)

// NewRouter создает новый роутер, регистрирует все маршруты и middleware.
func NewRouter(createUC *createUseCase.CreateUseCase, getRandomUC *getRandomUseCase.GetRandomUseCase) http.Handler {
	// Создаем экземпляр валидатора
	validator := NewValidator()

	// Создаем handlers для каждого use case
	createHandlerInstance := createHandler.NewCreateHandler(createUC, validator)
	getRandomHandlerInstance := getRandomHandler.NewGetRandomHandler(getRandomUC)

	router := chi.NewRouter()

	// Настройка Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Настройка CORS для локальной разработки
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Группа роутов для нашего API v1
	router.Route("/v1", func(r chi.Router) {
		r.Post("/plants", createHandlerInstance.CreatePlant)
		r.Get("/plants/random", getRandomHandlerInstance.GetRandomPlants)
	})

	return router
}
