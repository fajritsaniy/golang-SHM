package controller

import (
	"net/http"

	"github.com/fajritsaniy/golang-SHM/delivery/api"
	"github.com/fajritsaniy/golang-SHM/delivery/middleware"
	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/model/dto"
	"github.com/fajritsaniy/golang-SHM/usecase"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	router  *gin.Engine
	usecase usecase.TransactionUseCase
	api.BaseApi
}

func (e *TransactionController) createHandler(c *gin.Context) {
	var payload model.Transaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := e.usecase.RegisterNewTransaction(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, payload, "OK")
}

func (e *TransactionController) listHandler(c *gin.Context) {
	transactions, err := e.usecase.FindAllTransaction()

	if err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var transactionInterface []interface{}
	for _, v := range transactions {
		transactionInterface = append(transactionInterface, v)
	}
	e.NewSuccessPageResponse(c, transactionInterface, "OK", dto.Paging{})
}

func (e *TransactionController) getByIDHandler(c *gin.Context) {
	id := c.Param("id")
	transaction, err := e.usecase.FindByTransaction(id)
	if err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, transaction, "OK")
}

func NewTransactionController(r *gin.Engine, usecase usecase.TransactionUseCase, authMiddleware middleware.AuthTokenMiddleware) *TransactionController {
	controller := TransactionController{
		router:  r,
		usecase: usecase,
	}
	r.GET("/transactions", authMiddleware.RequireToken(), controller.listHandler)
	r.GET("/transactions/:id", authMiddleware.RequireToken(), controller.getByIDHandler)
	r.POST("/transactions", authMiddleware.RequireToken(), controller.createHandler)
	return &controller
}
