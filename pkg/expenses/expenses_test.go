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
	fakeData = []Expenses{
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
	mockData = Expenses{
		Title:  "strawberry smoothie",
		Amount: 79.00,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}
)

func TestCreateExpenseMock(t *testing.T) {
	db, mock, _ := sqlmock.New()
	expMockSql := "INSERT INTO expenses (title, amount, note, tags) values($1, $2, $3, $4) RETURNING id"
	expMockRow := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(expMockSql)).WithArgs(mockData.Title, mockData.Amount, mockData.Note, pq.Array(&mockData.Tags)).
		WillReturnRows(expMockRow)

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
	h := &handler{db}

	if assert.NoError(t, h.CreateExpense(c)) {
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
	db, mock, _ := sqlmock.New()
	expID := "1"
	expMockSql := "SELECT id, title, amount, note, tags FROM expenses WHERE id = $1"
	expMockRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expID, mockData.Title, mockData.Amount, mockData.Note, pq.Array(&mockData.Tags))

	mock.ExpectPrepare(regexp.QuoteMeta(expMockSql)).ExpectQuery().WithArgs(expID).WillReturnRows((expMockRow))

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
	db, mock, _ := sqlmock.New()
	expID := "1"
	expMockSql := "UPDATE expenses SET title = $2, amount = $3, note = $4, tags = $5 WHERE id = $1"
	expMockRow := sqlmock.NewResult(1, 1)

	mock.ExpectPrepare(regexp.QuoteMeta(expMockSql)).ExpectExec().
		WithArgs(expID, "apple smoothie", 89.00, "no discount", pq.Array(&[]string{"beverage"})).
		WillReturnResult(expMockRow)

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
	h := &handler{db}

	if assert.NoError(t, h.UpdateExpense(c)) {
		assert.Equal(t, http.StatusAccepted, res.Code)
	}
}

func TestGetAllExpenses(t *testing.T) {

	db, mock, _ := sqlmock.New()
	expMockSql := "SELECT id, title, amount, note, tags FROM expenses"
	expMockRow := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(fakeData[0].ID, fakeData[0].Title, fakeData[0].Amount, fakeData[0].Note, pq.Array(&fakeData[0].Tags)).
		AddRow(fakeData[1].ID, fakeData[1].Title, fakeData[1].Amount, fakeData[1].Note, pq.Array(&fakeData[1].Tags))

	mock.ExpectPrepare(regexp.QuoteMeta(expMockSql)).ExpectQuery().WillReturnRows((expMockRow))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expense", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	h := &handler{db}

	if assert.NoError(t, h.GetAllExpenses(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}
