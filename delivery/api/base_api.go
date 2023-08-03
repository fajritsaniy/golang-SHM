package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/delivery/api/response"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/mitchellh/mapstructure"
)

type BaseApi struct{}

func (b *BaseApi) ParseRequestBody(c *gin.Context, payload interface{}) error {
	if err := c.ShouldBindJSON(payload); err != nil {
		return err
	}
	return nil
}

func (b *BaseApi) ParseRequestFormData(c *gin.Context, requestModel interface{}, postFormKey ...string) error {
	mapRes := make(map[string]interface{})
	for _, v := range postFormKey {
		mapRes[v] = c.PostForm(v)
	}
	err := mapstructure.Decode(mapRes, &requestModel)
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseApi) NewSuccessSingleResponse(c *gin.Context, data interface{}, responseType string) {
	response.SendSingleResponse(c, data, responseType)
}

func (b *BaseApi) NewSuccessPageResponse(c *gin.Context, data []interface{}, responseType string, paging dto.Paging) {
	response.SendPageResponse(c, data, responseType, paging)
}

func (b *BaseApi) NewErrorErrorResponse(c *gin.Context, code int, errorMessage string) {
	response.SendErrorResponse(c, code, errorMessage)
}

func (b *BaseApi) NewSuccessFileResponse(c *gin.Context, fileName string, responseType string) {
	response.SendFileResponse(c, fileName, responseType)
}
