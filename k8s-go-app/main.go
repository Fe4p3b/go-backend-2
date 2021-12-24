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

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := "8080"

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := Serve(ctx, port)
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

func Serve(ctx context.Context, port string) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Recover())

	InitHandlers(e)

	go func() {
		e.Logger.Infof("start server on port: %s", port)
		err := e.Start(":" + port)
		if err != nil {
			e.Logger.Errorf("start server error: %v", err)
		}
	}()

	<-ctx.Done()

	return e.Shutdown(ctx)
}

func InitHandlers(e *echo.Echo) {
	e.GET("/", handler)
	e.GET("/__heartbeat__", heartbeatHandler)

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
