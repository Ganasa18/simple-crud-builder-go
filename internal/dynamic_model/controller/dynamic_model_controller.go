package controller

import (
	"github.com/gin-gonic/gin"
)

type DynamicModelController interface {
	CreateModel(ctx *gin.Context)
}
