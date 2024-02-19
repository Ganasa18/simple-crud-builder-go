package controller

import "github.com/gin-gonic/gin"

type DynamicTableController interface {
	CreateRecord(ctx *gin.Context)
	ListRecord(ctx *gin.Context)
	GetRecord(ctx *gin.Context)
	UpdateRecord(ctx *gin.Context)
	DeleteRecord(ctx *gin.Context)
	SoftDeleteRecord(ctx *gin.Context)
}
