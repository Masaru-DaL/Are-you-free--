package main

import (
	"src/internal/db"
	"src/internal/route"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// sqlite3とのコネクション接続
	db := db.ConnectDatabase()
	defer db.Close()

	e := route.InitRouting(db)
	/*
		Logger: リクエスト単位のログを出力する
		Recover: 予期せぬpanicを起こしてもサーバを落とさない
		CORS: アクセスを許可するオリジン(デフォルト)とメソッドの設定
	*/
	// e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339_nano} ${host} ${method} ${uri} ${status} ${header:my-header}` + "\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(":8080"))
}
