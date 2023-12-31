package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fajritsaniy/golang-SHM/delivery/middleware"

	"github.com/fajritsaniy/golang-SHM/delivery/api"
	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/usecase"
	"github.com/fajritsaniy/golang-SHM/utils/common"
	"github.com/gin-gonic/gin"
)

type VehicleController struct {
	router  *gin.Engine
	usecase usecase.VehicleUseCase
	api.BaseApi
}

func (v *VehicleController) createHandler(c *gin.Context) {
	vehicle := c.PostForm("vehicle")
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusBadRequest, "Failed Get File")
	}
	log.Println(fileHeader.Filename)
	fileName := strings.Split(fileHeader.Filename, ".")
	if len(fileName) != 2 {
		v.NewErrorErrorResponse(c, http.StatusBadRequest, "Unrecognized file extension")
	}
	var payload model.Vehicle
	err = json.Unmarshal([]byte(vehicle), &payload)
	if err != nil {
		log.Println("failed to unmarshal")
	}
	if err := v.usecase.UploadImage(&payload, file, fileName[1]); err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	v.NewSuccessSingleResponse(c, payload, "OK")
}

func (v *VehicleController) updateHandler(c *gin.Context) {
	var payload model.Vehicle
	if err := c.ShouldBindJSON(&payload); err != nil {
		v.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := v.usecase.SaveData(&payload); err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	v.NewSuccessSingleResponse(c, payload, "OK")
}

func (v *VehicleController) listHandler(c *gin.Context) {
	requestQueryParams, err := common.ValidateRequestQueryParams(c)
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	vehicles, paging, err := v.usecase.Paging(requestQueryParams)
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var vehicleInterface []interface{}
	for _, v := range vehicles {
		vehicleInterface = append(vehicleInterface, v)
	}
	v.NewSuccessPageResponse(c, vehicleInterface, "OK", paging)
}

func (v *VehicleController) getByIDHandler(c *gin.Context) {
	id := c.Param("id")
	vehicle, err := v.usecase.FindById(id)
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	v.NewSuccessSingleResponse(c, vehicle, "OK")
}

func (v *VehicleController) getImageByIDHandler(c *gin.Context) {
	id := c.Param("id")
	vehicle, err := v.usecase.FindById(id)
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	v.NewSuccessFileResponse(c, vehicle.ImgPath, "OK")
}

func (v *VehicleController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := v.usecase.DeleteData(id)
	if err != nil {
		v.NewErrorErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusNoContent, "")
}

func NewVehicleController(r *gin.Engine, usecase usecase.VehicleUseCase, authMiddleware middleware.AuthTokenMiddleware) *VehicleController {
	controller := VehicleController{
		router:  r,
		usecase: usecase,
	}

	const vehicleEndpoint = "/vehicles"
	r.GET(vehicleEndpoint, controller.listHandler)
	r.POST(vehicleEndpoint, authMiddleware.RequireToken(), controller.createHandler)
	r.PUT(vehicleEndpoint, authMiddleware.RequireToken(), controller.updateHandler)
	r.GET("/vehicles/:id", controller.getByIDHandler)
	r.GET("/vehicles/image/:id", authMiddleware.RequireToken(), controller.getImageByIDHandler)
	r.DELETE("/vehicles/:id", authMiddleware.RequireToken(), controller.deleteHandler)
	return &controller
}
