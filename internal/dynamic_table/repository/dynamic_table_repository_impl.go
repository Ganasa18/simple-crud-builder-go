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

func (repository *DynamicTableRepositoryImpl) CreateRecord(ctx *gin.Context, columns []string, values []string) (map[string]interface{}, error) {
	tableName := ctx.Param("table")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *", tableName, strings.Join(columns, ", "), strings.Join(values, ", "))

	var requestData map[string]interface{}

	result := repository.DB.Raw(query).Scan(&requestData)

	if result.Error != nil {
		return nil, errors.New("error inserting data: " + result.Error.Error())
	}

	return requestData, nil
}
