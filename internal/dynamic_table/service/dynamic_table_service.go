package service

import "github.com/gin-gonic/gin"

type DynamicTableService interface {
	CreateRecord(ctx *gin.Context, request interface{}) (map[string]interface{}, error)
	ListRecord(ctx *gin.Context) ([]map[string]interface{}, error)
}
