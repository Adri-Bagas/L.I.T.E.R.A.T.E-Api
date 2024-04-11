package models

import (
	"net/http"
	A "perpus_api/db"
	H "perpus_api/helpers"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Member struct {
	ID          int64   `db:"id" json:"id"`
	Username    string  `db:"username" json:"username"`
	FullName    string  `db:"full_name" json:"full_name"`
	Email       string  `db:"email" json:"email"`
	Password    string  `db:"password" json:"password"`
	PhoneNumber string  `db:"phone_number" json:"phone_number"`
	Address     string  `db:"address" json:"address"`
	CreatedAt   *string `db:"created_at" json:"created_at"`
	UpdatedAt   *string `db:"updated_at" json:"updated_at"`
	DeletedAt   *string `db:"deleted_at" json:"deleted_at"`
	LastActive  *string `db:"last_active" json:"last_active"`
	ProfilePic  *string `db:"profile_pic" json:"profile_pic"`
}

type MemberSafe struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	FullName    string  `db:"full_name" json:"full_name"`
	Email       string  `json:"email"`
	LastActive  *string `json:"last_active"`
	DeletedAt   *string `json:"deleted_at"`
	PhoneNumber string  `json:"phone_number"`
	Address     string  `json:"address"`
	ProfilePic  *string `json:"profile_pic"`
}

type MemberForm struct {
	Username    string  `json:"username" validate:"required"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required"`
	Role        int64   `json:"role"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
}

var Memberlock = sync.Mutex{}

func GetAllMember(getThrashed bool) (ResponseMultiple, error) {
	var obj Member
	var arrobj []MemberSafe
	var res ResponseMultiple

	con := A.GetDB()

	var sql string

	if !getThrashed {
		sql = `
			SELECT t1.*, t2.location as profile_pic FROM public.members t1 left join public.medias t2 on t1.id = t2.model_id where t1.deleted_at and model_name = "member" is null;
		`
	} else {
		sql = `
			SELECT t1.*, t2.location as profile_pic FROM public.members t1 left join public.medias t2 on t1.id = t2.model_id where t1.deleted_at and model_name = "member" is not null;
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
			&obj.FullName,
			&obj.Email,
			&obj.Password,
			&obj.PhoneNumber,
			&obj.Address,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.LastActive,
			&obj.ProfilePic,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return res, err
		}

		arrobj = append(arrobj, MemberSafe{
			ID:          obj.ID,
			Username:    obj.Username,
			Email:       obj.Email,
			LastActive:  obj.LastActive,
			DeletedAt:   obj.DeletedAt,
			ProfilePic:  obj.ProfilePic,
			PhoneNumber: obj.PhoneNumber,
			FullName:    obj.FullName,
			Address:     obj.Address,
		})
	}

	res.Status = http.StatusOK
	res.Msg = "Users founded!"
	res.Success = true
	res.Datas = arrobj

	return res, nil
}

func FindMember(id int64) (Response, error) {
	var obj Member
	// var objSafe MemberSafe
	var res Response

	con := A.GetDB()

	sql := `
		SELECT t1.*, t2.location as profile_pic FROM public.members t1 left join public.medias t2 on t1.id = t2.model_id WHERE t1.id = $1;
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
			&obj.FullName,
			&obj.Email,
			&obj.Password,
			&obj.PhoneNumber,
			&obj.Address,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.LastActive,
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
	res.Msg = "Members founded!"
	res.Success = true
	res.Data = obj

	return res, nil
}

func CreateMember(d MemberForm) (ResponseNoData, error) {

	Memberlock.Lock()
	defer Memberlock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	hashed, _ := H.HashPassword(d.Password)

	sql := `
		INSERT INTO public.members(
			username, full_name, email, password, phone_number, address)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

	_, err := con.Exec(sql, d.Username, d.FullName, d.Email, hashed, d.PhoneNumber, d.Address)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Member created successfully"
	res.Success = true

	return res, nil
}

func DeleteMember(id int64) (ResponseNoData, error) {
	Memberlock.Lock()
	defer Memberlock.Unlock()

	var res ResponseNoData

	con := A.GetDB()

	sql := `
		UPDATE public.members SET deleted_at = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Member trashed successfully"
	res.Success = true

	return res, nil
}

func UpdateMember(id int64, password string, username string, email string, role int64) (ResponseNoData, error) {
	Memberlock.Lock()
	defer Memberlock.Unlock()

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
		return res, err
	}

	res.Status = http.StatusOK
	res.Msg = "Member updated successfully"
	res.Success = true

	return res, nil
}

func WhereMember(col string, val string) (*Member, error) {

	var obj Member

	con := A.GetDB()

	sql := "SELECT * FROM members WHERE " + col + " = $1;"

	rows, err := con.Query(sql, val)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.ID,
			&obj.Username,
			&obj.FullName,
			&obj.Email,
			&obj.Password,
			&obj.PhoneNumber,
			&obj.Address,
			&obj.CreatedAt,
			&obj.UpdatedAt,
			&obj.DeletedAt,
			&obj.LastActive,
			&obj.ProfilePic,
		)

		if err != nil {
			return nil, err
		}

	}

	return &obj, nil
}

func SetMemberLastActive(id int64) error {
	Memberlock.Lock()
	defer Memberlock.Unlock()

	con := A.GetDB()

	sql := `
		UPDATE public.members SET last_active = NOW() WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	return err
}

func SetMemberProfilePic(path string, id int64) error {
	Memberlock.Lock()
	defer Memberlock.Unlock()

	con := A.GetDB()

	sql := `INSERT INTO public.medias (model_name, model_id, media_type, location) VALUES ('member', $1, 'image', $2);`

	_, err := con.Exec(sql, id, path)

	return err
}

func GetMemberDataByJWT(c echo.Context) (*JwtCustomClaims, *string) {

	var errM string

	user, ok := c.Get("user").(*jwt.Token)

	if !ok {
		errM = "JWT token missing or invalid"
		return nil, &errM
	}

	claims := user.Claims.(*JwtCustomClaims)

	return claims, nil
}
