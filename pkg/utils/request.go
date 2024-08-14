package utils

import (
	"net"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetNetworkRequest(c echo.Context) (ip string, port int, userAgent string) {
	ip = c.RealIP()

	remoteAddr := c.Request().RemoteAddr
	if strings.Contains(remoteAddr, ":") {
		_, portStr, _ := net.SplitHostPort(remoteAddr)
		port, _ = strconv.Atoi(portStr)
	}
	userAgent = c.Request().Header.Get("User-Agent")
	return ip, port, userAgent
}
