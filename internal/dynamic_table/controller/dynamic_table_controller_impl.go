package controller

import (
	"net/http"

	response "github.com/Ganasa18/simple-crud-builder-go/internal/base/models/web"
	"github.com/Ganasa18/simple-crud-builder-go/internal/dynamic_table/service"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/helper"
	"github.com/gin-gonic/gin"
)

type DynamicTableControllerImpl struct {
	DynamicTableService service.DynamicTableService
}

func NewDynamicTableController(dynamicTableService service.DynamicTableService) DynamicTableController {
	return &DynamicTableControllerImpl{
		DynamicTableService: dynamicTableService,
	}
}

func (controller *DynamicTableControllerImpl) CreateRecord(ctx *gin.Context) {

	var result interface{}
	helper.ReadFromRequestBody(ctx.Request, &result)

	responseValue, err := controller.DynamicTableService.CreateRecord(ctx, result)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = responseValue
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)

}

func (controller *DynamicTableControllerImpl) ListRecord(ctx *gin.Context) {

	responseTable, err := controller.DynamicTableService.ListRecord(ctx)
	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = responseTable
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)
}
