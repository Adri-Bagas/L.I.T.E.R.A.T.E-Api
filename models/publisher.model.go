package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Publisher struct {
	Id          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Address     string  `json:"address" db:"address"`
	PhoneNumber string  `json:"phone_number" db:"phone_number"`
	Desc        string  `json:"desc" db:"desc"`
	CreatedAt   *string `db:"created_at" json:"created_at"`
	UpdatedAt   *string `db:"updated_at" json:"updated_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at"`
	CreatedBy   *int64  `db:"created_by" json:"created_by"`
	UpdatedBy   *int64  `db:"updated_by" json:"updated_by"`
	DeletedBy   *int64  `db:"deleted_by" json:"deleted_by"`
}

var PublisherLock = sync.Mutex{}

func GetAllPublisher() (ResponseMultiple, error) {
	var obj Publisher
	var arrobj []Publisher
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT * FROM publishers WHERE deleted_at IS NULL;
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
			&obj.Name,
			&obj.PhoneNumber,
			&obj.Desc,
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
	res.Msg = "Publishers founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func FindPublisher(id int64) (Response, error) {
	var obj Publisher
	var res Response

	con := A.GetDB()

	sql := `
		SELECT * FROM publishers WHERE id = $1;
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
			&obj.Name,
			&obj.PhoneNumber,
			&obj.Desc,
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
	res.Msg = "Publisher founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func CreatePublisher(publisher *Publisher) (ResponseNoData, error) {

	PublisherLock.Lock()
	defer PublisherLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		INSERT INTO public.publishers(
			name, address, phone_number, "desc", created_at, created_by)
		VALUES ($1, $2, $3, $4, NOW(), $5);
	`

	_, err := con.Exec(sql, publisher.Name, publisher.Address, publisher.PhoneNumber, publisher.Desc, publisher.CreatedBy)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Publisher created successfully"
	res.Success = true

	return res, nil
}

func DeletePublisher(id int64) (ResponseNoData, error) {
	PublisherLock.Lock()
	defer PublisherLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.publishers SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Publisher trashed successfully"
	res.Success = true

	return res, nil
}

func UpdatePublisher(publisher *Publisher) (ResponseNoData, error) {
	PublisherLock.Lock()
	defer PublisherLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.publishers
			SET name = $1, address = $2, phone_number = $3, "desc" = $4, updated_at = NOW(), updated_by=$5
		WHERE id = $6;
	`

	_, err := con.Exec(sql, publisher.Name, publisher.Address, publisher.Desc, publisher.UpdatedBy, publisher.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "Publisher updated successfully"
	res.Success = true

	return res, nil
}

func WherePublisher(col string, val string) (*Publisher, error) {

	var obj Publisher

	con := A.GetDB()

	sql := "SELECT * FROM publishers WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Name,
			&obj.PhoneNumber,
			&obj.Desc,
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
