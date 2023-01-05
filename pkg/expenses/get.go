package expenses

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpense(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query expense statment"})
	}
	row := stmt.QueryRow(id)
	exp := Expenses{}
	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, (*pq.StringArray)(&exp.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNoContent, Err{Message: "Expense Data Not Found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan expense"})
	}
}

func (h *handler) GetAllExpenses(c echo.Context) error {
	stmt, err := h.db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query all expenses statment"})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't query all users"})
	}

	exps := []Expenses{}
	for rows.Next() {
		var exp Expenses
		err = rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, (*pq.StringArray)(&exp.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan expenses"})
		}
		exps = append(exps, exp)
	}
	return c.JSON(http.StatusOK, exps)
}
