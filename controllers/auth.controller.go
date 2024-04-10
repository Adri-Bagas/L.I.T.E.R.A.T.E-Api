package controllers

import (
	"net/http"
	"perpus_api/config"
	H "perpus_api/helpers"
	M "perpus_api/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {

	conf := config.GetConfig()

	var res M.Response

	email := c.FormValue("email")
	password := c.FormValue("password")

	userData, err := M.WhereUser("email", email)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false
		return c.JSON(http.StatusInternalServerError, res)
	}

	if userData.ID == 0 {
		res.Status = http.StatusInternalServerError
		res.Msg = "User not found!"
		res.Success = false
		return c.JSON(http.StatusInternalServerError, res)
	}

	check := H.CheckPasswordHash(password, userData.Password)

	if !check {
		res.Status = http.StatusBadRequest
		res.Msg = "Password do not match!"
		res.Success = false

		return c.JSON(http.StatusBadRequest, res)
	}

	err = M.SetUserLastActive(userData.ID)

	if(err != nil){
		res.Status = http.StatusBadRequest
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusBadRequest, res)
	}

	if(userData.Role == 0){
		res.Status = http.StatusUnauthorized
		res.Msg = "Unauthorized!"
		res.Success = false

		return c.JSON(http.StatusUnauthorized, res)
	}

	claims := &M.JwtCustomClaims{
		ID: int(userData.ID),
		Name:  userData.Username,
		Email: userData.Email,
		Role: userData.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(conf.SECRET_TOKEN_A))

	if err != nil {

		res.Status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Success = false

		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Status = http.StatusOK
	res.Msg = "Login Success!"
	res.Success = true
	res.Data = echo.Map{
		"token": t,
		"user":  M.UserSafe{
			ID: userData.ID,
			Username: userData.Username,
			Email: userData.Email,
			LastActive: userData.LastActive,
		},
	}

	return c.JSON(http.StatusOK, res)
}

func GetMe(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return c.JSON(http.StatusBadRequest, "JWT token missing or invalid")
	}

	claims := user.Claims.(*M.JwtCustomClaims)

	return c.JSON(http.StatusOK, echo.Map{
		"name": claims.Name,
		"email": claims.Email,
		"id": claims.ID,
	} )
}
