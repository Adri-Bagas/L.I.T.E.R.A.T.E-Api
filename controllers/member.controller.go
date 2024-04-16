package controllers

import (
	"io"
	"net/http"
	"os"
	M "perpus_api/models"
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