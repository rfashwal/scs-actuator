package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfashwal/scs-actuator/internal"
	"github.com/rfashwal/scs-utilities/net"
)

type IpResponse struct {
	Ip string `json:"ip"`
}

func NewRouter(s internal.Service) *gin.Engine {
	r := gin.Default()

	r.POST("/ping", pingHandler(s))

	return r
}

func pingHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		ip, err := net.GetIP(true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, &IpResponse{Ip: ip})
	}
}
