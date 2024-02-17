package repository

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/domain"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DynamicModelRepositoryImpl struct {
	DB *gorm.DB
}

func NewDynamicModelRepository(db *gorm.DB) DynamicModelRepository {
	return &DynamicModelRepositoryImpl{
		DB: db,
	}
}

// CreateModel implements DynamicModelRepository.
func (repository *DynamicModelRepositoryImpl) CreateModel(ctx *gin.Context, model domain.Model) (domain.Model, error) {
	err := repository.DB.Create(&model).Error
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repository *DynamicModelRepositoryImpl) CreateTable(ctx *gin.Context, model domain.Model) {
	// Create a dynamic struct
	fields := []reflect.StructField{
		{Name: "Id", Type: reflect.TypeOf(0), Tag: reflect.StructTag(`gorm:"primaryKey"`)},
		{Name: "CreatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:current_timestamp"`)},
		{Name: "UpdatedAt", Type: reflect.TypeOf(time.Time{}), Tag: reflect.StructTag(`gorm:"default:current_timestamp"`)},
		{Name: "DeletedAt", Type: reflect.TypeOf(gorm.DeletedAt{})},
	}
	for _, field := range model.Fields {
		var fieldType reflect.Type

		switch field.Type {
		case "number":
			fieldType = reflect.TypeOf(0)
		case "string":
			fieldType = reflect.TypeOf("")
		case "boolean":
			fieldType = reflect.TypeOf(true)
		default:
			// Handle unsupported types
			panic(fmt.Sprintf("Unsupported type for field %s", helper.CamelToSnake(field.Name)))
		}

		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`gorm:"column:%s"`, helper.CamelToSnake(field.Name))),
		})
	}

	resultType := reflect.StructOf(fields)

	// Create an instance of the dynamic struct
	instance := reflect.New(resultType).Interface()

	tableName := model.Name
	repository.DB.Table(tableName).AutoMigrate(instance)
}
