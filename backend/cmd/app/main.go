package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heartmarshall/digital-forest/backend/internal/config"
	"github.com/heartmarshall/digital-forest/backend/internal/repository/plant"
	transportHTTP "github.com/heartmarshall/digital-forest/backend/internal/transport/http"
	"github.com/heartmarshall/digital-forest/backend/internal/usecase"
)

func main() {
	// Настраиваем graceful shutdown:
	// Создаем контекст, который будет отменен при получении сигнала SIGINT или SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 1. Инициализация конфигурации
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Подключение к базе данных
	// Собираем DSN (Data Source Name) из отдельных полей конфигурации.
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Проверяем, что соединение с БД действительно установлено.
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	log.Println("database connection successful")

	// 3. Сборка всех зависимостей (Dependency Injection)
	// Идем "изнутри наружу": Repository -> UseCase -> Handler -> Router
	plantRepo := plant.NewPlantRepo(dbPool)
	plantUseCase := usecase.NewPlantUseCase(plantRepo)
	router := transportHTTP.NewRouter(plantUseCase) // Роутер создается с зависимостью от use case

	// 4. Настройка и запуск HTTP-сервера
	server := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: router,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("starting server on port %s", cfg.HTTP.Port)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("service stopped gracefully")
}
