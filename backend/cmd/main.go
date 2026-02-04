package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type SimpleJsonData = struct {
	Hello      string `json:"hello"`
	Foo        string `json:"foo"`
	ApiVersion int    `json:"apiVersion"`
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())

	e.GET("/hello-world", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, SimpleJsonData{
			Hello:      "world",
			Foo:        "bar",
			ApiVersion: 2,
		})
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
