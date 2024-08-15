package middleware

import (
	"fmt"
	"gateway/internal/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func PseudoHandler(host string) gin.HandlerFunc {
	return func (c *gin.Context) {
		path := c.Request.URL.Path

		url := host + path
		log.Debug().Msg(fmt.Sprintf("url value is: %s", url))

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error creating request: %s", err.Error()))
			utils.NewResponse(
				utils.ErrCreateRequest.Error(),
				http.StatusInternalServerError,
				nil,
			).Write(c.Writer)
			c.Abort()
			return
		}

		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		resp, err := http.DefaultClient.Do(req)
		fmt.Printf("%+v", resp)

		if err != nil {
			utils.NewResponse(
				utils.ErrClientSend.Error(),
				http.StatusInternalServerError,
				nil,
			).Write(c.Writer)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		for k, v := range resp.Header {
			c.Header(k, v[0])
		}

		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
