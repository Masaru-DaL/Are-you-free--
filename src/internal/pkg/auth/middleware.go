package auth

import (
	"src/internal/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Cookie Store
var CookieStore = sessions.NewCookieStore([]byte(config.Config.AUTH.SessionKey))

/* セッションを設定するミドルウェア */
func SessionMiddleware(store *sessions.CookieStore) echo.MiddlewareFunc {
	// セッション設定を作成する
	sessCfg := session.Config{
		Store: CookieStore,
	}

	// ミドルウェア関数を返す
	return session.MiddlewareWithConfig(sessCfg)
}
