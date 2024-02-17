package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DynamicTableRepositoryImpl struct {
	DB *gorm.DB
}

func NewDynamicTableRepository(db *gorm.DB) DynamicTableRepository {
	return &DynamicTableRepositoryImpl{
		DB: db,
	}
}

// CreateRecord implements DynamicTableRepository.
func (repository *DynamicTableRepositoryImpl) CreateRecord(ctx *gin.Context, values []string, columns []string) (requestData map[string]interface{}, err error) {
	tableName := ctx.Param("table")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))
	err = repository.DB.Exec(query).Error
	if err != nil {
		var requestEmptyData map[string]interface{}
		return requestEmptyData, errors.New("error insert data")
	}

	return requestData, nil

}
