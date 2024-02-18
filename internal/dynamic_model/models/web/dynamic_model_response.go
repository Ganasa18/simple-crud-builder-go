package web

import (
	"time"

	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/domain"
	"gorm.io/gorm"
)

type ModelRequest struct {
	Name   string         `json:"name"`
	Fields []domain.Field `json:"fields"`
}
type ModelResponse struct {
	Name      string         `json:"name"`
	Fields    []domain.Field `json:"fields"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func ToDynamicModelResponseWithError(model domain.Model, err error) (ModelResponse, error) {
	var modelResponse = ModelResponse{
		Name:      model.Name,
		Fields:    model.Fields,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}

	return modelResponse, err
}
func ToDynamicModelResponse(model domain.Model) ModelResponse {
	var modelResponse = ModelResponse{
		Name:      model.Name,
		Fields:    model.Fields,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}

	return modelResponse
}

func ToDynamicModelResponses(models []domain.Model, err error) ([]ModelResponse, error) {
	var modelResponse []ModelResponse
	for _, model := range models {
		modelResponse = append(modelResponse, ToDynamicModelResponse(model))
	}
	return modelResponse, err
}
