package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func saveItem(c echo.Context) error {

	i := new(Item)
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := i.Save(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, i)
}

func getItems(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("name")
	price := c.FormValue("price")
	quantity := c.FormValue("quantity")
	description := c.FormValue("description")

	i := new(Item)
	items, err := i.Search(id, name, description, price, quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, items)
}
