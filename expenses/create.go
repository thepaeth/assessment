package expenses

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) CreateExpense(c echo.Context) error {
	var exp Expenses
	err := c.Bind(&exp)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusBadRequest, Err{Message: "Bad Request!!!"})
	}
	row := h.db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't Create Expense!!!!"})
	}
	return c.JSON(http.StatusCreated, exp)
}
