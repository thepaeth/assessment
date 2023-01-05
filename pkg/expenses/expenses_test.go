package expenses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var bodyExpense = `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`

func TestCreateExpenseMock(t *testing.T) {
	// setup echo server
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(bodyExpense))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	// Assertions
	if assert.NoError(t, CreateExpense(c)) {
		exp := &Expenses{}
		err := json.Unmarshal(res.Body.Bytes(), exp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, 0, exp.ID)
		assert.Equal(t, "strawberry smoothie", exp.Title)
		assert.Equal(t, 79.00, exp.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
		assert.Equal(t, []string{"food", "beverage"}, exp.Tags)
	}
}
