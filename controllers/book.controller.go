package controllers

import (
	"net/http"
	M "perpus_api/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookForm struct {
	ISBN        string `json:"ISBN" validate:"required"`
	Title       string `json:"address" validate:"required"`
	Lang        string `json:"phone_number" validate:"required"`
	NumOfPages  int    `json:"num_of_pages" validate:"required"`
	AuthorId    int64    `json:"author_id" validate:"required"`
	PublisherId int64    `json:"publisher_id" validate:"required"`
}

func GetAllBook(c echo.Context) error {
	res, err := M.GetAllBook()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateBook(c echo.Context) error {

	parsed, err := strconv.Atoi(c.FormValue("num_of_pages"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	b := &BookForm{
		ISBN:       c.FormValue("isbn"),
		Title:      c.FormValue("title"),
		Lang:       c.FormValue("lang"),
		NumOfPages: parsed,
	}

	if err := c.Validate(b); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedP, _ := strconv.Atoi(c.FormValue("price"))

	parsedCB, _ := strconv.ParseInt(c.FormValue("created_by"), 10, 64)

	bo := &M.Book{
		ISBN:       b.ISBN,
		Title:      b.Title,
		Lang:       b.Lang,
		NumOfPages: b.NumOfPages,
		Price:      &parsedP,
		Desc:       c.FormValue("desc"),
		Sypnosis:   c.FormValue("sypnosis"),
		CreatedBy:  &parsedCB,
	}

	res, err := M.CreateBook(bo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateBook(c echo.Context) error {

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsed, err := strconv.Atoi(c.FormValue("num_of_pages"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	authorId, err := strconv.ParseInt(c.FormValue("author_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	publisherId, err := strconv.ParseInt(c.FormValue("publisher_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	b := &BookForm{
		ISBN:        c.FormValue("isbn"),
		Title:       c.FormValue("title"),
		Lang:        c.FormValue("lang"),
		NumOfPages:  parsed,
		AuthorId:    authorId,
		PublisherId: publisherId,
	}

	if err := c.Validate(b); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	parsedP, _ := strconv.Atoi(c.FormValue("price"))

	parsedCB, _ := strconv.ParseInt(c.FormValue("updated_by"), 10, 64)

	bo := &M.Book{
		Id:          int(bookId),
		ISBN:        b.ISBN,
		Title:       b.Title,
		Lang:        b.Lang,
		NumOfPages:  b.NumOfPages,
		Price:       &parsedP,
		Desc:        c.FormValue("desc"),
		Sypnosis:    c.FormValue("sypnosis"),
		CreatedBy:   &parsedCB,
		AuthorId:    &b.AuthorId,
		PublisherId: &b.PublisherId,
	}

	res, err := M.UpdateBook(bo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedBook(c echo.Context) error {

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeleteBook(bookId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindBook(c echo.Context) error {

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindBook(bookId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}
