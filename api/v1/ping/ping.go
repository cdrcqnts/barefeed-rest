package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping returns "Pong" as a response to a successful request.
// Endpoint: GET "feeds/:sid"
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Pong"})
	}
}
