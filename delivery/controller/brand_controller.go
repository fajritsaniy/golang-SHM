package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/delivery/api"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/delivery/middleware"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/usecase"
	"net/http"
)

type BrandController struct {
	router         *gin.Engine
	usecase        usecase.BrandUseCase
	authMiddleware middleware.AuthTokenMiddleware
	api.BaseApi
}

func (b *BrandController) createHandler(c *gin.Context) {
	var payload model.Brand
	if err := c.ShouldBindJSON(&payload); err != nil {
		b.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := b.usecase.SaveData(&payload); err != nil {
		b.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	b.NewSuccessSingleResponse(c, payload, "OK")
}

func (b *BrandController) updateHandler(c *gin.Context) {
	var payload model.Brand
	if err := c.ShouldBindJSON(&payload); err != nil {
		b.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := b.usecase.SaveData(&payload); err != nil {
		b.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	b.NewSuccessSingleResponse(c, payload, "OK")
}

func (b *BrandController) listHandler(c *gin.Context) {
	vehicles, err := b.usecase.FindAll()

	if err != nil {
		b.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var brandInterface []interface{}
	for _, v := range vehicles {
		brandInterface = append(brandInterface, v)
	}
	b.NewSuccessPageResponse(c, brandInterface, "OK", dto.Paging{})
}

func (b *BrandController) getByIDHandler(c *gin.Context) {
	id := c.Param("id")
	vehicle, err := b.usecase.FindById(id)
	if err != nil {
		b.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	b.NewSuccessSingleResponse(c, vehicle, "OK")
}

func (b *BrandController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := b.usecase.DeleteData(id)
	if err != nil {
		b.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewBrandController(r *gin.Engine, usecase usecase.BrandUseCase, authMiddleware middleware.AuthTokenMiddleware) *BrandController {
	controller := BrandController{
		router:         r,
		usecase:        usecase,
		authMiddleware: authMiddleware,
	}
	r.GET("/brands", controller.listHandler)
	r.GET("/brands/:id", controller.getByIDHandler)
	r.POST("/brands", authMiddleware.RequireToken(), controller.createHandler)
	r.PUT("/brands", authMiddleware.RequireToken(), controller.updateHandler)
	r.DELETE("/brands/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
