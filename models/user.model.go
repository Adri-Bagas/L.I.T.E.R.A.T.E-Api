package models

import (
	"net/http"
	A "perpus_api/db"
	H "perpus_api/helpers"
	"sync"
)

type User struct {
	ID         int64   `db:"id" json:"id"`
	Username   string  `db:"username" json:"username"`
	Email      string  `db:"email" json:"email"`
	Password   string  `db:"password" json:"password"`
	LastActive *string `db:"last_active" json:"last_active"`
	CreatedAt  *string `db:"created_at" json:"created_at"`
	UpdatedAt  *string `db:"updated_at" json:"updated_at"`
	DeletedAt  *string `db:"deleted_at" json:"deleted_at"`
	CreatedBy  *int64 `db:"created_by" json:"created_by"`
	UpdatedBy  *int64 `db:"updated_by" json:"updated_by"`
	DeletedBy  *int64 `db:"deleted_by" json:"deleted_by"`
	Role       int     `db:"role" json:"role"`
	ProfilePic *string `db:"profile_pic" json:"profile_pic"`
}

type UserSafe struct {
	ID         int64   `json:"id"`
	Username   string  `json:"username"`
	Email      string  `json:"email"`
	LastActive *string `json:"last_active"`
	DeletedAt  *string `json:"deleted_at"`
	ProfilePic *string `json:"profile_pic"`
}

var lock = sync.Mutex{}

func GetAllUser(getThrashed bool) (ResponseMultiple, error) {
	var obj User
	var arrobj []UserSafe
	var res ResponseMultiple

	con := A.GetDB()

	var sql string

	if !getThrashed {
		sql = `
			SELECT t1.*, t2.location as profile_pic FROM public.users t1 left join public.medias t2 on t1.id = t2.model_id where t1.deleted_at is null;
		`
	} else {
		sql = `
			SELECT t1.*, t2.location as profile_pic FROM public.users t1 left join public.medias t2 on t1.id = t2.model_id where t1.deleted_at is not null;
		`
	}

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
			&obj.ID,
			&obj.Username,
			&obj.Email,
			&obj.Password,
			&obj.LastActive,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.Role,
			&obj.ProfilePic,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		arrobj = append(arrobj, UserSafe{
			ID:         obj.ID,
			Username:   obj.Username,
			Email:      obj.Email,
			LastActive: obj.LastActive,
			DeletedAt:  obj.DeletedAt,
			ProfilePic: obj.ProfilePic,
		})
	}

	res.Status = http.StatusOK
	res.Msg = "Users founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func FindUser(id int64) (Response, error) {
	var obj User
	// var objSafe UserSafe
	var res Response

	con := A.GetDB()

	sql := `
		SELECT t1.*, t2.location as profile_pic FROM public.users t1 left join public.medias t2 on t1.id = t2.model_id WHERE t1.id = $1;
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
			&obj.ID,
			&obj.Username,
			&obj.Email,
			&obj.Password,
			&obj.LastActive,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.Role,
			&obj.ProfilePic,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		// objSafe = UserSafe{
		// 	ID: obj.ID,
		// 	Username: obj.Username,
		// 	Email: obj.Email,
		// 	LastActive: obj.LastActive,
		// }

	}

	res.Status = http.StatusOK
	res.Msg = "Users founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func CreateUser(password string, username string, email string, role int64) (ResponseNoData, error) {

	lock.Lock()
	defer lock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	hashed, _ := H.HashPassword(password)

	sql := `
		INSERT INTO public.users(username, email, password, role)
		VALUES ($1, $2, $3, $4);
	`

	_, err := con.Exec(sql, username, email, hashed, role)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "User created successfully"
	res.Success = true

	return res, nil
}

func DeleteUser(id int64) (ResponseNoData, error) {
	lock.Lock()
	defer lock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.users SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "User trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateUser(id int64, password string, username string, email string, role int64) (ResponseNoData, error) {
	lock.Lock()
	defer lock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	hashed, _ := H.HashPassword(password)

	sql := `
		UPDATE public.users SET username = $1, email = $2, "password" = $3, role = $4, updated_at = NOW() WHERE id = $5;
	`

	_, err := con.Exec(sql, username, email, hashed, role, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, nil
	}

	res.Status = http.StatusOK
	res.Msg = "User updated successfully"
	res.Success = true

	return res, nil
}

func WhereUser(col string, val string) (*User, error) {

	var obj User

	con := A.GetDB()

	sql := "SELECT * FROM users WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.ID,
			&obj.Username,
			&obj.Email,
			&obj.Password,
			&obj.LastActive,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.CreatedBy,
			&obj.UpdatedBy,
			&obj.DeletedAt,
			&obj.DeletedBy,
			&obj.Role,
		)

		if err != nil {
			return nil, err
		}

	}

	return &obj, nil
}

func SetUserLastActive(id int64) error {
	lock.Lock()
	defer lock.Unlock()

	con := A.GetDB()

	sql := `
		UPDATE public.users SET last_active = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	return err
}

func SetUserProfilePic(path string, id int64) error {
	lock.Lock()
	defer lock.Unlock()

	con := A.GetDB()

	sql := `INSERT INTO public.medias (model_name, model_id, media_type, location) VALUES ('user', $1, 'image', $2);`

	_, err := con.Exec(sql, id, path)

	return err
}