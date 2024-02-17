package service

import (
	"fmt"

	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DynamicTableServiceImpl struct {
	DynamicTableRepository repository.DynamicTableRepository
	Validate               *validator.Validate
}

func NewDynamicTableService(dynamicTableRepository repository.DynamicTableRepository, validate *validator.Validate) DynamicTableService {
	return &DynamicTableServiceImpl{
		DynamicTableRepository: dynamicTableRepository,
		Validate:               validate,
	}
}

// CreateRecord implements DynamicTableService.
func (service *DynamicTableServiceImpl) CreateRecord(ctx *gin.Context, result interface{}) (requestData map[string]interface{}, err error) {
	requestData = result.(map[string]interface{})

	var columns, values []string
	for key, value := range requestData {
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("'%v'", value))
	}

	resultData, err := service.DynamicTableRepository.CreateRecord(ctx, columns, values)

	return resultData, err
}
