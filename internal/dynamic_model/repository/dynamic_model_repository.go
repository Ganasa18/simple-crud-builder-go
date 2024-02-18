package repository

import (
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/domain"
	"github.com/gin-gonic/gin"
)

type DynamicModelRepository interface {
	CreateModel(ctx *gin.Context, model domain.Model) (domain.Model, error)
	ListModel(ctx *gin.Context) ([]domain.Model, error)
	CreateTable(ctx *gin.Context, model domain.Model)
}
