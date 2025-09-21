package testutil

import (
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
)

// TestPlants содержит тестовые данные для растений
var TestPlants = []domain.Plant{
	{
		ID:        1,
		Author:    "test_author_1",
		ImageData: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	},
	{
		ID:        2,
		Author:    "test_author_2",
		ImageData: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		CreatedAt: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	},
	{
		ID:        3,
		Author:    "test_author_3",
		ImageData: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		CreatedAt: time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
	},
}

// TestCreatePlantRequests содержит тестовые данные для запросов создания растений
var TestCreatePlantRequests = []dto.CreatePlantRequest{
	{
		Author:    "valid_author",
		ImageData: "valid_base64_image_data",
	},
	{
		Author:    "author_with_special_chars_!@#$%",
		ImageData: "image_data_with_special_chars",
	},
	{
		Author:    "author_with_unicode_测试",
		ImageData: "image_data_with_unicode",
	},
}

// InvalidCreatePlantRequests содержит невалидные тестовые данные
var InvalidCreatePlantRequests = []dto.CreatePlantRequest{
	{
		Author:    "", // Empty author
		ImageData: "valid_data",
	},
	{
		Author:    "valid_author",
		ImageData: "", // Empty image data
	},
	{
		Author:    string(make([]byte, 256)), // Author too long
		ImageData: "valid_data",
	},
}

// GetTestPlant возвращает тестовое растение с указанным ID
func GetTestPlant(id int) domain.Plant {
	for _, plant := range TestPlants {
		if plant.ID == id {
			return plant
		}
	}
	return domain.Plant{}
}

// GetTestCreatePlantRequest возвращает тестовый запрос создания растения
func GetTestCreatePlantRequest(index int) dto.CreatePlantRequest {
	if index >= 0 && index < len(TestCreatePlantRequests) {
		return TestCreatePlantRequests[index]
	}
	return dto.CreatePlantRequest{}
}

// GetInvalidCreatePlantRequest возвращает невалидный тестовый запрос
func GetInvalidCreatePlantRequest(index int) dto.CreatePlantRequest {
	if index >= 0 && index < len(InvalidCreatePlantRequests) {
		return InvalidCreatePlantRequests[index]
	}
	return dto.CreatePlantRequest{}
}
