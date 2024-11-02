package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Reviews struct {
	Id         int64   `db:"id" json:"id"`
	MemberId   int64   `db:"member_id" json:"member_id"`
	MemberName string  `json:"member_name"`
	BookId     int64   `db:"book_id" json:"book_id" validate:"required"`
	Review     string  `db:"review" json:"review" validate:"required"`
	Rating     float32 `db:"rating" json:"rating" validate:"required"`
	CreatedAt  string  `db:"created_at" json:"created_at"`
	UpdatedAt  *string `db:"updated_at" json:"updated_at,omitempty"`
}

var ReviewsLock = sync.Mutex{}

func GetAllReviews() (ResponseMultiple, error) {
	var obj Reviews
	var arrobj []Reviews
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT * FROM reviews WHERE deleted_at IS NULL;
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
			&obj.MemberId,
			&obj.BookId,
			&obj.Review,
			&obj.Rating,
			&obj.CreatedAt,
			&obj.UpdatedAt,
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
	res.Msg = "Reviews founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

type AvgWithReviews struct {
	Avg     *float32   `json:"avg"`
	Reviews *[]Reviews `json:"reviews"`
}

func GetBookAvgRating(bookId int) (*float32, error) {
	var avg *float32

	con := A.GetDB()

	sql := `
		SELECT AVG(rating) FROM reviews WHERE book_id = $1;
	`

	err := con.QueryRow(sql, bookId).Scan(&avg)

	if err != nil {
		return nil, err
	}

	return avg, err

}

func GetAllReviewsBookIdObj(bookId int) (*[]Reviews, error) {
	var obj Reviews
	var arrobj []Reviews

	con := A.GetDB()

	sql := `
		SELECT t1.*, t2.username FROM reviews t1 inner join members t2 on t1.member_id = t2.id WHERE book_id = $1;
	`

	rows, err := con.Query(sql, bookId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.MemberId,
			&obj.BookId,
			&obj.Review,
			&obj.Rating,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.MemberName,
		)

		if err != nil {
			return nil, err
		}

		arrobj = append(arrobj, obj)
	}

	return &arrobj, nil
}

func StoreReviews(r Reviews, claims JwtCustomClaims) (ResponseNoData, error) {

	var res ResponseNoData

	con := A.GetDB()

	query := `
		INSERT INTO reviews (member_id, book_id, review, rating, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := con.Exec(
		query,
		claims.ID,
		r.BookId,
		r.Review,
		r.Rating,
	)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Reviews created!"
	res.Success = true

	return res, nil

}
