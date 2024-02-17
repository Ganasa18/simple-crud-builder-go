package service

import (
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/web"
	"github.com/gin-gonic/gin"
)

type DynamicModelService interface {
	CreateModel(ctx *gin.Context, request web.ModelRequest) (web.ModelResponse, error)
}
