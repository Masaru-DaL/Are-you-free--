package auth

import (
	"src/internal/infra/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Cookie Store
var CookieStore = sessions.NewCookieStore([]byte(config.Config.Session.SecretKey))

/* セッションを設定するミドルウェア */
func SessionMiddleware(store *sessions.CookieStore) echo.MiddlewareFunc {
	// セッション設定を作成する
	sessCfg := session.Config{
		Store: CookieStore,
	}

	// ミドルウェア関数を返す
	return session.MiddlewareWithConfig(sessCfg)
}

/* 非認証ルートにアクセスした場合に処理を行うミドルウェア */
func UnAuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// セッション情報の取得
		session, err := session.Get(config.Config.Session.Name, c)
		if err != nil {
			return HandleUnAuthError(c)
		}

		// セッションIDが空ではなかった場合
		if session.ID != "" {
			return HandleUnAuthError(c)
		}

		return next(c)
	}
}

/* 認証ルートにアクセスした場合にアクセスした場合に処理を行うミドルウェア */
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// セッション情報の取得
		sess, err := session.Get(config.Config.Session.Name, c)
		if err != nil {
			return HandleAuthError(c)
		}

		// セッションの値が空だった場合
		if sess.ID == "" {
			return HandleAuthError(c)
		}

		return next(c)
	}
}
