package main

import (
	"Cart/controllers"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	//Load ENV file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Load echo framework and start using route.
	e := echo.New()

	controllers.Route(e)

	//log.Print("Server is running at -" + os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	controllers.CleanDB(0)
	expire, _ := strconv.ParseInt(os.Getenv("EXPIRE"), 10, 64)
	go func() {
		for range time.Tick(time.Millisecond * time.Duration(expire)) {
			controllers.CleanDB(expire)
		}
	}()

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
