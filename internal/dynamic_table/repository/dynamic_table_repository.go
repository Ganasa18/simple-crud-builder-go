package repository

import "github.com/gin-gonic/gin"

type DynamicTableRepository interface {
	CreateRecord(ctx *gin.Context, values []string, columns []string) (map[string]interface{}, error)
	ListRecord(ctx *gin.Context, records []map[string]interface{}) ([]map[string]interface{}, error)
}
