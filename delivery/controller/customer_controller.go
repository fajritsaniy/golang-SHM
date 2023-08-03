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

type CustomerController struct {
	router         *gin.Engine
	usecase        usecase.CustomerUseCase
	authMiddleware middleware.AuthTokenMiddleware
	api.BaseApi
}

func (cc *CustomerController) createHandler(c *gin.Context) {
	var payload model.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		cc.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cc.usecase.SaveData(&payload); err != nil {
		cc.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cc.NewSuccessSingleResponse(c, payload, "OK")
}

func (cc *CustomerController) updateHandler(c *gin.Context) {
	var payload model.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		cc.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := cc.usecase.SaveData(&payload); err != nil {
		cc.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cc.NewSuccessSingleResponse(c, payload, "OK")
}

func (cc *CustomerController) listHandler(c *gin.Context) {
	customers, err := cc.usecase.FindAll()

	if err != nil {
		cc.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var customerInterface []interface{}
	for _, v := range customers {
		customerInterface = append(customerInterface, v)
	}
	cc.NewSuccessPageResponse(c, customerInterface, "OK", dto.Paging{})
}

func (cc *CustomerController) getByIDHandler(c *gin.Context) {
	id := c.Param("id")
	vehicle, err := cc.usecase.FindById(id)
	if err != nil {
		cc.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cc.NewSuccessSingleResponse(c, vehicle, "OK")
}

func (cc *CustomerController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := cc.usecase.DeleteData(id)
	if err != nil {
		cc.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewCustomerController(r *gin.Engine, usecase usecase.CustomerUseCase, authMiddleware middleware.AuthTokenMiddleware) *CustomerController {
	controller := CustomerController{
		router:  r,
		usecase: usecase,
	}
	r.GET("/customers", authMiddleware.RequireToken(), controller.listHandler)
	r.GET("/customers/:id", authMiddleware.RequireToken(), controller.getByIDHandler)
	r.POST("/customers", authMiddleware.RequireToken(), controller.createHandler)
	r.PUT("/customers", authMiddleware.RequireToken(), controller.updateHandler)
	r.DELETE("/customers/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
