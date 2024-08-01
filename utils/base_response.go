package utils

import "github.com/gin-gonic/gin"

type BaseResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	IsSuccess  bool   `json:"is_success"`
	Message    string `json:"message"`
	Data       T      `json:"data,omitempty"`
}

func SendResponse[T any](c *gin.Context, statusCode int, isSuccess bool, message string, data T) {
	response := BaseResponse[T]{
		StatusCode: statusCode,
		IsSuccess:  isSuccess,
		Message:    message,
		Data:       data,
	}
	c.JSON(statusCode, response)
}

func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	SendResponse[interface{}](c, statusCode, false, message, nil)
}
