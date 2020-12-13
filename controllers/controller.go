package controllers

import (
	"Cart/models"
	"Cart/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db = utils.ConnectDB()

//Route - Route function for routing
func Route(e *echo.Echo) {

	e.GET("/:url", GetData)
	e.POST("/", PostData)
	e.File("/", "assets/index.html")
	e.Use(CORSMiddlewareWrapper)
	e.Use(RateLimitWithConfig(RateLimitConfig{
		Limit: 2,
		Burst: 2,
	}))

}

//GetData - Fetch data from db using url
func GetData(e echo.Context) error {
	response := &models.Request{}

	url := e.Param("url")
	db.Where("url = ?", url).Find(response)
	if response.Data == "" {
		return e.JSON(http.StatusNotFound, `{message : "404 Not Found"}`)
	}
	return e.String(http.StatusOK, response.Data)

}

//PostData - Receive Data and store in db
func PostData(e echo.Context) error {
	request := &models.Request{} //new(models.Request)

	//Read request body and decode to parts and store in db
	b, _ := ioutil.ReadAll(e.Request().Body)
	request.Data = string(b)
	request.URL = URLGen()

	//StoreRequest(request)
	createdRequest := db.Create(request)
	if createdRequest.Error != nil {
		return e.String(http.StatusInternalServerError, "403 Internal Server Error")
	}

	return e.String(http.StatusOK, ("http://" + os.Getenv("HOST") + "/" + request.URL + "\n"))
}

//CleanDB - clear expired data
func CleanDB(expire int64) {
	Now := time.Now().UnixNano() / int64(time.Millisecond)
	db.Unscoped().Where("created_at < ?", (Now - expire)).Delete(&models.Request{})
}

//URLGen - Generate unoverlapped ramdomnized url
func URLGen() string {
	request := &models.Request{}
	randString := utils.RandStringGen()
	if err := db.Where("url = ?", randString).First(request).Error; err != nil {
		return randString
	}
	return URLGen()
}

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins: []string{req.Header.Get("Origin")},
			AllowHeaders: []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-With", "Origin", "proxy", "proxy_pass"},
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}

type (
	RateLimitConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		Limit   int
		Burst   int
	}
)

var DefaultRateLimitConfig = RateLimitConfig{
	Skipper: middleware.DefaultSkipper,
	Limit:   2,
	Burst:   1,
}

func RateLimitMiddleware() echo.MiddlewareFunc {
	return RateLimitWithConfig(DefaultRateLimitConfig)
}

func RateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRateLimitConfig.Skipper
	}
	var limiter = rate.NewLimiter(rate.Limit(config.Limit), config.Burst)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			if limiter.Allow() == false {
				return echo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}
