package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

func (repository *DynamicTableRepositoryImpl) GetRecord(ctx *gin.Context, record map[string]interface{}) (map[string]interface{}, error) {
	tableName := ctx.Param("table")
	idRecord := ctx.Param("id")

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	err := repository.DB.Raw(query, idRecord).Scan(&record).Error
	if err != nil {
		return record, errors.New("errot to get record: " + err.Error())
	}

	if len(record) == 0 {
		return record, errors.New("record not found")
	}

	return record, nil

}

func (repository *DynamicTableRepositoryImpl) DeleteRecord(ctx *gin.Context) error {
	tableName := ctx.Param("table")
	idRecord := ctx.Param("id")
	// Check if the record exists
	if err := repository.recordExists(tableName, idRecord); err != nil {
		return err
	}

	// Delete the record
	if err := repository.DB.Table(tableName).Where("id = ?", idRecord).Delete(nil).Error; err != nil {
		return errors.New("error deleting record: " + err.Error())
	}

	return nil
}

// SoftDeleteRecord implements DynamicTableRepository.
func (repository *DynamicTableRepositoryImpl) SoftDeleteRecord(ctx *gin.Context) error {
	tableName := ctx.Param("table")
	idRecord := ctx.Param("id")

	// Check if the record exists
	err := repository.recordExists(tableName, idRecord)
	if err != nil {
		return err
	}

	// Perform the soft delete operation
	err = repository.DB.Table(tableName).Where("id = ?", idRecord).Updates(map[string]interface{}{"deleted_at": time.Now()}).Error
	if err != nil {
		return errors.New("error soft deleting record: " + err.Error())
	}
	return nil
}

// UpdateRecord implements DynamicTableRepository.
func (repository *DynamicTableRepositoryImpl) UpdateRecord(ctx *gin.Context, record map[string]interface{}) (map[string]interface{}, error) {

	tableName := ctx.Param("table")
	idRecord := ctx.Param("id")

	var emptyRecord map[string]interface{}

	// Check if the record exists
	err := repository.recordExists(tableName, idRecord)
	if err != nil {
		return emptyRecord, err
	}

	record["updated_at"] = time.Now()

	// Perform the update operation
	err = repository.DB.Table(tableName).Where("id = ?", idRecord).Updates(record).Error

	if err != nil {
		return emptyRecord, errors.New("error to update record: " + err.Error())
	}

	// Fetch the updated record from the database as a map
	err = repository.DB.Table(tableName).Where("id = ?", idRecord).Scan(&record).Error

	if err != nil {
		return emptyRecord, errors.New("error to get record: " + err.Error())
	}

	return record, nil
}

// recordExists checks if a record exists in the specified table.
func (repository *DynamicTableRepositoryImpl) recordExists(tableName, idRecord string) error {
	var count int64
	if err := repository.DB.Table(tableName).Where("id = ?", idRecord).Count(&count).Error; err != nil {
		return errors.New("error finding record: " + err.Error())
	}

	if count == 0 {
		return errors.New("record not found")
	}

	return nil
}
