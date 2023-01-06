//go:build integration

package expenses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
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

func setup(t *testing.T) func() {
	e := echo.New()
	ExpRouter(e)
	go func() {
		e.Start(os.Getenv("PORT"))
	}()

	teardown := func() {
		ctx, down := context.WithTimeout(context.Background(), 10*time.Second)
		defer down()
		err := e.Shutdown(ctx)
		assert.NoError(t, err)
	}

	return teardown
}

func seedExpense(t *testing.T) Expenses {
	var exp Expenses
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	err := request(t, http.MethodPost, uri("expenses"), body).Decode(&exp)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return exp
}

func uri(paths ...string) string {
	host := fmt.Sprint("http://localhost", os.Getenv("PORT"))
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(t *testing.T, method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
