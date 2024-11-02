package controllers

import (
	"net/http"
	M "perpus_api/models"

	"github.com/labstack/echo/v4"
)
func CreateReviews(c echo.Context) error {
	requestBody := new(M.Reviews)

	if err := c.Bind(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Validate(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"msg": *errm})
	}

	res, err := M.StoreReviews(*requestBody, *claims)

	if err != nil {
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, res)
}