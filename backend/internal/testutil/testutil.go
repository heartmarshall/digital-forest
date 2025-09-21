package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDBConfig содержит конфигурацию для тестовой базы данных
type TestDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
}

// SetupTestDB создает тестовую базу данных с помощью testcontainers
func SetupTestDB(t *testing.T) (*pgxpool.Pool, *TestDBConfig, testcontainers.Container) {
	ctx := context.Background()

	// Создаем PostgreSQL контейнер
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %v", err)
	}

	// Получаем строку подключения
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Создаем пул соединений
	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Проверяем соединение
	if err := dbPool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Создаем таблицы
	if err := createTestTables(ctx, dbPool); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Извлекаем конфигурацию из строки подключения
	config := &TestDBConfig{
		DSN: connStr,
	}

	return dbPool, config, postgresContainer
}

// CleanupTestDB закрывает соединения и останавливает контейнер
func CleanupTestDB(t *testing.T, dbPool *pgxpool.Pool, container testcontainers.Container) {
	if dbPool != nil {
		dbPool.Close()
	}
	if container != nil {
		ctx := context.Background()
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// createTestTables создает необходимые таблицы для тестов
func createTestTables(ctx context.Context, db *pgxpool.Pool) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS plants (
		id SERIAL PRIMARY KEY,
		author VARCHAR(255) NOT NULL,
		image_data TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(ctx, createTableSQL)
	return err
}

// TruncateTables очищает все таблицы для изоляции тестов
func TruncateTables(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, "TRUNCATE TABLE plants RESTART IDENTITY CASCADE")
	return err
}
