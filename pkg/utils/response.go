package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, APIResponse{
		Data: data,
	})
}

func SuccessResponseWithMeta(c *gin.Context, code int, data interface{}, meta interface{}) {
	c.JSON(code, APIResponse{
		Data: data,
		Meta: meta,
	})
}

func SuccessMessage(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Message: message,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Error: message,
	})
}