package httputil

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewError(ctx *gin.Context, status string, err error) {
	er := ErrorResponseModel{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(http.StatusOK, er)
}

// HTTPError example
type ErrorResponseModel struct {
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func (err *ErrorResponseModel) Error(ctx *gin.Context, code string, errMessage string) {
	err.Code = code
	err.Message = fmt.Sprint(errMessage)
}

// HTTPResponse
type ResponseModel struct {
	Data   any  `json:"data,omitempty"`
	Result bool `json:"result"`
	Error  any  `json:"error,omitempty"`
}

// PageModel
type PageModel struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}
