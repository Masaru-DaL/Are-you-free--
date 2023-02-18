package freetime

import (
	"context"
	"net/http"
	"src/internal/config"
	"src/internal/entity"
	"src/internal/pkg/models/freetimes"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CreateFreeTime(ctx context.Context, db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(config.Config.Session.Name, c)
		userID := sess.Values[config.Config.Session.KeyName].(int)

		yearString, monthString, dayString := c.FormValue("year"), c.FormValue("month"), c.FormValue("day")

		year, _ := strconv.Atoi(yearString)
		month, _ := strconv.Atoi(monthString)
		day, _ := strconv.Atoi(dayString)

		dateFreeTime := &entity.DateFreeTime{
			UserID: userID,
			Year:   year,
			Month:  month,
			Day:    day,
		}

		dateFreeTime, err := freetimes.CreateDateFreeTime(ctx, db, dateFreeTime)
		if err != nil {
			return c.Render(http.StatusOK, "create-free-time", echo.Map{
				"error_message": entity.ERR_INTERNAL_SERVER_ERROR,
			})
		}

		return c.JSON(http.StatusOK, dateFreeTime)
	}
}
