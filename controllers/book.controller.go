package controllers

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"path/filepath"
	M "perpus_api/models"
	"strconv"

	"github.com/karmdip-mi/go-fitz"
	"github.com/labstack/echo/v4"
)

type BookForm struct {
	ISBN        string    `json:"ISBN" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Lang        string    `json:"lang" validate:"required"`
	NumOfPages  *int      `json:"num_of_pages"`
	AuthorId    []int64   `json:"author_id" validate:"required"`
	PublisherId int64     `json:"publisher_id" validate:"required"`
	CategoryId  []int64   `json:"category_id" validate:"required"`
	Tags        *[]string `json:"tags"`
	Price       *int      `json:"price"`
	Desc        *string   `json:"desc"`
	IsEnabled   bool      `json:"is_enabled"`
	IsOnline    bool      `json:"is_online"`
}

func GetAllBook(c echo.Context) error {
	res, err := M.GetAllBook()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateBook(c echo.Context) error {

	requestBody := new(BookForm)

	if err := c.Bind(requestBody); err != nil {
		return err
	}

	if err := c.Validate(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	bo := &M.Book{
		ISBN:       requestBody.ISBN,
		Title:      requestBody.Title,
		Lang:       requestBody.Lang,
		NumOfPages: requestBody.NumOfPages,
		Price:      requestBody.Price,
		Desc:       requestBody.Desc,
		CreatedBy:  &parsedCB,
		PublisherId: &requestBody.PublisherId,
		IsEnabled: requestBody.IsEnabled,
		IsOnline: requestBody.IsOnline,
	}

	_, res, err := M.CreateBook(bo, requestBody.AuthorId, requestBody.CategoryId, *requestBody.Tags)

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

	requestBody := new(BookForm)

	if err := c.Bind(requestBody); err != nil {
		return err
	}

	if err := c.Validate(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	bo := &M.Book{
		Id:          int(bookId),
		ISBN:        requestBody.ISBN,
		Title:       requestBody.Title,
		Lang:        requestBody.Lang,
		NumOfPages:  requestBody.NumOfPages,
		Price:       requestBody.Price,
		Desc:        requestBody.Desc,
		PublisherId: &requestBody.PublisherId,
		IsEnabled: requestBody.IsEnabled,
		IsOnline: requestBody.IsOnline,
		UpdatedBy: &parsedCB,
	}

	res, err := M.UpdateBook(bo, requestBody.AuthorId, requestBody.CategoryId, *requestBody.Tags)

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

func UploadBookPdfToImage(c echo.Context) error {

	var res M.ResponseNoData

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	bookTitle := c.FormValue("title")

	re := regexp.MustCompile(`^\s+|\s+`)
	cleanedStr :=  strings.Replace(strings.ToLower(re.ReplaceAllString(bookTitle, ""))," ", "-", -1)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	file, err := c.FormFile("file")

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	src, err := file.Open()

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	defer src.Close()

	location := "uploads/books/" + cleanedStr + strconv.Itoa(int(bookId))

	err = os.MkdirAll(location, 0777)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	err = os.MkdirAll(location+"/pdf", 0777)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	err = os.MkdirAll(location+"/images", 0777)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	dst, err := os.Create(location + "/pdf/" + file.Filename)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {

		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)

	}

	doc, err := fitz.New(location + "/pdf/" + file.Filename)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	defer doc.Close()

	for i := 0; i < doc.NumPage(); i++ {

		img, err := doc.Image(i)

		if err != nil {
			res.Status = http.StatusBadRequest
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusBadRequest, res)
		}

		f, err := os.Create(filepath.Join(location+"/images/", fmt.Sprintf("image-%05d.jpg", i)))

		if err != nil {
			res.Status = http.StatusBadRequest
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusBadRequest, res)
		}

		err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})

		if err != nil {
			res.Status = http.StatusBadRequest
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusBadRequest, res)
		}

		f.Close()
	}

	err = M.CreateMedia(location+"/images/", bookId, "book", "images")

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	err = M.UpdateBookStatusAndPagesNum(int(bookId), doc.NumPage())

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Status = http.StatusOK
	res.Msg = "Berhasil disimpan!"
	res.Success = true

	return c.JSON(http.StatusOK, res)

}
