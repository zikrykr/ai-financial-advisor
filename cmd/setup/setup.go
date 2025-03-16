package setup

import (
	"github.com/ai-financial-advisor/config"
	"github.com/ai-financial-advisor/config/db"
	healthHandler "github.com/ai-financial-advisor/internal/healthz/handler"
	healthPort "github.com/ai-financial-advisor/internal/healthz/ports"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SetupData struct {
	ConfigData  config.Config
	InternalApp InternalAppStruct
	DB          db.DBConfig
}

type InternalAppStruct struct {
	Repositories initRepositoriesApp
	Services     initServicesApp
	Handler      InitHandlerApp
}

// Repositories
type initRepositoriesApp struct {
	dbInstance *gorm.DB
}

// Services
type initServicesApp struct {
}

// Handler
type InitHandlerApp struct {
	dbInstance     *gorm.DB
	HealthzHandler healthPort.IHealthHandler
}

// CloseDB close connection to db
var CloseDB func() error

func InitSetup() SetupData {
	configData := config.GetConfig()

	//DB INIT
	dbConn, err := db.Init()
	if err != nil {
		logrus.Fatal("database error", err)
	}

	CloseDB = func() error {
		if err := dbConn.CloseConnection(); err != nil {
			return err
		}

		return nil
	}

	internalAppVar := initInternalApp(dbConn.GormDB)

	return SetupData{
		ConfigData:  configData,
		InternalApp: internalAppVar,
		DB:          dbConn,
	}
}

func initInternalApp(gormDB *db.GormDB) InternalAppStruct {
	var internalAppVar InternalAppStruct

	initAppRepo(gormDB, &internalAppVar)
	initAppService(&internalAppVar)
	initAppHandler(gormDB, &internalAppVar)

	return internalAppVar
}

func initAppRepo(gormDB *db.GormDB, initializeApp *InternalAppStruct) {
	// Get Gorm instance
	initializeApp.Repositories.dbInstance = gormDB.DB
}

func initAppService(initializeApp *InternalAppStruct) {
}

func initAppHandler(gormDB *db.GormDB, initializeApp *InternalAppStruct) {
	// Get Gorm instance
	initializeApp.Handler.dbInstance = gormDB.DB

	// Healthz Handler
	initializeApp.Handler.HealthzHandler = healthHandler.NewHealthHandler(gormDB.DB)
}
