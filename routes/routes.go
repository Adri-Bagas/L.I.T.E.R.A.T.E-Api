package routes

import (
	"net/http"
	"perpus_api/config"
	CR "perpus_api/controllers"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
	M "perpus_api/models"
)

type (
	Greetings struct {
		Msg        string `json:"msg"`
		StatusCode int    `json:"status_code"`
		Version    int    `json:"api_version"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func Init() *echo.Echo {

	conf := config.GetConfig()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {

		datas := &Greetings{
			Msg:        "Hello and Welcome, This is the main endpoint for this app. please use us kindly!",
			StatusCode: 200,
			Version:    1,
		}

		return c.JSON(http.StatusOK, datas)
	})

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(M.JwtCustomClaims)
		},
		SigningKey: []byte(conf.SECRET_TOKEN_A),
	}

	//static files
	e.Static("/uploads", "uploads")

	// User Group Start
	userRoute := e.Group("/user")
	userRoute.Use(echojwt.WithConfig(config))
	userRoute.GET("", CR.GetAllUser)
	userRoute.POST("", CR.CreateUser)
	userRoute.PUT("/:id", CR.UpdateUser)
	userRoute.DELETE("/:id", CR.DeletedUser)
	userRoute.GET("/:id", CR.FindUser)
	userRoute.GET("/thrashed", CR.GetAllThrashedUser)
	userRoute.POST("/upload/prof", CR.SetUserProfilePic)
	// User Group End

	//Book Group Start
	bookRoute := e.Group("/book")
	bookRoute.Use(echojwt.WithConfig(config))
	bookRoute.GET("", CR.GetAllBook)
	bookRoute.POST("", CR.CreateBook)
	bookRoute.PUT("/:id", CR.UpdateBook)
	bookRoute.DELETE("/:id", CR.DeletedBook)
	bookRoute.GET("/:id", CR.FindBook)
	//Book Group End

	//Author Group Start
	authorRoute := e.Group("/author")
	authorRoute.Use(echojwt.WithConfig(config))
	authorRoute.GET("", CR.GetAllAuthor)
	authorRoute.POST("", CR.CreateAuthor)
	authorRoute.PUT("/:id", CR.UpdateAuthor)
	authorRoute.DELETE("/:id", CR.DeletedAuthor)
	authorRoute.GET("/:id", CR.FindAuthor)
	//Author Group End

	//Publisher Group Start
	publisherRoute := e.Group("/publisher")
	publisherRoute.Use(echojwt.WithConfig(config))
	publisherRoute.GET("", CR.GetAllPublisher)
	publisherRoute.POST("", CR.CreatePublisher)
	publisherRoute.PUT("/:id", CR.UpdatePublisher)
	publisherRoute.DELETE("/:id", CR.DeletedPublisher)
	publisherRoute.GET("/:id", CR.FindPublisher)
	//Publisher Group End

	//Login
	e.POST("/auth/login", CR.Login)
	//Auth
	authRouteSafe := e.Group("/auth")
	authRouteSafe.Use(echojwt.WithConfig(config))
	authRouteSafe.POST("/me", CR.GetMe)

	return e
}
