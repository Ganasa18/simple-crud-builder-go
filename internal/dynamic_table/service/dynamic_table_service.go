package service

import "github.com/gin-gonic/gin"

type DynamicTableService interface {
	CreateRecord(ctx *gin.Context, request interface{}) (map[string]interface{}, error)
	ListRecord(ctx *gin.Context) ([]map[string]interface{}, error)
	GetRecord(ctx *gin.Context) (map[string]interface{}, error)
	DeleteRecord(ctx *gin.Context) error
	SoftDeleteRecord(ctx *gin.Context) error
	UpdateRecord(ctx *gin.Context, request map[string]interface{}) (map[string]interface{}, error)
}
