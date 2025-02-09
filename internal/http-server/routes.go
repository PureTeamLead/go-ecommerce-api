package http_server

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

func initRoutes(router *echo.Echo) {

	router.GET("/", HomePage)
}

func StartServer(host, port string) {
	r := echo.New()

	addr := fmt.Sprintf("%s:%s", host, port)

	initRoutes(r)
	r.Logger.Fatal(r.Start(addr))
}
