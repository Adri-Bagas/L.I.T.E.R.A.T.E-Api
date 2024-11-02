package controllers

import (
	"io"
	"net/http"
	"os"
	M "perpus_api/models"
	A "perpus_api/db"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetAllMember(c echo.Context) error {
	res, err := M.GetAllMember(false)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func GetAllThrashedMember(c echo.Context) error {
	res, err := M.GetAllMember(true)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindMemberWithTransaction(c echo.Context) error {
	memberId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var res M.Response

	con := A.GetDB()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	member, err := M.FindMemberObj(memberId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusInternalServerError, res)
	}

	sql := `
	select 
		t1.id as trans_id, 
		t1.date, 
		t1.detail, 
		t1.created_at, 
		t1.expected_return_date,
		t1.approval_status,
		t1.is_returned,
		t2.* 
	from public.transactions t1 
		inner join public.users t2 on t1.approver_id = t2.id 
	where transaction_type = 'LOAN' and member_id = $1;
	`

	var obj M.TransactionLoan
	var arrobj []M.TransactionLoan

	rows, err := con.Query(sql, memberId)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}

	defer rows.Close()

	for rows.Next() {

		var user M.User

		err := rows.Scan(
			&obj.Id,
			&obj.Date,
			&obj.Detail,
			&obj.CreatedAt,
			&obj.ExpectedReturnDate,
			&obj.ApprovalStatus,
			&obj.IsReturned,
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.LastActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.DeletedBy,
			&user.Role,
		)

		if err != nil {
			res.Status = http.StatusInternalServerError
			res.Msg = err.Error()
			res.Success = false
			return c.JSON(http.StatusInternalServerError, res)
		}

		*user.Password = ""

		obj.Approver = &user

		arrobj = append(arrobj, obj)
	}

	member.Transaction = &arrobj

	res.Status = http.StatusOK
	res.Msg = "Founded Member and transactions!"
	res.Success = true
	res.Data = member

	return c.JSON(http.StatusOK, res)
}

func GetAllMemberIdName(c echo.Context) error {
	category, err := M.GetAllMemberObj()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"err": err.Error(),
		})
	}

	datas := make(map[int]string)
	
	for _, val := range category {
		datas[int(val.ID)] = val.FullName
	}

	return c.JSON(http.StatusOK, datas)
}

func CreateMember(c echo.Context) error {

	var u M.MemberForm

		var (
			full_name string
			phone_number string
			address string
		)
		full_name = c.FormValue("full_name")
		phone_number = c.FormValue("phone_number")
		address = c.FormValue("address")
		u = M.MemberForm{
			Username:    c.FormValue("username"),
			Email:       c.FormValue("email"),
			Password:    c.FormValue("password"),
			FullName:    &full_name,
			PhoneNumber: &phone_number,
			Address:     &address,
		}
	

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.CreateMember(u)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateMember(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	var u M.MemberForm

		var (
			full_name string
			phone_number string
			address string
		)
		full_name = c.FormValue("full_name")
		phone_number = c.FormValue("phone_number")
		address = c.FormValue("address")
		u = M.MemberForm{
			Username:    c.FormValue("username"),
			Email:       c.FormValue("email"),
			Password:    c.FormValue("password"),
			FullName:    &full_name,
			PhoneNumber: &phone_number,
			Address:     &address,
		}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.UpdateMember(userID, u)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedMember(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeleteMember(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindMember(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindMember(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func SetMemberProfilePic(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return c.JSON(http.StatusBadRequest, "JWT token missing or invalid")
	}

	claims := user.Claims.(*M.JwtCustomClaims)

	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": err.Error(),
			},
		)
	}

	if file == nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": "File Null!",
			},
		)
	}

	find, _ := M.FindMedia("member", claims.ID)

	if find != nil {
		err = os.Remove(find.Location)

		if err != nil {
			return c.JSON(
				500,
				echo.Map{
					"msg": err.Error(),
				},
			)
		}

		err = M.DeleteMedia(int64(find.Id))

		if err != nil {
			return c.JSON(
				500,
				echo.Map{
					"msg": err.Error(),
				},
			)
		}
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": err.Error(),
			},
		)
	}
	defer src.Close()

	dst, err := os.Create("uploads/profiles/" + file.Filename)
	if err != nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": err.Error(),
			},
		)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": err.Error(),
			},
		)
	}

	err = M.SetMemberProfilePic(dst.Name(), int64(claims.ID))

	if err != nil {
		return c.JSON(
			500,
			echo.Map{
				"msg": err.Error(),
			},
		)
	}

	return c.JSON(200, echo.Map{
		"path":    dst.Name(),
		"success": true,
	})

}