package get_random

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/transport/http/dto"
)

const defaultRandomCount = 15
const maxRandomCount = 50

// GetRandomUseCase - интерфейс для use case получения случайных растений.
type GetRandomUseCase interface {
	GetRandom(ctx context.Context, count int) ([]domain.Plant, error)
}

// GetRandomHandler - HTTP обработчик для получения случайных растений.
type GetRandomHandler struct {
	uc GetRandomUseCase
}

// NewGetRandomHandler - конструктор для хендлера.
func NewGetRandomHandler(uc GetRandomUseCase) *GetRandomHandler {
	return &GetRandomHandler{
		uc: uc,
	}
}

// GetRandomPlants - обработчик для GET /v1/plants/random
func (h *GetRandomHandler) GetRandomPlants(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count := defaultRandomCount

	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid count parameter. Must be a positive integer",
			})
			return
		}

		if count > maxRandomCount {
			count = maxRandomCount
		}
	}

	plants, err := h.uc.GetRandom(r.Context(), count)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get random plants"})
		return
	}

	responses := make([]dto.PlantResponse, len(plants))
	for i, plant := range plants {
		responses[i] = dto.ToPlantResponse(plant)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"plants": responses,
		"count":  len(responses),
	})
}

// respondJSON - хелпер для отправки JSON-ответов.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
