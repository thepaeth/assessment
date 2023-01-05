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
	fakeData = map[string]*Expenses{
		"1": &Expenses{
			1,
			"strawberry smoothie",
			79.00,
			"night market promotion discount 10 bath",
			[]string{"food", "beverage"},
		},
		"2": &Expenses{
			2,
			"iPhone 14 Pro Max 1TB",
			66900.00,
			"birthday gift from my love",
			[]string{"gadget"},
		},
	}
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

func TestGetExpenseByIDSuccess(t *testing.T) {
	// Setup echo server
	expID := "1"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(expID)
	h := &expHandler{fakeData}

	if assert.NoError(t, h.GetExpense(c)) {
		exp := &Expenses{}
		err := json.Unmarshal(res.Body.Bytes(), exp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, fakeData[expID].ID, exp.ID)
		assert.Equal(t, fakeData[expID].Title, exp.Title)
		assert.Equal(t, fakeData[expID].Amount, exp.Amount)
		assert.Equal(t, fakeData[expID].Note, exp.Note)
		assert.Equal(t, fakeData[expID].Tags, exp.Tags)
	}

}
