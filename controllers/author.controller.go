package controllers

import (
	"net/http"
	M "perpus_api/models"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthorForm struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc" validate:"required"`
}

func GetAllAuthor(c echo.Context) error {
	res, err := M.GetAllAuthor()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateAuthor(c echo.Context) error {

	p := &AuthorForm{
		Name: c.FormValue("name"),
		Desc: c.FormValue("desc"),
	}

	user, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return c.JSON(http.StatusBadRequest, "JWT token missing or invalid")
	}

	claims := user.Claims.(*M.JwtCustomClaims)

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedCB := int64(claims.ID)

	pu := &M.Author{
		Name:      c.FormValue("name"),
		Desc:      c.FormValue("desc"),
		CreatedBy: &parsedCB,
	}

	res, err := M.CreateAuthor(pu)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateAuthor(c echo.Context) error {

	AuthorId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	p := &AuthorForm{
		Name: c.FormValue("name"),
		Desc: c.FormValue("desc"),
	}

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedCB, _ := strconv.ParseInt(c.FormValue("updated_by"), 10, 64)

	bo := &M.Author{
		Id:        int(AuthorId),
		Name:      c.FormValue("name"),
		Desc:      c.FormValue("desc"),
		CreatedBy: &parsedCB,
	}

	res, err := M.UpdateAuthor(bo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedAuthor(c echo.Context) error {

	AuthorId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeleteAuthor(AuthorId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindAuthor(c echo.Context) error {

	AuthorId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindAuthor(AuthorId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}