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

type EmployeeController struct {
	router         *gin.Engine
	usecase        usecase.EmployeeUseCase
	authMiddleware middleware.AuthTokenMiddleware
	api.BaseApi
}

func (e *EmployeeController) createHandler(c *gin.Context) {
	var payload model.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := e.usecase.SaveData(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, payload, "OK")
}

func (e *EmployeeController) updateHandler(c *gin.Context) {
	var payload model.Employee
	if err := c.ShouldBindJSON(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := e.usecase.SaveData(&payload); err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, payload, "OK")
}

func (e *EmployeeController) listHandler(c *gin.Context) {
	employees, err := e.usecase.FindAll()

	if err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var employeeInterface []interface{}
	for _, v := range employees {
		employeeInterface = append(employeeInterface, v)
	}
	e.NewSuccessPageResponse(c, employeeInterface, "OK", dto.Paging{})
}

func (e *EmployeeController) getByIDHandler(c *gin.Context) {
	id := c.Param("id")
	employee, err := e.usecase.FindById(id)
	if err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	e.NewSuccessSingleResponse(c, employee, "OK")
}

func (e *EmployeeController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := e.usecase.DeleteData(id)
	if err != nil {
		e.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewEmployeeController(r *gin.Engine, usecase usecase.EmployeeUseCase, authMiddleware middleware.AuthTokenMiddleware) *EmployeeController {
	controller := EmployeeController{
		router:  r,
		usecase: usecase,
	}
	r.GET("/employees", authMiddleware.RequireToken(), controller.listHandler)
	r.GET("/employees/:id", authMiddleware.RequireToken(), controller.getByIDHandler)
	r.POST("/employees", authMiddleware.RequireToken(), controller.createHandler)
	r.PUT("/employees", authMiddleware.RequireToken(), controller.updateHandler)
	r.DELETE("/employees/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
