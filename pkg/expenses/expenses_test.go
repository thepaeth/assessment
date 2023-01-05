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
	newFakeData = []Expenses{
		{
			ID:     1,
			Title:  "strawberry smoothie",
			Amount: 79.00,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		},
		{
			ID:     2,
			Title:  "iPhone 14 Pro Max 1TB",
			Amount: 66900.00,
			Note:   "birthday gift from my love",
			Tags:   []string{"gadget"},
		},
	}
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
	// Mock
	db, mock, _ := sqlmock.New()
	expID := "1"
	expMockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
	expMockRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expID, mockData.Title, mockData.Amount, mockData.Note, pq.Array(&mockData.Tags))

	mock.ExpectPrepare(regexp.QuoteMeta(expMockSql)).ExpectQuery().WithArgs(expID).WillReturnRows((expMockRow))

	// Setup echo server
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(expID)
	h := &handler{db}

	if assert.NoError(t, h.GetExpense(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}

}

func TestUpdateExpense(t *testing.T) {
	// Setup echo server
	expID := "1"

	body := `{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(expID)
	h := &newHandler{newFakeData}
	// Assertions
	if assert.NoError(t, h.UpdateExpense(c)) {
		assert.Equal(t, http.StatusAccepted, res.Code)
	}
}
