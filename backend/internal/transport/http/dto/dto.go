package dto

import (
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
)

// CreatePlantRequest - DTO для запроса на создание растения.
// Теги `json:"..."` используются для сериализации/десериализации.
type CreatePlantRequest struct {
	Author    string `json:"author"`
	ImageData string `json:"imageData"`
}

// PlantResponse - DTO для ответа клиенту.
// Мы отделяем эту структуру от доменной, чтобы иметь полный контроль
// над тем, как наши данные выглядят в API.
type PlantResponse struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	ImageData string    `json:"imageData"`
	CreatedAt time.Time `json:"createdAt"`
}

// ToPlantResponse преобразует доменную модель в DTO для ответа.
func ToPlantResponse(p domain.Plant) PlantResponse {
	return PlantResponse{
		ID:        p.ID,
		Author:    p.Author,
		ImageData: p.ImageData,
		CreatedAt: p.CreatedAt,
	}
}
