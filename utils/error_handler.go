package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardError struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

type ValidationErrorResponse struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
	Field            string `json:"field"`
	Message          string `json:"message"`
}

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (eh *ErrorHandler) NotFound(c *gin.Context, message ...string) {
	defaultMessage := "RESOURCE_NOT_FOUND"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	
	c.JSON(http.StatusNotFound, StandardError{
		ErrorCode:        404,
		ErrorDescription: defaultMessage,
	})
}

func (eh *ErrorHandler) BadRequest(c *gin.Context, message ...string) {
	defaultMessage := "BAD_REQUEST"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	
	c.JSON(http.StatusBadRequest, StandardError{
		ErrorCode:        400,
		ErrorDescription: defaultMessage,
	})
}

func (eh *ErrorHandler) Unauthorized(c *gin.Context, message ...string) {
	defaultMessage := "UNAUTHORIZED"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	
	c.JSON(http.StatusUnauthorized, StandardError{
		ErrorCode:        401,
		ErrorDescription: defaultMessage,
	})
}

func (eh *ErrorHandler) Forbidden(c *gin.Context, message ...string) {
	defaultMessage := "FORBIDDEN"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	
	c.JSON(http.StatusForbidden, StandardError{
		ErrorCode:        403,
		ErrorDescription: defaultMessage,
	})
}

func (eh *ErrorHandler) InternalServerError(c *gin.Context, message ...string) {
	defaultMessage := "INTERNAL_SERVER_ERROR"
	if len(message) > 0 {
		defaultMessage = message[0]
	}
	
	c.JSON(http.StatusInternalServerError, StandardError{
		ErrorCode:        500,
		ErrorDescription: defaultMessage,
	})
}

func (eh *ErrorHandler) ValidationError(c *gin.Context, field string, message string) {
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		ErrorCode:        400,
		ErrorDescription: "VALIDATION_ERROR",
		Field:            field,
		Message:          message,
	})
}

func (eh *ErrorHandler) CustomError(c *gin.Context, statusCode int, errorCode int, message string) {
	c.JSON(statusCode, StandardError{
		ErrorCode:        errorCode,
		ErrorDescription: message,
	})
}

var globalErrorHandler = NewErrorHandler()

func ReportNotFound(c *gin.Context, message ...string) {
	globalErrorHandler.NotFound(c, message...)
}

func ReportBadRequest(c *gin.Context, message ...string) {
	globalErrorHandler.BadRequest(c, message...)
}

func ReportUnauthorized(c *gin.Context, message ...string) {
	globalErrorHandler.Unauthorized(c, message...)
}

func ReportForbidden(c *gin.Context, message ...string) {
	globalErrorHandler.Forbidden(c, message...)
}

func ReportInternalServerError(c *gin.Context, message ...string) {
	globalErrorHandler.InternalServerError(c, message...)
}

func ReportValidationError(c *gin.Context, field string, message string) {
	globalErrorHandler.ValidationError(c, field, message)
}

func ReportCustomError(c *gin.Context, statusCode int, errorCode int, message string) {
	globalErrorHandler.CustomError(c, statusCode, errorCode, message)
}

func HandleNoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		ReportNotFound(c, "Route not found")
	}
}