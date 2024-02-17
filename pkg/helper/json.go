package helper

import (
	"encoding/json"
	"net/http"

	"github.com/Ganasa18/simple-crud-builder-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	utils.IsErrorDoPanic(err)
}

func WriteToResponseBody(ctx *gin.Context, statusCode int, response interface{}) {
	ctx.JSON(statusCode, response)
}
