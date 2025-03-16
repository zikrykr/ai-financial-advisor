package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ai-financial-advisor/config"
	"github.com/ai-financial-advisor/constants"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	*gorm.DB
}

type DBConfig struct {
	GormDB       *GormDB
	ConnectionDB *sql.DB
}

func (db DBConfig) CloseConnection() error {
	return db.ConnectionDB.Close()
}

func Init() (DBConfig, error) {
	var (
		dbConfigVar DBConfig
		loggerGorm  logger.Interface
	)
	configData := config.GetConfig()

	loggerGorm = logger.Default.LogMode(logger.Silent)
	if configData.App.Env == constants.DEV {
		loggerGorm = logger.Default.LogMode(logger.Info)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", configData.DB.Host, configData.DB.User, configData.DB.Pass, configData.DB.Name),
	}), &gorm.Config{
		Logger:                 loggerGorm,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return dbConfigVar, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return dbConfigVar, err
	}

	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(configData.DB.MaxIdletimeConn))
	sqlDB.SetMaxIdleConns(configData.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(configData.DB.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(configData.DB.MaxLifetimeConn))
	dbConfigVar.ConnectionDB = sqlDB

	dbConfigVar.GormDB = &GormDB{gormDB}
	logrus.Info("database connected")

	return dbConfigVar, nil
}

func (db DBConfig) HealthDBCheck(ctx context.Context) error {
	tracer := otel.Tracer("health-db")
	_, span := tracer.Start(ctx, "HealthDBCheck")
	defer span.End()
	if err := db.ConnectionDB.Ping(); err != nil {
		return err
	}

	return nil
}
