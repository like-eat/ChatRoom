package ssl

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"kama_chat_server/pkg/zlog"
	"strconv"
)

func TlsHandler(host string, port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     host + ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			zlog.Fatal(err.Error())
			return
		}

		c.Next()
	}
}
