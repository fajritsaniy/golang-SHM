package delivery

import (
	"fmt"

	"github.com/fajritsaniy/golang-SHM/config"
	"github.com/fajritsaniy/golang-SHM/delivery/controller"
	"github.com/fajritsaniy/golang-SHM/delivery/middleware"
	"github.com/fajritsaniy/golang-SHM/manager"
	"github.com/fajritsaniy/golang-SHM/usecase"
	"github.com/fajritsaniy/golang-SHM/utils/security"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ucManager    manager.UseCaseManager
	authUseCase  usecase.AuthenticationUseCase
	tokenService security.AccessToken
	engine       *gin.Engine
	host         string
	log          *logrus.Logger
}

func (s *Server) initController() {
	s.engine.Use(middleware.LogRequestMiddleware(s.log))
	authMiddleware := middleware.NewTokenValidator(s.tokenService)
	controller.NewVehicleController(s.engine, s.ucManager.VehicleUseCase(), authMiddleware)
	controller.NewBrandController(s.engine, s.ucManager.BrandUseCase(), authMiddleware)
	controller.NewCustomerController(s.engine, s.ucManager.CustomerUseCase(), authMiddleware)
	controller.NewEmployeeController(s.engine, s.ucManager.EmployeeUseCase(), authMiddleware)
	controller.NewTransactionController(s.engine, s.ucManager.TransactionUseCase(), authMiddleware)
	controller.NewAuthController(s.engine, s.authUseCase)
}

func NewServer() *Server {
	c, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
	}
	// infra manager
	infraManager, _ := manager.NewInfraManager(c)
	// repo manager
	repoManager := manager.NewRepositoryManager(infraManager)
	// use case manager
	useCaseManager := manager.NewUseCaseManager(repoManager)
	// token
	tokenService := security.NewAccessToken(c.TokenConfig)
	authUseCase := usecase.NewAuthenticationUseCase(repoManager.UserRepo(), tokenService)

	r := gin.Default()
	host := fmt.Sprintf("%s:%s", c.ApiHost, c.ApiPort)
	return &Server{
		ucManager:    useCaseManager,
		authUseCase:  authUseCase,
		tokenService: tokenService,
		engine:       r,
		host:         host,
		log:          infraManager.Log(),
	}
}

func (s *Server) Run() {
	s.initController()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}
