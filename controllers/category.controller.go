package controllers

import (
	"net/http"
	M "perpus_api/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryForm struct {
	Name string `json:"name" validate:"required"`
	Desc *string `json:"desc"`
}

func GetAllCategory(c echo.Context) error {
	res, err := M.GetAllCategory()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func GetAllCategoryIdName(c echo.Context) error {
	category, err := M.GetAllCategoryObj()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"err": err.Error(),
		})
	}

	datas := make(map[int]string)
	
	for _, val := range category {
		datas[val.Id] = val.Name
	}

	return c.JSON(http.StatusOK, datas)
}

func CreateCategory(c echo.Context) error {

	desc := c.FormValue("desc")

	p := &CategoryForm{
		Name: c.FormValue("name"),
		Desc: &desc,
	}

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	ca := &M.Category{
		Name: p.Name,
		Desc: p.Desc,
		CreatedBy: &parsedCB,
	}

	res, err := M.CreateCategory(ca)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateCategory(c echo.Context) error {
	CategoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	desc := c.FormValue("desc")

	p := &CategoryForm{
		Name: c.FormValue("name"),
		Desc: &desc,
	}

	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedUB := int64(claims.ID)

	ca := &M.Category{
		Id: int(CategoryId),
		Name: p.Name,
		Desc: p.Desc,
		UpdatedBy: &parsedUB,
	}

	res, err := M.UpdateCategory(ca)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedCategory(c echo.Context) error {

	CategoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeleteCategory(&CategoryId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindCategory(c echo.Context) error {

	CategoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindCategory(CategoryId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}