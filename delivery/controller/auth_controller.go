package controller

import (
	"fmt"
	"net/http"

	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/usecase"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	router  *gin.Engine
	usecase usecase.AuthenticationUseCase
}

func (a *AuthController) loginHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	fmt.Println(payload.UserName, payload.Password)
	token, err := a.usecase.Login(payload.UserName, payload.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code":  http.StatusCreated,
		"token": token,
	})
}

func (a *AuthController) registerHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	err := a.usecase.Register(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	successMessage := fmt.Sprintf("%s has been registered.", payload.UserName)

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": successMessage,
	})
}

func (a *AuthController) userActivationHandler(c *gin.Context) {
	var payload model.UserCredential
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	status, err := a.usecase.UserActivation(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	var successMessage string
	if status {
		successMessage = fmt.Sprintf("%s has been activated.", payload.UserName)
	} else {
		successMessage = fmt.Sprintf("%s has been disabled.", payload.UserName)
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": successMessage,
	})
}

func NewAuthController(r *gin.Engine, usecase usecase.AuthenticationUseCase) *AuthController {
	controller := AuthController{
		router:  r,
		usecase: usecase,
	}
	r.POST("/login", controller.loginHandler)
	r.POST("/register", controller.registerHandler)
	r.POST("/activation", controller.userActivationHandler)
	return &controller
}
