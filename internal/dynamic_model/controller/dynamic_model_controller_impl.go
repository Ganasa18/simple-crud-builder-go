package controller

import (
	"net/http"

	response "github.com/Ganasa18/simple-crud-builder-go/internal/base/models/web"
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/models/web"
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_model/service"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/helper"
	"github.com/gin-gonic/gin"
)

type DynamicModelControllerImpl struct {
	DynamicModelService service.DynamicModelService
}

func NewDynamicModelController(dynamicModelService service.DynamicModelService) DynamicModelController {
	return &DynamicModelControllerImpl{
		DynamicModelService: dynamicModelService,
	}
}

// CreateModel implements DynamicModelController.
func (controller *DynamicModelControllerImpl) CreateModel(ctx *gin.Context) {
	modelRequest := web.ModelRequest{}
	helper.ReadFromRequestBody(ctx.Request, &modelRequest)

	responseModel, err := controller.DynamicModelService.CreateModel(ctx, modelRequest)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = responseModel
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)

}
