package route

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"src/internal/handlers/account"
	"src/internal/handlers/admin"
	"src/internal/handlers/freetimes"
	"src/internal/handlers/templates"
	"src/internal/infra/auth"
	"src/internal/infra/config"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func InitRouting(db *sqlx.DB) *echo.Echo {
	e := echo.New()

	/* html/template対応 */
	initTemplate(e)
	e.Static("/css", config.Config.FilePath.CSS)
	e.Static("/img", config.Config.FilePath.Image)
	e.Static("/js", config.Config.FilePath.JS)
	e.Static("/icon", config.Config.FilePath.Icon)

	ctx := context.Background()
	e.GET("/admin/users", admin.GetUsers(ctx, db))

	/* Unauthorized routing group. */
	e.Use(auth.SessionMiddleware(auth.CookieStore))

	unAuthenticatedGroup := e.Group("/auth")
	unAuthenticatedGroup.Use(auth.UnAuthenticatedMiddleware)
	unAuthenticatedGroup.GET("/signup", templates.SignupPage)
	unAuthenticatedGroup.POST("/signup", account.Signup(ctx, db))
	unAuthenticatedGroup.GET("/login", templates.LoginPage)
	unAuthenticatedGroup.POST("/login", account.Login(ctx, db))

	/* Authorized routing group. */
	authenticatedGroup := e.Group("/")
	authenticatedGroup.Use(auth.AuthenticatedMiddleware)
	e.GET("index", templates.TopPage(ctx, db))
	// e.GET("index", templates.TopPage)

	// e.GET("/index", templates.TopPage)
	e.GET("/free-time/:id", templates.FreeTimePage(ctx, db))
	e.GET("/free-times", templates.FreeTimesPage)
	e.GET("/free-time/create", templates.CreateFreeTimePage)
	e.POST("/free-time/create", freetimes.CreateFreeTime(ctx, db))
	e.GET("/free-time/update", templates.UpdateFreeTimePage)
	e.GET("/share/with_someone", templates.ShareWithSomeonePage)
	e.GET("/share/with_me", templates.ShareWithMePage)
	// e.GET("/account/", templates.AccountPage)
	authenticatedGroup.GET("account/", templates.AccountPage)
	authenticatedGroup.GET("account/edit", templates.AccountEditPage)
	authenticatedGroup.GET("logout", account.Logout(ctx, db))

	e.GET("/account/password_reset", templates.PasswordResetPage)
	e.GET("/account/password_re_registration", templates.PasswordReRegistrationPage)

	return e
}

type Template struct {
	templates *template.Template
	json      func(v interface{}) (string, error)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if strings.Contains(c.Request().Header.Get("Accept"), "application/json") {
		if t.json == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "JSON renderer not configured"})
		}
		jsonStr, err := t.json(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		_, err = io.WriteString(w, jsonStr)
		return err
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *Template) SetJSON(f func(v interface{}) (string, error)) {
	t.json = f
}

func initTemplate(e *echo.Echo) {
	// renderer := &Template{
	// 	templates: template.Must(template.New("t").ParseGlob("internal/web/template/*.html")),
	// }
	// renderer.SetJSON(func(v interface{}) (string, error) {
	// 	b, err := json.Marshal(v)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return string(b), nil
	// })
	renderer := &Template{
		templates: template.Must(template.New("t").Funcs(template.FuncMap{
			"marshal": func(v interface{}) (string, error) {
				b, err := json.Marshal(v)
				if err != nil {
					return "", err
				}
				return string(b), nil
			},
		}).ParseGlob("internal/web/template/*.html")),
	}
	e.Renderer = renderer
	e.Pre(templates.MethodOverride)
}
