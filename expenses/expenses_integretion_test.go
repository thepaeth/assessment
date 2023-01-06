//go:build integration

package expenses

import (
	"bytes"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseIt(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	var exp Expenses
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	res := request(t, http.MethodPost, uri("expenses"), body)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "strawberry smoothie", exp.Title)
	assert.Equal(t, 79.00, exp.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
	assert.Equal(t, []string{"food", "beverage"}, exp.Tags)
}

func TestGETExpenseByIDIt(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	c := seedExpense(t)
	var exp Expenses
	res := request(t, http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGETAllExpensesIt(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	seedExpense(t)
	var exp []Expenses
	res := request(t, http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(exp), 0)
}

func TestUpdateExpenseIt(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	c := seedExpense(t)
	var exp Expenses
	res := request(t, http.MethodPut, uri("expenses", strconv.Itoa(c.ID)), body)
	err := res.Decode(&exp)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
