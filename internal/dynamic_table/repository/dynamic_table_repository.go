package repository

import "github.com/gin-gonic/gin"

type DynamicTableRepository interface {
	CreateRecord(ctx *gin.Context, values []string, columns []string) (map[string]interface{}, error)
	ListRecord(ctx *gin.Context, records []map[string]interface{}) ([]map[string]interface{}, error)
	GetRecord(ctx *gin.Context, record map[string]interface{}) (map[string]interface{}, error)
	DeleteRecord(ctx *gin.Context) error
	SoftDeleteRecord(ctx *gin.Context) error
	UpdateRecord(ctx *gin.Context, record map[string]interface{}) (map[string]interface{}, error)
}
