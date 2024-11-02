package models

import (
	"net/http"
	A "perpus_api/db"
	"regexp"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Book struct {
	Id             int            `json:"id" db:"id"`
	ISBN           string         `json:"ISBN" db:"ISBN"`
	Title          string         `json:"title" db:"title"`
	Lang           string         `json:"lang" db:"lang"`
	NumOfPages     *int           `json:"num_of_pages" db:"num_of_pages"`
	Price          *int           `json:"price" db:"price"`
	Desc           *string        `json:"desc" db:"desc"`
	CreatedAt      *string        `db:"created_at" json:"created_at"`
	UpdatedAt      *string        `db:"updated_at" json:"updated_at"`
	DeletedAt      *string        `db:"deleted_at" json:"deleted_at"`
	CreatedBy      *int64         `db:"created_by" json:"created_by"`
	UpdatedBy      *int64         `db:"updated_by" json:"updated_by"`
	DeletedBy      *int64         `db:"deleted_by" json:"deleted_by"`
	PublisherId    *int64         `db:"publisher_id" json:"publisher_id"`
	IsEnabled      bool           `json:"is_enabled" db:"is_enabled"`
	IsOnline       bool           `json:"is_online" db:"is_online"`
	Stock          int            `db:"stock" json:"stock"`
	Authors        []string       `json:"authors,omitempty"`
	Publishers     string         `json:"publisher,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Categories     []string       `json:"categories,omitempty"`
	AuthorsId      []int64        `json:"authors_id,omitempty"`
	TagsId         []int64        `json:"tags_id,omitempty"`
	CategoriesId   []int64        `json:"categories_id,omitempty"`
	Books          []BookDetails  `json:"books,omitempty"`
	MediaLoc       *string        `json:"media_loc,omitempty"`
	ReviewsWithAvg AvgWithReviews `json:"reviews,omitempty"`
}

type BookSmallView struct {
	Id            int     `json:"id" db:"id"`
	Title         string  `json:"title" db:"title"`
	Desc          *string `json:"desc" db:"desc"`
	IsOnline      bool    `json:"is_online" db:"is_online"`
	MediaLoc      *string `json:"media_loc,omitempty"`
	BookDetailsId *int    `json:"book_detail_id,omitempty"`
}

type BookDetails struct {
	Id           int64  `json:"id"`
	Title        string `json:"title" db:"title"`
	SerialNumber string `json:"sn"`
	Condition    string `json:"condition"`
	Price        *int    `json:"price"`
	Status       string `json:"status"`
}

var BookLock = sync.Mutex{}

func GetAllBookDetailsNotBorrowedOrRemoved() (ResponseMultiple, error) {
	var obj BookDetails
	var arrobj []BookDetails
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT 
			t1.id, 
			t2.title, 
			t1.serial_number, 
			t1.condition, 
			t1.status 
		FROM public.book_details t1 
		inner join public.books t2 on t2.id = t1.book_id 
		where t1.status != 'REMOVED' AND t1.status != 'BORROWED';
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
			&obj.Title,
			&obj.SerialNumber,
			&obj.Condition,
			&obj.Status,
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
	res.Msg = "Books details founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func GetAllBook() (ResponseMultiple, error) {
	var obj Book
	var arrobj []Book
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
	SELECT 
		b.*,
		p.name AS publisher,
		array_agg(DISTINCT a.name) AS authors,
		array_agg(DISTINCT COALESCE(t.name, 'No Tags')) AS tags,
		array_agg(DISTINCT c.name) AS categories
	FROM 
		books b
		LEFT JOIN publishers p ON b.publisher_id = p.id
		LEFT JOIN author_book ab ON b.id = ab.book_id
		LEFT JOIN authors a ON ab.author_id = a.id
		LEFT JOIN book_tag bt ON b.id = bt.book_id
		LEFT JOIN tags t ON bt.tag_id = t.id
		LEFT JOIN book_category bc ON b.id = bc.book_id
		LEFT JOIN categories c ON bc.category_id = c.id
	WHERE 
		b.deleted_at IS NULL
	GROUP BY 
		b.id, p.name;
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
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.PublisherId,
			&obj.Desc,
			&obj.IsOnline,
			&obj.IsEnabled,
			&obj.Stock,
			&obj.Publishers,
			pq.Array(&obj.Authors),
			pq.Array(&obj.Tags),
			pq.Array(&obj.Categories),
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

func FindBookWithBookDetails(id int64) (Response, error) {
	var obj Book
	var res Response

	con := A.GetDB()

	sql := `
	SELECT 
		b.*,
		p.name AS publisher,
		array_agg(DISTINCT a.name) AS authors,
		array_agg(DISTINCT COALESCE(t.name, 'No Tags')) AS tags,
		array_agg(DISTINCT c.name) AS categories
	FROM 
		books b
		LEFT JOIN publishers p ON b.publisher_id = p.id
		LEFT JOIN author_book ab ON b.id = ab.book_id
		LEFT JOIN authors a ON ab.author_id = a.id
		LEFT JOIN book_tag bt ON b.id = bt.book_id
		LEFT JOIN tags t ON bt.tag_id = t.id
		LEFT JOIN book_category bc ON b.id = bc.book_id
		LEFT JOIN categories c ON bc.category_id = c.id
	WHERE 
		b.deleted_at IS NULL AND
		b.id = $1
	GROUP BY 
		b.id, p.name;
	`

	findMedia := `SELECT location FROM medias where model_id = $1 and model_name = 'book'`

	err := con.QueryRow(sql, id).Scan(
		&obj.Id,
		&obj.ISBN,
		&obj.Title,
		&obj.Lang,
		&obj.NumOfPages,
		&obj.Price,
		&obj.CreatedAt,
		&obj.UpdatedAt,
		&obj.CreatedBy,
		&obj.UpdatedBy,
		&obj.DeletedAt,
		&obj.DeletedBy,
		&obj.PublisherId,
		&obj.Desc,
		&obj.IsOnline,
		&obj.IsEnabled,
		&obj.Stock,
		&obj.Publishers,
		pq.Array(&obj.Authors),
		pq.Array(&obj.Tags),
		pq.Array(&obj.Categories),
	)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	_ = con.QueryRow(findMedia, id).Scan(&obj.MediaLoc)

	sql = `
	SELECT 
		t1.id, 
		t2.title, 
		t1.serial_number, 
		t1.condition, 
		t1.status 
	FROM public.book_details t1 
		inner join public.books t2 on t2.id = t1.book_id 
	where t1.book_id = $1;
	`

	rows, err := con.Query(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	defer rows.Close()

	var book BookDetails
	var books []BookDetails

	for rows.Next() {
		err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.SerialNumber,
			&book.Condition,
			&book.Status,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		books = append(books, book)
	}

	obj.Books = books

	rev, err := GetAllReviewsBookIdObj(int(id))

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	avg, err := GetBookAvgRating(int(id))

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	datas := AvgWithReviews{ Avg: avg, Reviews: rev}

	obj.ReviewsWithAvg = datas

	res.Status = http.StatusOK
	res.Msg = "Books founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func FindBook(id int64) (Response, error) {
	var obj Book
	var res Response

	con := A.GetDB()

	sql := `
	SELECT 
		b.*,
		p.name,
		array_agg(DISTINCT a.id) AS authors,
		array_agg(DISTINCT COALESCE(t.name, 'No Tags')) AS tags,
		array_agg(DISTINCT c.id) AS categories
	FROM 
		books b
		LEFT JOIN publishers p ON b.publisher_id = p.id
		LEFT JOIN author_book ab ON b.id = ab.book_id
		LEFT JOIN authors a ON ab.author_id = a.id
		LEFT JOIN book_tag bt ON b.id = bt.book_id
		LEFT JOIN tags t ON bt.tag_id = t.id
		LEFT JOIN book_category bc ON b.id = bc.book_id
		LEFT JOIN categories c ON bc.category_id = c.id
	WHERE 
		b.deleted_at IS NULL AND
		b.id = $1
	GROUP BY 
		b.id, p.name;
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
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.PublisherId,
			&obj.Desc,
			&obj.IsOnline,
			&obj.IsEnabled,
			&obj.Stock,
			&obj.Publishers,
			pq.Array(&obj.AuthorsId),
			pq.Array(&obj.Tags),
			pq.Array(&obj.CategoriesId),
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

func CreateBook(book *Book, authorId []int64, categoryId []int64, tags []string) (*int, Response, error) {

	BookLock.Lock()
	defer BookLock.Unlock()

	var res Response

	con := A.GetDB()

	// tx, err := con.Begin()

	sql := `
		INSERT INTO public.books(
		"ISBN", title, lang, num_of_pages, price, created_at, created_by, publisher_id, "desc", is_online, is_enabled)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7, $8, $9, $10) RETURNING id;
	`

	var id *int

	err := con.QueryRow(
		sql, book.ISBN, book.Title, book.Lang, book.NumOfPages, book.Price, book.CreatedBy, book.PublisherId, book.Desc, book.IsOnline, book.IsEnabled,
	).Scan(&id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return nil, res, err
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
			return nil, res, err
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
			return nil, res, err
		}
	}

	err = CreateBookTags(tags, *id, false)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return nil, res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Book created successfully"
	res.Success = true
	res.Data = &echo.Map{
		"id": id,
	}

	return nil, res, nil
}

func CreateBookTags(tags []string, bookId int, isEdit bool) error {

	con := A.GetDB()

	// tx, err := con.Begin()

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

	// tx, err := con.Begin()

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

	// tx, err := con.Begin()

	sql := `
		UPDATE public.books
			SET "ISBN" = $1, title = $2, lang = $3, num_of_pages = $4, price = $5, updated_at = NOW(), updated_by = $6, publisher_id = $7, "desc" = $8, is_online = $9, is_enabled = $10
		WHERE id= $11
	`

	_, err := con.Exec(sql, book.ISBN, book.Title, book.Lang, book.NumOfPages, book.Price, book.UpdatedBy, book.PublisherId, book.Desc, book.IsOnline, book.IsEnabled, book.Id)

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
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.PublisherId,
			&obj.Desc,
			&obj.IsOnline,
			&obj.IsEnabled,
			&obj.Stock,
		)

		if err != nil {
			return nil, err
		}

	}

	return &obj, nil
}

func UpdateBookStatusAndPagesNum(bookId int, num int) error {
	BookLock.Lock()
	defer BookLock.Unlock()

	con := A.GetDB()

	sql := `
		UPDATE public.books
			SET is_online = $1, num_of_pages = $2
		WHERE id= $3
	`

	_, err := con.Exec(sql, true, num, bookId)

	return err
}

// func UpdateBookStock(id int) error {

// 	con := A.GetDB()

// 	sql := `UPDATE books SET stock = (SELECT COUNT(id) FROM book_details WHERE status = 'STORED' AND book_id = $1) WHERE id = $1;`

// 	_, err := con.Exec(sql, id)

// 	return err

// }

func GetBookForRecom() (*[]BookSmallView, error) {

	var obj BookSmallView
	var arrobj []BookSmallView

	con := A.GetDB()

	sql := `
	SELECT t1.id, t1.title, t1.desc, t1.is_online, t2.location as location FROM books t1 
		inner join medias t2 on t1.id = t2.model_id
	where
		t2.model_name = 'book' and t1.stock > 0 and t1.is_enabled = true
	ORDER BY RANDOM()
	`

	rows, err := con.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Title,
			&obj.Desc,
			&obj.IsOnline,
			&obj.MediaLoc,
		)

		if err != nil {
			return nil, err
		}

		arrobj = append(arrobj, obj)

	}

	return &arrobj, nil

}
