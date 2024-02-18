package controller

import "github.com/gin-gonic/gin"

type DynamicTableController interface {
	CreateRecord(ctx *gin.Context)
	ListRecord(ctx *gin.Context)
}
