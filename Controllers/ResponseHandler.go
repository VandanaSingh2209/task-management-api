package Controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	StatusCode   int
	Message      string
	ErrorMessage error
	Body         map[string]interface{}
}

func ValidationResponse(c *gin.Context, Message string) {
	response := ResponseBody{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    Message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusUnprocessableEntity, response)
}

func NoDataFoundResponse(c *gin.Context, Message string) {
	response := ResponseBody{
		StatusCode: http.StatusNoContent,
		Message:    Message,
		Body:       map[string]interface{}{},
	}
	c.JSON(http.StatusOK, response)
}

func successResponse(c *gin.Context, Message string, Body map[string]interface{}) {
	response := ResponseBody{
		StatusCode: http.StatusOK,
		Message:    Message,
		Body:       Body,
	}
	c.JSON(http.StatusOK, response)
}
