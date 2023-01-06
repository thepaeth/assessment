package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpense(c echo.Context) error {
	id := c.Param("id")
	exp := Expenses{}

	if err := c.Bind(&exp); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	stmt, err := h.db.Prepare("UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare statment update"})
	}

	if _, err := stmt.Exec(id, exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error execute update"})
	}
	return c.JSON(http.StatusOK, exp)
}
