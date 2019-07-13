package main

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/syndtr/goleveldb/leveldb"
)

// SessionName for cookie store
const SessionName = "session"
const GyotakuDir = "gyotaku"

func main() {
	// create gyotaku directory if not exists
	_, err := os.Stat(GyotakuDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(GyotakuDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// open db
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbconn := &DBConn{db}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		log.Fatal(err)
	}
	store := sessions.NewCookieStore(secret)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	e.Use(session.Middleware(store))

	e.GET("/", IndexHandler(dbconn), LoginRequiredMiddleware)
	e.GET("/gyotaku", GyotakuListHandler(dbconn), LoginRequiredMiddleware)
	e.GET("/gyotaku/:gid", GyotakuViewHandler(dbconn), LoginRequiredMiddleware)
	e.GET("/flag", FlagHandler, InternalRequiredMiddleware)

	e.POST("/login", LoginHandler(dbconn))
	e.POST("/gyotaku", GyotakuHandler(dbconn), LoginRequiredMiddleware)

	e.Logger.Fatal(e.Start(":80"))
}
