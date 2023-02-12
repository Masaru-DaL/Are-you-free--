package route

import (
	"database/sql"
	"html/template"
	"io"
	"src/internal/config"
	"src/internal/handler/templates"

	"github.com/labstack/echo/v4"
)

func InitRouting(db *sql.DB) *echo.Echo {
	e := echo.New()

	/* html/template対応 */
	initTemplate(e)
	e.Static("/css", config.Config.FilePath.CSS)
	e.Static("/img", config.Config.FilePath.Image)
	e.Static("/js", config.Config.FilePath.JS)
	e.Static("/icon", config.Config.FilePath.Icon)

	/* No authentication required */
	// noAuthGroup := e.Group("/auth")
	// noAuthGroup.GET("/signup", templates.SignupPage)
	e.GET("/signup", templates.SignupPage)
	e.GET("/login", templates.LoginPage)
	e.GET("/index", templates.TopPage)
	e.GET("/free-time", templates.SchedulePage)
	e.GET("/free-times", templates.SchedulesPage)
	e.GET("/free-time/new", templates.CreateSchedulePage)
	e.GET("/free-time/update", templates.UpdateSchedulePage)
	e.GET("/share/with_someone", templates.SharingToSomeonePage)
	e.GET("/share/with_me", templates.SharingFromSomeonePage)
	e.GET("/account/", templates.AccountPage)
	e.GET("/account/password_reset", templates.PasswordResetPage)
	e.GET("/account/password_re_registration", templates.PasswordReRegistrationPage)

	return e
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func initTemplate(e *echo.Echo) {
	renderer := &Template{
		templates: template.Must(template.New("t").ParseGlob("internal/web/template/*.html")),
	}
	e.Renderer = renderer
	e.Pre(templates.MethodOverride)
}
