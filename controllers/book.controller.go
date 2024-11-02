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
	A "perpus_api/db"
	H "perpus_api/helpers"
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

func GetAllBookDetailsNotBorrowedOrRemoved(c echo.Context) error {
	res, err := M.GetAllBookDetailsNotBorrowedOrRemoved()

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
		ISBN:        requestBody.ISBN,
		Title:       requestBody.Title,
		Lang:        requestBody.Lang,
		NumOfPages:  requestBody.NumOfPages,
		Price:       requestBody.Price,
		Desc:        requestBody.Desc,
		CreatedBy:   &parsedCB,
		PublisherId: &requestBody.PublisherId,
		IsEnabled:   requestBody.IsEnabled,
		IsOnline:    requestBody.IsOnline,
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
		IsEnabled:   requestBody.IsEnabled,
		IsOnline:    requestBody.IsOnline,
		UpdatedBy:   &parsedCB,
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

func FindBookWithBookDetails(c echo.Context) error {

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindBookWithBookDetails(bookId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func GetBookRecom(c echo.Context) error {
	books, err := M.GetBookForRecom()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	return c.JSON(
		http.StatusOK,
		M.ResponseMultiple{
			Status:  http.StatusOK,
			Msg:     "books founded!",
			Success: true,
			Datas:   books,
		},
	)
}

func UploadBookPdfToImage(c echo.Context) error {

	var res M.ResponseNoData

	bookId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	bookTitle := c.FormValue("title")

	re := regexp.MustCompile(`^\s+|\s+`)
	cleanedStr := strings.Replace(strings.ToLower(re.ReplaceAllString(bookTitle, "")), " ", "-", -1)

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

func GetBookReaded(c echo.Context) error {
	var res M.ResponseMultiple

	var obj M.BookSmallView
	var arrobj []M.BookSmallView

	con := A.GetDB()

	sql := `
		SELECT 
		t1.book_details_id, t2.id, t2.title, t2.desc, t2.is_online, t3.location as location 
		from member_access t1 inner join books t2 on t2.id = t1.book_id
		inner join medias t3 on t2.id = t3.model_id where t3.model_name = 'book' and t1.member_id = $1
	`

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	rows, err := con.Query(sql, parsedCB)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	for rows.Next() {
		err := rows.Scan(
			&obj.BookDetailsId,
			&obj.Id,
			&obj.Title,
			&obj.Desc,
			&obj.IsOnline,
			&obj.MediaLoc,
		)

		if err != nil {
			res.Status = http.StatusBadRequest
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusBadRequest, res)
		}

		arrobj = append(arrobj, obj)
	}

	res.Msg = "books founded!"
	res.Status = http.StatusOK
	res.Success = true
	res.Datas = arrobj

	return c.JSON(
		http.StatusOK,
		res,
	)
}

func GetBookCollections(c echo.Context) error {
	var res M.ResponseMultiple

	var obj M.BookSmallView
	var arrobj []M.BookSmallView

	con := A.GetDB()

	sql := `
		SELECT 
		t2.id, t2.title, t2.desc, t2.is_online, t3.location as location 
		from member_collections t1 inner join books t2 on t2.id = t1.book_id
		inner join medias t3 on t2.id = t3.model_id where t3.model_name = 'book' and t1.member_id = $1
	`

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	rows, err := con.Query(sql, parsedCB)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Title,
			&obj.Desc,
			&obj.IsOnline,
			&obj.MediaLoc,
		)

		if err != nil {
			res.Status = http.StatusBadRequest
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusBadRequest, res)
		}

		arrobj = append(arrobj, obj)
	}

	res.Msg = "books founded!"
	res.Status = http.StatusOK
	res.Success = true
	res.Datas = arrobj

	return c.JSON(
		http.StatusOK,
		res,
	)
}

func GenBookAccessKey(c echo.Context) error {

	var res M.Response

	bookDetails, err := strconv.Atoi(c.FormValue("details_id"))

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	con := A.GetDB()

	claims, errm := M.GetUserDataByJWT(c)

	if errm != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": *errm})
	}

	parsedCB := int64(claims.ID)

	sql := `
	SELECT id
		FROM public.member_access where member_id = $1 AND book_details_id = $2;
	`

	var id *int

	err = con.QueryRow(sql, parsedCB, bookDetails).Scan(&id)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	if id == nil {
		res.Status = http.StatusBadRequest
		res.Msg = "Failed to auth"
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	genString, err := H.GenerateRandomString(24)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	sql = `
		Update member_access set access_key = $1 where id = $2;
	`

	_, err = con.Exec(sql, genString, id)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create Key!"
	res.Success = true
	res.Data = echo.Map{
		"access_key": genString,
	}

	return c.JSON(http.StatusOK, res)
}

type Page struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	URI    string `json:"uri"`
}

func generateArray(number int, location string) [][]Page {
	var result [][]Page

	if number > 0 {
        uri := fmt.Sprintf("http://localhost:1323/"+ location +"image-%05d.jpg", 0)
        result = append(result, []Page{{Width: 800, Height: 1200, URI: uri}})
    }

	for i := 1; i < number; i += 2 {
		var images []Page

		for j := i; j <= i+1 && j < number; j++ {
			uri := fmt.Sprintf("http://localhost:1323/"+ location +"image-%05d.jpg", j)
			images = append(images, Page{Width: 800, Height: 1200, URI: uri})
		}

		result = append(result, images)
	}

	return result
}

func GetBookDataFromAccessKey(c echo.Context) error {
	var res M.Response

	accessKey := c.FormValue("access_key")

	con := A.GetDB()

	var location *string
	var numOfPages *int

	sql := `
		SELECT t3."location", t2.num_of_pages from member_access t1 
			inner join books t2 on t1.book_id = t2.id
			inner join medias t3 on t2.id = t3.model_id
		where t1.access_key = $1 AND t3.model_name = 'book';
	`

	err := con.QueryRow(sql, accessKey).Scan(&location, &numOfPages)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusInternalServerError, res)
	}

	resultsArr := generateArray(*numOfPages, *location) 

	sql = `
		update member_access set access_key = null where access_key = $1;
	`

	_, err = con.Exec(sql, accessKey)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = http.StatusOK
	res.Msg = "Success to create Key!"
	res.Success = true
	res.Data = resultsArr

	return c.JSON(http.StatusOK, res)
}
