package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Book struct {
	Id          int     `json:"id" db:"id"`
	ISBN        string  `json:"ISBN" db:"ISBN"`
	Title       string  `json:"address" db:"address"`
	Lang        string  `json:"phone_number" db:"phone_number"`
	NumOfPages  int     `json:"num_of_pages" db:"num_of_pages"`
	Price       *int    `json:"price" db:"price"`
	Desc        string  `json:"desc" db:"desc"`
	Sypnosis    string  `json:"sypnosis" db:"sypnosis"`
	CreatedAt   *string `db:"created_at" json:"created_at"`
	UpdatedAt   *string `db:"updated_at" json:"updated_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at"`
	CreatedBy   *int64  `db:"created_by" json:"created_by"`
	UpdatedBy   *int64  `db:"updated_by" json:"updated_by"`
	DeletedBy   *int64  `db:"deleted_by" json:"deleted_by"`
	AuthorId    *int64  `db:"author_id" json:"author_id"`
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
			&obj.Sypnosis,
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
			&obj.Sypnosis,
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

func CreateBook(book *Book) (ResponseNoData, error) {

	BookLock.Lock()
	defer BookLock.Unlock()

	var res ResponseNoData

	var desc *string

	var syp *string

	if book.Desc == "" {desc = nil} else {desc = &book.Desc}

	if book.Sypnosis == "" {syp = nil} else {syp = &book.Sypnosis}

	con := A.GetDB()

	sql := `
		INSERT INTO public.books(
		"ISBN", title, lang, num_of_pages, price, created_at, created_by, author_id, publisher_id, "desc", sypnosis)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9, $10);
	`

	_, err := con.Exec(sql, book.ISBN, book.Title, book.Lang, book.Title, book.NumOfPages, book.Price, book.CreatedBy, book.AuthorId, book.PublisherId, desc, syp)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Book created successfully"
	res.Success = true

	return res, nil
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
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Book trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateBook(book *Book) (ResponseNoData, error) {
	BookLock.Lock()
	defer BookLock.Unlock()

	var res ResponseNoData

	var desc *string

	var syp *string

	if book.Desc == "" {desc = nil} else {desc = &book.Desc}

	if book.Sypnosis == "" {syp = nil} else {syp = &book.Sypnosis}

	con := A.GetDB()

	sql := `
		UPDATE public.books
			SET "ISBN" = $1, title = $2, lang = $3, num_of_pages = $4, price = $5, updated_at = NOW(), updated_by = $6, author_id = $7, publisher_id = $8, "desc" = $9, sypnosis = $10
		WHERE id= $12
	`

	_, err := con.Exec(sql, book.Id, book.Title, book.Lang, book.NumOfPages, book.Price, book.UpdatedBy, book.AuthorId, book.PublisherId, desc, syp)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Book updated successfully"
	res.Success = true

	return res, nil
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
			&obj.Sypnosis,
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
