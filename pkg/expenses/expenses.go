package expenses

import (
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

func (h *expHandler) GetExpense(c echo.Context) error {
	id := c.Param("id")
	exp := h.db[id]
	if exp == nil {
		return c.JSON(http.StatusNotFound, Err{Message: "Data Not Found"})
	}
	return c.JSON(http.StatusOK, exp)
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
