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

func (controller *DynamicTableControllerImpl) GetRecord(ctx *gin.Context) {
	getResponse, err := controller.DynamicTableService.GetRecord(ctx)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = getResponse
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)
}

func (controller *DynamicTableControllerImpl) DeleteRecord(ctx *gin.Context) {
	err := controller.DynamicTableService.DeleteRecord(ctx)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = "sucess delete record"
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)
}

func (controller *DynamicTableControllerImpl) SoftDeleteRecord(ctx *gin.Context) {
	err := controller.DynamicTableService.SoftDeleteRecord(ctx)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = "sucess soft delete record"
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)
}

func (controller *DynamicTableControllerImpl) UpdateRecord(ctx *gin.Context) {


	var request map[string]interface{}
	helper.ReadFromRequestBody(ctx.Request, &request)

	updateResponse, err := controller.DynamicTableService.UpdateRecord(ctx, request)

	var statusCode int
	var responseData interface{}

	if err != nil {
		statusCode = http.StatusBadRequest
		responseData = err.Error()
	} else {
		statusCode = http.StatusOK
		responseData = updateResponse
	}

	webResponse := response.WebResponse{
		Code:   statusCode,
		Status: http.StatusText(statusCode),
		Data:   responseData,
	}

	helper.WriteToResponseBody(ctx, statusCode, webResponse)
}
