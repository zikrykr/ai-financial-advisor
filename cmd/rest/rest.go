package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ai-financial-advisor/cmd/setup"
	appSetup "github.com/ai-financial-advisor/cmd/setup"
	"github.com/ai-financial-advisor/config"
	"github.com/ai-financial-advisor/constants"
	"github.com/ai-financial-advisor/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BaseURL base url of api
const BaseURL = "/api/v1"

func StartServer(setupData appSetup.SetupData) {
	conf := config.GetConfig()
	// appName := conf.App.Name
	if conf.App.Env == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	// GIN Init
	router := gin.Default()
	router.UseRawPath = true

	tracerProvider, err := setup.InitTracer()
	if err != nil {
		logrus.Fatal("failed to init tracer", err)
	}
	defer tracerProvider.Shutdown(context.Background())

	router.Use(middleware.CORSMiddleware())

	// init public route
	initPublicRoute(router, setupData.InternalApp)

	router.Use(middleware.JwtAuthMiddleware())

	//Init Main APP and Route
	initRoute(router, setupData.InternalApp)

	port := config.GetConfig().Http.Port
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil {
			logrus.Error(fmt.Printf("listen: %s\n", err))
		}
	}()
	logrus.Info("webserver started")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Panic("Server Shutdown:", err)
	}

	_ = appSetup.CloseDB()

	logrus.Info("Server exiting")
}

func initRoute(router *gin.Engine, internalAppStruct appSetup.InternalAppStruct) {
	// r := router.Group(BaseURL)
}

func initPublicRoute(router *gin.Engine, internalAppStruct appSetup.InternalAppStruct) {
	r := router.Group(BaseURL)
	r.GET("/healthz", internalAppStruct.Handler.HealthzHandler.Healthz)
}
