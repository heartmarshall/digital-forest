package plant

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	// Убедись, что путь импорта соответствует имени твоего Go-модуля
	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// PlantRepo - это реализация usecase.PlantRepository для работы с PostgreSQL.
type PlantRepo struct {
	db *pgxpool.Pool
}

// NewPlantRepo - конструктор для репозитория.
// Принимает пул соединений с базой данных в качестве зависимости.
func NewPlantRepo(db *pgxpool.Pool) *PlantRepo {
	return &PlantRepo{db: db}
}

// Create реализует метод интерфейса usecase.PlantRepository.
// Он вставляет новую запись о растении в таблицу "plants".
func (r *PlantRepo) Create(ctx context.Context, plant domain.Plant) (domain.Plant, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("plants").
		Columns("author", "image_data", "created_at").
		Values(plant.Author, plant.ImageData, plant.CreatedAt).
		Suffix("RETURNING id, author, image_data, created_at"). // Возвращаем все поля
		ToSql()
	if err != nil {
		return domain.Plant{}, fmt.Errorf("PlantRepo - Create - ToSql: %w", err)
	}

	var createdPlant domain.Plant
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&createdPlant.ID,
		&createdPlant.Author,
		&createdPlant.ImageData,
		&createdPlant.CreatedAt,
	)
	if err != nil {
		return domain.Plant{}, fmt.Errorf("PlantRepo - Create - QueryRow.Scan: %w", err)
	}
	return createdPlant, nil
}

// GetRandom реализует метод интерфейса usecase.PlantRepository.
// Он извлекает случайные записи из таблицы "plants".
func (r *PlantRepo) GetRandom(ctx context.Context, count int) ([]domain.Plant, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "author", "image_data", "created_at").
		From("plants").
		OrderBy("RANDOM()"). // ORDER BY RANDOM() - простой, но потенциально медленный способ для очень больших таблиц.
		Limit(uint64(count)).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("PlantRepo - GetRandom - ToSql: %w", err)
	}

	// Выполняем запрос для получения нескольких строк.
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PlantRepo - GetRandom - Query: %w", err)
	}
	defer rows.Close()

	plants := make([]domain.Plant, 0, count)

	// Итерируемся по результатам и сканируем каждую строку в структуру domain.Plant.
	for rows.Next() {
		var p domain.Plant
		if err := rows.Scan(&p.ID, &p.Author, &p.ImageData, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("PlantRepo - GetRandom - rows.Scan: %w", err)
		}
		plants = append(plants, p)
	}

	// Проверяем на наличие ошибок, которые могли возникнуть во время итерации.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PlantRepo - GetRandom - rows.Err: %w", err)
	}

	return plants, nil
}
