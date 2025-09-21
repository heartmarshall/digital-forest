package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/handlers"
	"github.com/heartmarshall/digital-forest/backend/internal/usecase"
)

// NewRouter создает новый роутер, регистрирует все маршруты и middleware.
func NewRouter(uc *usecase.PlantUseCase) http.Handler {
	// Создаем экземпляр валидатора
	validator := NewValidator()

	// Передаем валидатор в конструктор хендлера
	handler := handlers.NewPlantHandler(uc, validator)

	router := chi.NewRouter()

	// Настройка Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Настройка CORS для локальной разработки
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Группа роутов для нашего API v1
	router.Route("/v1", func(r chi.Router) {
		r.Post("/plants", handler.CreatePlant)
		r.Get("/plants/random", handler.GetRandomPlants)
	})

	return router
}
