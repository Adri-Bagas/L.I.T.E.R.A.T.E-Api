package models

import (
	"net/http"
	A "perpus_api/db"
	"sync"
)

type Category struct {
	Id        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	CreatedAt *string `db:"created_at" json:"created_at"`
	CreatedBy *int64  `db:"created_by" json:"created_by"`
	UpdatedAt *string `db:"updated_at" json:"updated_at"`
	UpdatedBy *int64  `db:"updated_by" json:"updated_by"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at"`
}

var CategoryLock = sync.Mutex{}

func GetAllCategory() (ResponseMultiple, error) {
	var obj Category
	var arrobj []Category
	var res ResponseMultiple

	con := A.GetDB()

	sql := `
		SELECT * FROM categories WHERE deleted_at IS NULL;
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
			&obj.CreatedAt,
			&obj.CreatedBy,
			&obj.UpdatedAt,
			&obj.UpdatedBy,
			&obj.DeletedAt,
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
	res.Msg = "Categories founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func FindCategory(id int64) (Response, error) {
	var obj Category
	var res Response

	con := A.GetDB()

	sql := `
		SELECT * FROM categories WHERE id = $1;
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
			&obj.CreatedAt,
			&obj.CreatedBy,
			&obj.UpdatedAt,
			&obj.UpdatedBy,
			&obj.DeletedAt,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

	}

	res.Status = http.StatusOK
	res.Msg = "Category founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func WhereCategory(col string, val string) (*Category, error) {

	var obj Category

	con := A.GetDB()

	sql := "SELECT * FROM categories WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.Name,
			&obj.CreatedAt,
			&obj.CreatedBy,
			&obj.UpdatedAt,
			&obj.UpdatedBy,
			&obj.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

	}

	return &obj, nil
}

func CreateCategory(category *Category) (ResponseNoData, error) {

	CategoryLock.Lock()
	defer CategoryLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		INSERT INTO public.categories(
			name, created_at, created_by)
		VALUES ($1, NOW(), $2);
	`

	_, err := con.Exec(sql, category.Name, category.CreatedBy)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Category created successfully"
	res.Success = true

	return res, nil
}

func DeleteCategory(id *int64) (ResponseNoData, error) {
	CategoryLock.Lock()
	defer CategoryLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.categories SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Category trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateCategory(category *Category) (ResponseNoData, error) {
	CategoryLock.Lock()
	defer CategoryLock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.categories
			SET name = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3;
	`

	_, err := con.Exec(sql, category.Name, category.UpdatedBy, category.Id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Category updated successfully"
	res.Success = true

	return res, nil
}