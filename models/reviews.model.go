package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Reviews struct {
	Id        int64   `db:"id" json:"id"`
	UserId    int64   `db:"user_id" json:"user_id"`
	BookId    int64   `db:"book_id" json:"book_id"`
	Review    string  `db:"review" json:"review"`
	Rating    float32 `db:"rating" json:"rating"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt string  `db:"updated_at" json:"updated_at"`
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
			&obj.UserId,
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