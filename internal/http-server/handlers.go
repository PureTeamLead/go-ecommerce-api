package http_server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const homePage = `
<html>
<body>
<h1>Hello on EDRIP-shop page</h1>
</body>
</html>`

func HomePage(c echo.Context) error {
	return c.HTML(http.StatusOK, homePage)
}
