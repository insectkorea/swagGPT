package example

import (
	"net/http"

	"github.com/gin-gonic/gin"
	echo "github.com/labstack/echo/v4"
)

// Helloworld handler function
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

// EchoHandler handler function
func EchoHandler(c echo.Context) error {
	return c.String(http.StatusOK, "hello echo")
}
