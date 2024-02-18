package service

import (
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/domain"
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/web"
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/repository"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DynamicModelServiceImpl struct {
	DynamicModelRepository repository.DynamicModelRepository
	Validate               *validator.Validate
}

func NewDynamicModelService(dynamicModelRepository repository.DynamicModelRepository, validate *validator.Validate) DynamicModelService {
	return &DynamicModelServiceImpl{
		DynamicModelRepository: dynamicModelRepository,
		Validate:               validate,
	}
}

// CreateModel implements DynamicModelService.
func (service *DynamicModelServiceImpl) CreateModel(ctx *gin.Context, request web.ModelRequest) (web.ModelResponse, error) {
	err := service.Validate.Struct(request)
	utils.PanicIfError(err)

	// LOGIC
	model := domain.Model{
		Name:   request.Name,
		Fields: request.Fields,
	}

	modelResponse, err := service.DynamicModelRepository.CreateModel(ctx, model)
	if err == nil {
		service.DynamicModelRepository.CreateTable(ctx, modelResponse)
	}

	return web.ToDynamicModelResponseWithError(modelResponse, err)
}

// ListModel implements DynamicModelService.
func (service *DynamicModelServiceImpl) ListModel(ctx *gin.Context) ([]web.ModelResponse, error) {
	roleResponse, err := service.DynamicModelRepository.ListModel(ctx)
	return web.ToDynamicModelResponses(roleResponse, err)
}
