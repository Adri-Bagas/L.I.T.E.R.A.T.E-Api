package models

import (
	"net/http"
	A "perpus_api/db"
	"regexp"
	"strings"
	"sync"
)

type Book struct {
	Id          int     `json:"id" db:"id"`
	ISBN        string  `json:"ISBN" db:"ISBN"`
	Title       string  `json:"address" db:"address"`
	Lang        string  `json:"phone_number" db:"phone_number"`
	NumOfPages  int     `json:"num_of_pages" db:"num_of_pages"`
	Price       *int    `json:"price" db:"price"`
	Desc        *string `json:"desc" db:"desc"`
	CreatedAt   *string `db:"created_at" json:"created_at"`
	UpdatedAt   *string `db:"updated_at" json:"updated_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at"`
	CreatedBy   *int64  `db:"created_by" json:"created_by"`
	UpdatedBy   *int64  `db:"updated_by" json:"updated_by"`
	DeletedBy   *int64  `db:"deleted_by" json:"deleted_by"`
	PublisherId *int64  `db:"publisher_id" json:"publisher_id"`
}

var BookLock = sync.Mutex{}

func GetAllBook() (ResponseMultiple, error) {
	var obj Book
	var arrobj []Book
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT * FROM books WHERE deleted_at IS NULL;
	`

	rows, err := con.Query(sql)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.ISBN,
			&obj.Title,
			&obj.Lang,
			&obj.NumOfPages,
			&obj.Price,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedBy,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		arrobj = append(arrobj, obj)
	}

	res.Status = http.StatusOK
	res.Msg = "Books founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func FindBook(id int64) (Response, error) {
	var obj Book
	var res Response

	con := A.GetDB()

	sql := `
		SELECT * FROM books WHERE id = $1;
	`

	rows, err := con.Query(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.ISBN,
			&obj.Title,
			&obj.Lang,
			&obj.NumOfPages,
			&obj.Price,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedBy,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

	}

	res.Status = http.StatusOK
	res.Msg = "Books founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func CreateBook(book *Book, authorId []int64, categoryId []int64, tags []string) (ResponseNoData, error) {

	BookLock.Lock()
	defer BookLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		INSERT INTO public.books(
		"ISBN", title, lang, num_of_pages, price, created_at, created_by, publisher_id, "desc")
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8) RETURNING id;
	`

	var id *int

	err := con.QueryRow(
		sql, book.ISBN, book.Title, book.Lang, book.Title, book.NumOfPages, book.Price, book.CreatedBy, book.PublisherId, book.Desc,
	).Scan(&id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	sql = `
		insert into public.author_book(author_id, book_id) values ($1, $2)
	`

	for _, e := range authorId {
		_, err := con.Exec(sql, e, id)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}
	}

	sql = `
		insert into public.book_category(category_id, book_id) values ($1, $2)
	`

	for _, e := range categoryId {
		_, err := con.Exec(sql, e, id)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}
	}

	err = CreateBookTags(tags, *id, false)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Book created successfully"
	res.Success = true

	return res, nil
}

func CreateBookTags(tags []string, bookId int, isEdit bool) error {

	con := A.GetDB()

	sql := "insert into public.book_tag(book_id, tag_id) values ($1, $2)"

	sqlCreateTag := "insert into public.tags(name) values ($1) RETURNING id"

	sqlCheckTag := "select id from public.tags where name = $1"

	if isEdit {
		sqlDeleteConnection := "delete from public.book_tag where book_id = $1"

		_, err := con.Exec(sqlDeleteConnection, bookId)

		if err != nil {
			return err
		}
	}

	for _, tag := range tags {
		var tagID *int

		re := regexp.MustCompile(`^\s+|\s+`)
		cleanedStr := re.ReplaceAllString(tag, "")

		formattedTagName := strings.Replace(strings.ToLower(cleanedStr), " ", "-", -1)

		con.QueryRow(sqlCheckTag, formattedTagName).Scan(&tagID)

		if tagID == nil {
			err := con.QueryRow(sqlCreateTag, formattedTagName).Scan(&tagID)

			if err != nil {
				return err
			}
		}

		_, err := con.Exec(sql, bookId, tagID)

		if err != nil {
			return err
		}

	}

	return nil

}

func DeleteBook(id int64) (ResponseNoData, error) {
	BookLock.Lock()
	defer BookLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.books SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Book trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateBook(book *Book, authorId []int64, categoryId []int64, tags []string) (ResponseNoData, error) {
	BookLock.Lock()
	defer BookLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.books
			SET "ISBN" = $1, title = $2, lang = $3, num_of_pages = $4, price = $5, updated_at = NOW(), updated_by = $6, publisher_id = $7, "desc" = $8
		WHERE id= $9
	`

	_, err := con.Exec(sql, book.ISBN, book.Title, book.Lang, book.NumOfPages, book.Price, book.UpdatedBy, book.PublisherId, book.Desc, book.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	err = checkAuthorBook(book.Id, authorId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	err = checkCategoryBook(book.Id, categoryId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	err = CreateBookTags(tags, book.Id, true)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Book updated successfully"
	res.Success = true

	return res, nil
}

func checkAuthorBook(bookId int, authorIds []int64) error {
	con := A.GetDB()

	sql := "insert into public.author_book(author_id, book_id) values ($1, $2)"

	sqlDeleteAll := "DELETE FROM public.author_book WHERE book_id = $1;"

	_, err := con.Exec(sqlDeleteAll, bookId)

		if err != nil {
			return err
		}

	for _, author := range authorIds {
		_, err := con.Exec(sql, author, bookId)

		if err != nil {
			return err
		}
	}

	return nil
}

func checkCategoryBook(bookId int, categoryIds []int64) error {
	con := A.GetDB()

	sql := "insert into public.book_category(category_id, book_id) values ($1, $2)"

	sqlDeleteAll := "DELETE FROM public.book_category WHERE book_id = $1;"

	_, err := con.Exec(sqlDeleteAll, bookId)

		if err != nil {
			return err
		}

	for _, category := range categoryIds {
		_, err := con.Exec(sql, category, bookId)

		if err != nil {
			return err
		}
	}

	return nil
}

func WhereBook(col string, val string) (*Book, error) {

	var obj Book

	con := A.GetDB()

	sql := "SELECT * FROM Books WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.ISBN,
			&obj.Title,
			&obj.Lang,
			&obj.NumOfPages,
			&obj.Price,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedBy,
		)

		if err != nil {
			return nil, err
		}

	}

	return &obj, nil
}
