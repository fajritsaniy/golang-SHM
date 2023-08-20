package response

import (
	"net/http"

	"github.com/fajritsaniy/golang-SHM/model/dto"
	"github.com/gin-gonic/gin"
)

func SendSingleResponse(c *gin.Context, data interface{}, responseType string) {
	c.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: responseType,
		},
		Data: data,
	})
}

func SendPageResponse(c *gin.Context, data []interface{}, responseType string, paging dto.Paging) {
	c.JSON(http.StatusOK, &PagedResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: responseType,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendErrorResponse(c *gin.Context, code int, errorMessage string) {
	c.AbortWithStatusJSON(code, &Status{
		Code:        code,
		Description: errorMessage,
	})
}

func SendFileResponse(c *gin.Context, fileName string, responseType string) {
	c.JSON(http.StatusOK, &FileResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: responseType,
		},
		FileName: fileName,
	})
}
