package api

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mr, err := miniredis.Run()
	assert.NoError(t, err)

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	API := New(sqlxDB, redisClient)
	assert.NotNil(t, API)
	assert.NotNil(t, API.e)
}

func TestHealthy(t *testing.T) {

	opt := sqlmock.MonitorPingsOption(true)
	mockDB, mock, err := sqlmock.New(opt)
	assert.NoError(t, err)
	defer mockDB.Close()


	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")


	rows := sqlmock.NewRows([]string{"bool"}).AddRow(true)
	mock.
		ExpectPing().
		WillReturnError(nil)
	mock.
		ExpectQuery(`SELECT true`).
		WillReturnRows(rows).
		WillReturnError(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthy", http.NoBody)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	healthy(sqlxDB)(c)

	var body = struct {
		Message string `json:"message"`
	}{}

	assert.Equal(t, http.StatusOK, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusOK), body.Message)

}