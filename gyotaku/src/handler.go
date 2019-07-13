package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IndexHandler(db *DBConn) echo.HandlerFunc {
	return func(c echo.Context) error {
		banner := `
                            ,､           __
                         ／/,＼     ／//＼
                     __,..-‐‐'..'..'─‐‐─'--'.'--.､｀= ､__                     , - ｱ
           _,..-‐::´;;;;;;;;::::::::::        :::;:::;;;;;:｀::‐.-..､_＿,.-‐ニ-/
    _,..-´__ ::::､:.丶  :::::::::::_＿ -‐ｧ      :::::;;;;;;;;;::::::::__ !三  〈
.∠_ = 〈●〉 ...〉....〕-‐ '´ﾐ_, -‐´             ::::::::::::::_,-‐'´￣  ｀＼ﾐ､ ヽ
.＼ ﾐヽ, ､_, .....Z../.......￣.´................_. _,....ｧ'´                 ｀‐- ゝ
    ｀ﾞ‐-  _＿  /」:::::::::::::::＿_＿ ,, -‐ﾍ､､､ヽ｀‐ '´
                ￣￣ ＼ヽ.l              ヽ､／
                          ＼|

                            Welcome to Gyotaku service!
		`
		return c.String(http.StatusOK, banner)
	}
}

func LoginHandler(db *DBConn) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(SessionName, c)
		if ok := sess.Values["username"]; ok != nil {
			return c.JSON(http.StatusOK, "alreday logged in")
		}

		username := c.FormValue("username")
		password := c.FormValue("password")

		if len(username) < 8 || len(password) < 8 {
			return c.JSON(http.StatusBadRequest, "both username and password must have 8 characters at least")
		}

		ret, err := db.UserExists(username)
		if err != nil {
			return err
		}

		if !ret { // create new user
			err := db.CreateUser(username, password)
			if err != nil {
				return err
			}
			sess.Values["username"] = username
			sess.Save(c.Request(), c.Response())
			return c.JSON(http.StatusOK, "success")
		}

		// verify user credentials
		ret, err = db.VerifyUser(username, password)
		if err != nil {
			return err
		}

		if !ret {
			return c.JSON(http.StatusBadRequest, "invalid username or password")
		}

		// success
		sess.Values["username"] = username
		sess.Save(c.Request(), c.Response())
		return c.JSON(http.StatusOK, "success")
	}
}

func GyotakuListHandler(db *DBConn) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(SessionName, c)
		username := sess.Values["username"].(string)

		gyotaku, err := db.GetGyotakuList(username)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, gyotaku)
	}
}

type GyotakuData struct {
	URL      string `json:"url"`
	Data     string `json:"data"`
	UserName string `json:"username"`
}

func GyotakuHandler(db *DBConn) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(SessionName, c)
		username := sess.Values["username"].(string)

		url := c.FormValue("url")

		// generate gyotaku id
		gid := fmt.Sprintf("%x", sha256.Sum256([]byte(url)))

		_, err := os.Stat(path.Join(GyotakuDir, gid))
		if !os.IsNotExist(err) {
			return c.JSON(http.StatusConflict, "this gyotaku has already been taken")
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// save gyotaku
		gyotakudata := &GyotakuData{
			URL:      url,
			Data:     string(body),
			UserName: username,
		}

		buf := bytes.NewBuffer(nil)
		err = gob.NewEncoder(buf).Encode(gyotakudata)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path.Join(GyotakuDir, gid), buf.Bytes(), 0644)
		if err != nil {
			return err
		}

		err = db.AddGyotakuList(username, gid)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, gid)
	}
}

func GyotakuViewHandler(db *DBConn) echo.HandlerFunc {
	return func(c echo.Context) error {
		// sess, _ := session.Get(SessionName, c)
		// username := sess.Values["username"].(string)
		gid := c.Param("gid")

		_, err := os.Stat(path.Join(GyotakuDir, gid))
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, "no such gyotaku")
		}

		_, err = ioutil.ReadFile(path.Join(GyotakuDir, gid))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusNotImplemented, "sorry but I couldn't make it by the submission deadline :P")

		// var gyotakudata GyotakuData
		// buf := bytes.NewBuffer(data)
		// err = gob.NewDecoder(buf).Decode(&gyotakudata)
		// if err != nil {
		// 	return err
		// }

		// if username != gyotakudata.UserName {
		// 	return c.JSON(http.StatusForbidden, "this is not your gyotaku")
		// }
	}
}

func FlagHandler(c echo.Context) error {
	data, err := ioutil.ReadFile("flag")
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(data))
}
