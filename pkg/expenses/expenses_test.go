package expenses

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	bodyExpense = `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	mockData = Expenses{
		Title:  "strawberry smoothie",
		Amount: 79.00,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}
)

func TestCreateExpenseMock(t *testing.T) {
	// db mock
	mockdb, mock, _ := sqlmock.New()
	expMockSql := "INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id"
	expMockRow := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(expMockSql)).WithArgs(mockData.Title, mockData.Amount, mockData.Note, pq.Array(&mockData.Tags)).
		WillReturnRows(expMockRow)

	// setup echo server
	body, err := json.Marshal(mockData)
	if err != nil {
		t.Error(err)
		return
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(string(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	db = mockdb

	// Assertions
	if assert.NoError(t, CreateExpense(c)) {
		exp := &Expenses{}
		err := json.Unmarshal(res.Body.Bytes(), exp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, 1, exp.ID)
		assert.Equal(t, "strawberry smoothie", exp.Title)
		assert.Equal(t, 79.00, exp.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
		assert.Equal(t, []string{"food", "beverage"}, exp.Tags)
	}
}
