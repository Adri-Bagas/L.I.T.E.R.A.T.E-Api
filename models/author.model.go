package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Author struct {
	Id        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Desc      string  `json:"desc" db:"desc"`
	CreatedAt *string `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at"`
	CreatedBy *int64  `db:"created_by" json:"created_by"`
	UpdatedBy *int64  `db:"updated_by" json:"updated_by"`
	DeletedBy *int64  `db:"deleted_by" json:"deleted_by"`
}

var AuthorLock = sync.Mutex{}

func GetAllAuthor() (ResponseMultiple, error) {
	var obj Author
	var arrobj []Author
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT * FROM authors WHERE deleted_at IS NULL;
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
	res.Msg = "Authors founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func GetAllAuthorObj() ([]Author, error) {
	var obj Author
	var arrobj []Author

	con := A.GetDB()

	sql := `
		SELECT * FROM authors WHERE deleted_at IS NULL;
	`

	rows, err := con.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Name,
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

		arrobj = append(arrobj, obj)
	}


	return arrobj, nil
}

func FindAuthor(id int64) (Response, error) {
	var obj Author
	var res Response

	con := A.GetDB()

	sql := `
		SELECT * FROM authors WHERE id = $1;
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
	res.Msg = "Author founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func CreateAuthor(author *Author) (ResponseNoData, error) {

	AuthorLock.Lock()
	defer AuthorLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		INSERT INTO public.authors(
			name, "desc", created_at, created_by)
		VALUES ($1, $2, NOW(), $3);
	`

	_, err := con.Exec(sql, author.Name, author.Desc, author.CreatedBy)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Author created successfully"
	res.Success = true

	return res, nil
}

func DeleteAuthor(id int64) (ResponseNoData, error) {
	AuthorLock.Lock()
	defer AuthorLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.authors SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Author trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateAuthor(author *Author) (ResponseNoData, error) {
	AuthorLock.Lock()
	defer AuthorLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.authors
			SET name = $1, "desc" = $2, updated_at = NOW(), updated_by = $3
		WHERE id = $4;
	`

	_, err := con.Exec(sql, author.Name, author.Desc, author.UpdatedBy, author.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Author updated successfully"
	res.Success = true

	return res, nil
}

func WhereAuthor(col string, val string) (*Author, error) {

	var obj Author

	con := A.GetDB()

	sql := "SELECT * FROM author WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Name,
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
