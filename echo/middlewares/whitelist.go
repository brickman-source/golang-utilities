package middlewares

import (
	"github.com/brickman-source/golang-utilities/log"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////

func WhiteList(getList func() []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			list := getList()
			if len(list) == 0 {
				return next(c)
			}
			remoteIP := getRemoteIP(c.Request())
			if remoteIP == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			log.Debugf( "remote ip:%s", remoteIP)
			ipAddr := net.ParseIP(remoteIP)

			for _, networkAddr := range list {
				_, ipNetworkAddr, err := net.ParseCIDR(networkAddr)
				if err != nil {
					return c.NoContent(http.StatusUnauthorized)
				}
				if ipNetworkAddr.Contains(ipAddr) {
					log.Debugf( "%s contain %s?\n", ipNetworkAddr, ipAddr)
					return next(c)
				} else {
					log.Debugf( "%s not contain %s?\n", ipNetworkAddr, ipAddr)
				}
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
func getRemoteIP(req *http.Request) string {
	if req.Header.Get("X-Forwarded-For") != "" {
		return req.Header.Get("X-Forwarded-For")
	}
	if req.Header.Get("X-Real-IP") != "" {
		return req.Header.Get("X-Real-IP")
	}
	if req.Header.Get("Proxy-Client-IP") != "" {
		return req.Header.Get("Proxy-Client-IP")
	}
	if req.Header.Get("WL-Proxy-Client-IP") != "" {
		return req.Header.Get("WL-Proxy-Client-IP")
	}
	return ""
}
