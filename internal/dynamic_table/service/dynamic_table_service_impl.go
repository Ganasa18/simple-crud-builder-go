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

func (service *DynamicTableServiceImpl) ListRecord(ctx *gin.Context) ([]map[string]interface{}, error) {

	var records []map[string]interface{}
	recordResponse, err := service.DynamicTableRepository.ListRecord(ctx, records)

	return recordResponse, err
}

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

func (service *DynamicTableServiceImpl) GetRecord(ctx *gin.Context) (map[string]interface{}, error) {
	record := make(map[string]interface{})

	resultData, err := service.DynamicTableRepository.GetRecord(ctx, record)

	return resultData, err

}

func (service *DynamicTableServiceImpl) DeleteRecord(ctx *gin.Context) error {
	err := service.DynamicTableRepository.DeleteRecord(ctx)
	return err
}

func (service *DynamicTableServiceImpl) SoftDeleteRecord(ctx *gin.Context) error {
	err := service.DynamicTableRepository.SoftDeleteRecord(ctx)
	return err
}

func (service *DynamicTableServiceImpl) UpdateRecord(ctx *gin.Context, record map[string]interface{}) (map[string]interface{}, error) {
	result, err := service.DynamicTableRepository.UpdateRecord(ctx, record)
	return result, err
}
