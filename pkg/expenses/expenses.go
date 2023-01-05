package expenses

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

type handler struct {
	db *sql.DB
}

type newHandler struct {
	db []Expenses
}

type expHandler struct {
	db map[string]*Expenses
}

func CreateExpense(c echo.Context) error {
	var exp Expenses
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Bad Request!!!"})
	}
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't Create Expense!!!!"})
	}
	return c.JSON(http.StatusCreated, exp)
}

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

func (h *newHandler) UpdateExpense(c echo.Context) error {
	id := c.Param("id")
	newExp := Expenses{}
	if err := c.Bind(&newExp); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	for _, k := range h.db {
		if fmt.Sprint(k.ID) == id {
			return c.JSON(http.StatusAccepted, &Expenses{
				ID:     k.ID,
				Title:  newExp.Title,
				Amount: newExp.Amount,
				Note:   newExp.Note,
				Tags:   newExp.Tags,
			})
		}
	}
	return c.JSON(http.StatusAccepted, newExp)
}
