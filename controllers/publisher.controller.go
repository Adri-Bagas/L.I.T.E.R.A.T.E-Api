package controllers

import (
	"net/http"
	M "perpus_api/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PublisherForm struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Desc        string `json:"desc" validate:"required"`
}

func GetAllPublisher(c echo.Context) error {
	res, err := M.GetAllPublisher()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func CreatePublisher(c echo.Context) error {

	p := &PublisherForm{
		Name:        c.FormValue("name"),
		Address:     c.FormValue("address"),
		PhoneNumber: c.FormValue("phone_number"),
		Desc:        c.FormValue("desc"),
	}

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedCB, _ := strconv.ParseInt(c.FormValue("created_by"), 10, 64)

	pu := &M.Publisher{
		Name:        c.FormValue("name"),
		Address:     c.FormValue("address"),
		PhoneNumber: c.FormValue("phone_number"),
		Desc:        c.FormValue("desc"),
		CreatedBy:  &parsedCB,
	}

	res, err := M.CreatePublisher(pu)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdatePublisher(c echo.Context) error {

	PublisherId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}


	p := &PublisherForm{
		Name:        c.FormValue("name"),
		Address:     c.FormValue("address"),
		PhoneNumber: c.FormValue("phone_number"),
		Desc:        c.FormValue("desc"),
	}

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedCB, _ := strconv.ParseInt(c.FormValue("updated_by"), 10, 64)

	bo := &M.Publisher{
		Id: int(PublisherId),
		Name:        c.FormValue("name"),
		Address:     c.FormValue("address"),
		PhoneNumber: c.FormValue("phone_number"),
		Desc:        c.FormValue("desc"),
		CreatedBy: &parsedCB,
	}

	res, err := M.UpdatePublisher(bo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedPublisher(c echo.Context) error {

	PublisherId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeletePublisher(PublisherId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindPublisher(c echo.Context) error {

	PublisherId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindPublisher(PublisherId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}
