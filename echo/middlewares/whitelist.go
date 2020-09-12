package middlewares

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////

func WhiteList(getList func() []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			list := getList()
			if len(list) == 0 {
				return next(c)
			}
			log.Debugf(nil, "remote ip:%s", c.Request().Header.Get("X-Real-IP"))
			for _, ipPattern := range list {
				ipPattern = strings.Replace(ipPattern, "*", `\d{1,3}`, -1)
				ipPattern = strings.Replace(ipPattern, ".", `\.`, -1)
				ipPattern = fmt.Sprintf(`^%s$`, ipPattern)
				log.Debugf(nil, "ipPattern:%s", ipPattern)
				matched, err := regexp.MatchString(ipPattern, c.Request().Header.Get("X-Real-IP"))
				if err == nil && matched {
					return next(c)
				}
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}
