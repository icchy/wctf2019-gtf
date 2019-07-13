package main

import (
	"net"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func LoginRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(SessionName, c)
		if ok := sess.Values["username"]; ok != nil {
			return next(c)
		}
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

func InternalRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := net.ParseIP(c.RealIP())
		localip := net.ParseIP("127.0.0.1")
		if !ip.Equal(localip) {
			return echo.NewHTTPError(http.StatusForbidden)
		}
		return next(c)
	}
}
