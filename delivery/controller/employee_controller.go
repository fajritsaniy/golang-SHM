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

type EmployeeController struct {
	router  *gin.Engine
	usecase usecase.EmployeeUseCase
	api.BaseApi
}

func (e *EmployeeController) createUpdateHandler(c *gin.Context) {
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

	const employeeEndpoint = "/employee"
	r.GET(employeeEndpoint, authMiddleware.RequireToken(), controller.listHandler)
	r.GET("/employees/:id", authMiddleware.RequireToken(), controller.getByIDHandler)
	r.POST(employeeEndpoint, authMiddleware.RequireToken(), controller.createUpdateHandler)
	r.PUT(employeeEndpoint, authMiddleware.RequireToken(), controller.createUpdateHandler)
	r.DELETE("/employees/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
