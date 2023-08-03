package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
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
