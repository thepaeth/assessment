package expenses

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func ExpRouter(e *echo.Echo) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	h := expService(db)
	h.InitDB()
	exp := e.Group("/expenses")
	{
		exp.POST("", h.CreateExpense)
		exp.GET("/:id", h.GetExpense)
		exp.GET("", h.GetAllExpenses)
		exp.PUT("/:id", h.UpdateExpense)
	}

}
