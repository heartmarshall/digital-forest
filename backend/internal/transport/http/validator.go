package http

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator - это обертка над стандартным валидатором.
type Validator struct {
	validate *validator.Validate
}

// NewValidator создает новый экземпляр валидатора.
func NewValidator() *Validator {
	return &Validator{validate: validator.New()}
}

// ValidateStruct выполняет валидацию переданной структуры.
// В случае ошибки возвращает map с удобным для JSON ответа описанием.
func (v *Validator) ValidateStruct(s interface{}) map[string]string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	// Преобразуем ошибку в ValidationErrors для детального анализа
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Если это не ошибка валидации, возвращаем как есть
		return map[string]string{"error": "internal validation error"}
	}

	// Создаем map для хранения ошибок
	errorMessages := make(map[string]string)

	for _, fieldErr := range validationErrors {
		// fieldErr.Field() - имя поля, fieldErr.Tag() - правило, которое не сработало
		fieldName := strings.ToLower(fieldErr.Field())
		switch fieldErr.Tag() {
		case "required":
			errorMessages[fieldName] = fmt.Sprintf("field '%s' is required", fieldName)
		case "max":
			errorMessages[fieldName] = fmt.Sprintf("field '%s' is too long (max: %s)", fieldName, fieldErr.Param())
		default:
			errorMessages[fieldName] = fmt.Sprintf("field '%s' is not valid", fieldName)
		}
	}

	return errorMessages
}
