package postgres

import (
	"context"
	"strings"
	"testing"
	"time"

	domain "github.com/heartmarshall/digital-forest/backend/internal/domain/plant"
	"github.com/heartmarshall/digital-forest/backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlantRepo_Create(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	repo := NewPlantRepo(dbPool)
	ctx := context.Background()

	tests := []struct {
		name          string
		plant         domain.Plant
		expectedError bool
	}{
		{
			name: "successful plant creation",
			plant: domain.Plant{
				Author:    "test_author",
				ImageData: "base64_image_data",
				CreatedAt: time.Now().UTC(),
			},
			expectedError: false,
		},
		{
			name: "plant with long author name",
			plant: domain.Plant{
				Author:    "very_long_author_name_that_might_cause_issues",
				ImageData: "base64_image_data",
				CreatedAt: time.Now().UTC(),
			},
			expectedError: false,
		},
		{
			name: "plant with empty author",
			plant: domain.Plant{
				Author:    "",
				ImageData: "base64_image_data",
				CreatedAt: time.Now().UTC(),
			},
			expectedError: false, // Database allows empty strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := repo.Create(ctx, tt.plant)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, result.ID)
				assert.Equal(t, tt.plant.Author, result.Author)
				assert.Equal(t, tt.plant.ImageData, result.ImageData)
				assert.WithinDuration(t, tt.plant.CreatedAt, result.CreatedAt, time.Second)
			}
		})
	}
}

func TestPlantRepo_GetRandom(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	repo := NewPlantRepo(dbPool)
	ctx := context.Background()

	// Insert test data
	testPlants := []domain.Plant{
		{Author: "author1", ImageData: "data1", CreatedAt: time.Now().UTC()},
		{Author: "author2", ImageData: "data2", CreatedAt: time.Now().UTC()},
		{Author: "author3", ImageData: "data3", CreatedAt: time.Now().UTC()},
		{Author: "author4", ImageData: "data4", CreatedAt: time.Now().UTC()},
		{Author: "author5", ImageData: "data5", CreatedAt: time.Now().UTC()},
	}

	for _, plant := range testPlants {
		_, err := repo.Create(ctx, plant)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		count          int
		expectedMinLen int
		expectedMaxLen int
	}{
		{
			name:           "get 3 random plants",
			count:          3,
			expectedMinLen: 3,
			expectedMaxLen: 3,
		},
		{
			name:           "get 10 random plants (more than available)",
			count:          10,
			expectedMinLen: 5, // Should return all available plants
			expectedMaxLen: 5,
		},
		{
			name:           "get 0 random plants",
			count:          0,
			expectedMinLen: 0,
			expectedMaxLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := repo.GetRandom(ctx, tt.count)

			// Assert
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(result), tt.expectedMinLen)
			assert.LessOrEqual(t, len(result), tt.expectedMaxLen)

			// Verify that all returned plants have valid data
			for _, plant := range result {
				assert.NotZero(t, plant.ID)
				assert.NotEmpty(t, plant.Author)
				assert.NotEmpty(t, plant.ImageData)
				assert.NotZero(t, plant.CreatedAt)
			}
		})
	}
}

func TestPlantRepo_Integration(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	repo := NewPlantRepo(dbPool)
	ctx := context.Background()

	t.Run("create and retrieve plants", func(t *testing.T) {
		// Create multiple plants
		createdPlants := make([]domain.Plant, 0, 3)
		for i := 0; i < 3; i++ {
			plant := domain.Plant{
				Author:    "author" + string(rune('1'+i)),
				ImageData: "data" + string(rune('1'+i)),
				CreatedAt: time.Now().UTC(),
			}
			created, err := repo.Create(ctx, plant)
			require.NoError(t, err)
			createdPlants = append(createdPlants, created)
		}

		// Retrieve random plants
		randomPlants, err := repo.GetRandom(ctx, 2)
		require.NoError(t, err)
		assert.Len(t, randomPlants, 2)

		// Verify that retrieved plants are among created ones
		createdIDs := make(map[int]bool)
		for _, plant := range createdPlants {
			createdIDs[plant.ID] = true
		}

		for _, plant := range randomPlants {
			assert.True(t, createdIDs[plant.ID], "Retrieved plant ID %d was not among created plants", plant.ID)
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		// Test concurrent creation
		done := make(chan bool, 5)
		errors := make(chan error, 5)

		for i := 0; i < 5; i++ {
			go func(i int) {
				plant := domain.Plant{
					Author:    "concurrent_author_" + string(rune('1'+i)),
					ImageData: "concurrent_data_" + string(rune('1'+i)),
					CreatedAt: time.Now().UTC(),
				}
				_, err := repo.Create(ctx, plant)
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			select {
			case <-done:
				// Success
			case err := <-errors:
				t.Fatalf("Concurrent operation failed: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatal("Concurrent operations timed out")
			}
		}

		// Verify all plants were created
		plants, err := repo.GetRandom(ctx, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(plants), 5)
	})
}

func TestPlantRepo_EdgeCases(t *testing.T) {
	// Arrange
	dbPool, _, container := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, dbPool, container)

	repo := NewPlantRepo(dbPool)
	ctx := context.Background()

	t.Run("get random from empty table", func(t *testing.T) {
		// Ensure table is empty
		err := testutil.TruncateTables(ctx, dbPool)
		require.NoError(t, err)

		// Try to get random plants
		plants, err := repo.GetRandom(ctx, 5)
		assert.NoError(t, err)
		assert.Empty(t, plants)
	})

	t.Run("create plant with special characters", func(t *testing.T) {
		plant := domain.Plant{
			Author:    "author with spaces and symbols!@#$%",
			ImageData: "data with\nnewlines\tand\ttabs",
			CreatedAt: time.Now().UTC(),
		}

		result, err := repo.Create(ctx, plant)
		assert.NoError(t, err)
		assert.Equal(t, plant.Author, result.Author)
		assert.Equal(t, plant.ImageData, result.ImageData)
	})

	t.Run("create plant with very long data", func(t *testing.T) {
		// Create a long string with valid UTF-8 characters
		longString := strings.Repeat("a", 10000) // 10KB string with valid UTF-8
		plant := domain.Plant{
			Author:    "author",
			ImageData: longString,
			CreatedAt: time.Now().UTC(),
		}

		result, err := repo.Create(ctx, plant)
		assert.NoError(t, err)
		assert.Equal(t, plant.ImageData, result.ImageData)
	})
}
