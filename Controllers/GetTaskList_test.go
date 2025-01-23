package Controllers

import (
	"golang-assesment/Database"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetTaskList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM "AmplTaskList"`).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status"}).
		AddRow(1, "Test Task", "Test Description", "pending"))

	Database.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	gin.Default()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/tasks?page=1&limit=10", nil)

	GetTaskList(c)

	// Validate the response
	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"success":true,"message":"Retrieve tasks successfully!!","data":[{"id":1,"title":"Test Task","description":"Test Description","status":"pending"}],"page":1,"limit":10,"total":1}`
	assert.JSONEq(t, expected, w.Body.String())

	// Ensure the mock expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}
