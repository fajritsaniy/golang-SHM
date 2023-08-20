package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/delivery/api"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/delivery/middleware"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/usecase"
)

type CustomerController struct {
	router  *gin.Engine
	usecase usecase.CustomerUseCase
	api.BaseApi
}

func (cc *CustomerController) createUpdateHandler(c *gin.Context) {
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

	const customerEndpoint = "/customers"
	r.GET(customerEndpoint, authMiddleware.RequireToken(), controller.listHandler)
	r.GET("/customers/:id", authMiddleware.RequireToken(), controller.getByIDHandler)
	r.POST(customerEndpoint, authMiddleware.RequireToken(), controller.createUpdateHandler)
	r.PUT(customerEndpoint, authMiddleware.RequireToken(), controller.createUpdateHandler)
	r.DELETE("/customers/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
