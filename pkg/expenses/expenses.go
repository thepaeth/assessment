package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

func CreateExpense(c echo.Context) error {
	// var exp Expenses
	// err := c.Bind(&exp)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, Err{Message: "Bad Request!!!"})
	// }
	// row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	// err = row.Scan(&exp.ID)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, Err{Message: "Can't Create Expense!!!!"})
	// }
	// return c.JSON(http.StatusCreated, exp)
	exp := new(Expenses)
	if err := c.Bind(exp); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "errror"})
	}
	return c.JSON(http.StatusCreated, exp)
}
