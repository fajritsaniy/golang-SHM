package manager

import (
	"fmt"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/config"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InfraManager interface {
	Conn() *gorm.DB
	Migrate(model ...any) error
	Log() *logrus.Logger
	LogFilePath() string
	UploadLocation() string
}

type infraManager struct {
	db  *gorm.DB
	cfg *config.Config
	log *logrus.Logger
}

func (i *infraManager) UploadLocation() string {
	return i.cfg.UploadLocation
}

func (i *infraManager) LogFilePath() string {
	return i.cfg.LogFilePath
}

func (i *infraManager) Log() *logrus.Logger {
	logger := logrus.New()
	return logger
}

func (i *infraManager) Conn() *gorm.DB {
	return i.db
}

func (i *infraManager) Migrate(model ...any) error {
	err := i.Conn().AutoMigrate(model...)
	if err != nil {
		return err
	}
	return nil
}

func (i *infraManager) initDb() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		i.cfg.Host,
		i.cfg.Port,
		i.cfg.User,
		i.cfg.Password,
		i.cfg.Name,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	i.db = conn
	if i.cfg.FileConfig.Env == "MIGRATION" {
		i.db = conn.Debug()
		err := i.Migrate(
			&model.Brand{},
			&model.Vehicle{},
			&model.UserCredential{},
			&model.Customer{},
			&model.Employee{},
			&model.Transaction{},
		)
		if err != nil {
			return err
		}
	} else if i.cfg.FileConfig.Env == "DEV" {
		i.db = conn.Debug()
	} else {
		// production / release
	}
	return nil
}

func NewInfraManager(cfg *config.Config) (InfraManager, error) {
	conn := &infraManager{cfg: cfg}
	err := conn.initDb()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
