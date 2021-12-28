package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	StaticsPath string `envconfig:"STATICS_PATH" default:"./static"`
}

func main() {
	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalf("Can't process config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := Serve(ctx, config)
		if err != nil {
			log.Println(fmt.Errorf("serve: %w", err))
			return
		}
	}()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("OS interrupting signal has received")

	cancel()

}

func Serve(ctx context.Context, config *Config) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Recover())

	InitHandlers(e, config.StaticsPath)

	go func() {
		e.Logger.Infof("start server on port: %s", config.Port)
		err := e.Start(":" + config.Port)
		if err != nil {
			e.Logger.Errorf("start server error: %v", err)
		}
	}()

	<-ctx.Done()

	return e.Shutdown(ctx)
}

func InitHandlers(e *echo.Echo, staticsPath string) {
	e.GET("/", handler)
	e.GET("/__heartbeat__", heartbeatHandler)
	e.Static("/static", staticsPath)

	e.Any("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}

func handler(c echo.Context) error {
	now := time.Now().Format(time.RFC822)
	return c.String(http.StatusOK, "Current time: "+now)
}

func heartbeatHandler(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
