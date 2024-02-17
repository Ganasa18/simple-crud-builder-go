package exception

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Ganasa18/simple-crud-builder-go/internal/base/models/web"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/helper"
	"github.com/Ganasa18/simple-crud-builder-go/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Middleware to recover from panics and handle errors
func ExceptionRecoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Log the panic details
			fmt.Println("Panic recovered:", err)
			debug.PrintStack()

			// Handle the error (you can customize this function)
			handleError(c, err)

			// Abort the request to prevent further execution
			c.Abort()
		}
	}()

	// Continue with the request handling
	c.Next()
}

func handleError(ctx *gin.Context, err interface{}) {
	if notFoundError(ctx, err) {
		return
	}

	if validationErrors(ctx, err) {
		return
	}

	internalServerError(ctx, err)
}

type ValidationError struct {
	FieldName string `json:"field_name"`
}

func validationErrors(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)

	if ok {
		ctx.Writer.Header().Set(utils.HEADER_CONTENT_TYPE, utils.CONTENT_TYPE_APPLICATION_JSON)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		// Convert the validation errors into an array
		errorArray := make([]ValidationError, len(exception))
		for i, fieldError := range exception {
			errorArray[i] = ValidationError{
				FieldName: "Field " + helper.CamelToSnake(fieldError.Field()) + " " + fieldError.Tag() + " error",
			}
		}
		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   map[string]interface{}{"errors": errorArray},
		}
		helper.WriteToResponseBody(ctx, http.StatusBadRequest, webResponse)
		return true
	} else {
		return false
	}
}
func notFoundError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		ctx.Writer.Header().Set(utils.HEADER_CONTENT_TYPE, utils.CONTENT_TYPE_APPLICATION_JSON)
		ctx.Writer.WriteHeader(http.StatusNotFound)
		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error,
		}

		helper.WriteToResponseBody(ctx, http.StatusNotFound, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(ctx *gin.Context, err interface{}) {
	ctx.Writer.Header().Set(utils.HEADER_CONTENT_TYPE, utils.CONTENT_TYPE_APPLICATION_JSON)
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
	webResponse := web.WebResponse{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data:   err,
	}

	helper.WriteToResponseBody(ctx, http.StatusInternalServerError, webResponse)
}
