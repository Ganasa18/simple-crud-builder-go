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

func (repository *DynamicTableRepositoryImpl) ListRecord(ctx *gin.Context, records []map[string]interface{}) ([]map[string]interface{}, error) {
	tableName := ctx.Param("table")
	query := fmt.Sprintf("SELECT * FROM %s WHERE deleted_at IS NULL", tableName)

	rows, err := repository.DB.Raw(query).Rows()

	if err != nil {
		return records, errors.New("error querying table get records ")
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return records, errors.New("error querying table get rows column")
	}

	values := make([]interface{}, len(columns))
	valueKey := make([]interface{}, len(columns))
	for i := range columns {
		valueKey[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valueKey...); err != nil {
			return records, errors.New("error querying table scan column")
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			rowData[col] = val
		}

		records = append(records, rowData)
	}

	return records, nil
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
