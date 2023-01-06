package expenses

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func checkAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == os.Getenv("AUTH_TOKEN") {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, Err{Message: "Authorized Failed"})
	}
}

func ExpRouter(e *echo.Echo) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	h := expService(db)
	h.InitDB()
	exp := e.Group("/expenses")
	exp.Use(checkAuthorized)
	exp.POST("", h.CreateExpense)
	exp.GET("/:id", h.GetExpense)
	exp.GET("", h.GetAllExpenses)
	exp.PUT("/:id", h.UpdateExpense)
}
