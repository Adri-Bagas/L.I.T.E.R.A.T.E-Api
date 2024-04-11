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


func GetAllUser(c echo.Context) error {
	res, err := M.GetAllUser(false)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func GetAllThrashedUser(c echo.Context) error {
	res, err := M.GetAllUser(true)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func CreateUser(c echo.Context) error {

	role, err := strconv.ParseInt(c.FormValue("role"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	var u M.UserForm

	if role == 0 {
		var (
			full_name string
			phone_number string
			address string
		)
		full_name = c.FormValue("full_name")
		phone_number = c.FormValue("phone_number")
		address = c.FormValue("address")
		u = M.UserForm{
			Username:    c.FormValue("username"),
			Email:       c.FormValue("email"),
			Password:    c.FormValue("password"),
			Role:        role,
			FullName:    &full_name,
			PhoneNumber: &phone_number,
			Address:     &address,
		}
	} else {
		u = M.UserForm{
			Username: c.FormValue("username"),
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			Role:     role,
		}
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.CreateUser(u)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateUser(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	role, err := strconv.ParseInt(c.FormValue("role"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	u := &M.UserForm{
		Username: c.FormValue("username"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		Role:     role,
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.UpdateUser(userID, u.Password, u.Username, u.Email, u.Role)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletedUser(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.DeleteUser(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func FindUser(c echo.Context) error {

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	res, err := M.FindUser(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, res)
	}

	return c.JSON(http.StatusOK, res)
}

func SetUserProfilePic(c echo.Context) error {
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

	find, _ := M.FindMedia("user", claims.ID)

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

	err = M.SetUserProfilePic(dst.Name(), int64(claims.ID))

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
