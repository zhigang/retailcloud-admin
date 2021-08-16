package main

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/zhigang/retailcloud-admin/config"
	"github.com/zhigang/retailcloud-admin/controllers"
	"github.com/zhigang/retailcloud-admin/factory"
)

var (
	globalConfig *config.Config
	startedAt    time.Time
	client       *retailcloud.Client
)

func main() {

	startedAt = time.Now()

	globalConfig = factory.GlobalConfig()
	client = factory.GetRetailCloudClient()

	startEchoServer()
}

func startEchoServer() {
	e := initEchoServer()

	// Start server
	go func() {
		s := &http.Server{
			Addr:         globalConfig.Service.Address,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		}
		if err := e.StartServer(s); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 30 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initEchoServer() *echo.Echo {

	e := echo.New()

	if strings.ToLower(globalConfig.Log.Level) == "debug" {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(globalConfig.Service.BasicAuth.Username)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(globalConfig.Service.BasicAuth.Password)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	// 	KeyLookup: "header:api-key",
	// 	Validator: func(key string, c echo.Context) (bool, error) {
	// 		return key == globalConfig.Service.APIKey, nil
	// 	},
	// }))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"startedAt": startedAt,
		})
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	bindingAPI(e)

	return e
}

func bindingAPI(e *echo.Echo) {
	apiV1 := e.Group("/v1")

	var app controllers.AppController

	v1App := apiV1.Group("/app")
	v1App.GET("", app.GetAppList)
	v1App.GET("/:id/env", app.GetEnvList)
	v1App.PUT("/:id/deploy", app.DeployApp)
}
