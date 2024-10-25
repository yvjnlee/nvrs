package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nvrs-gateway/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	var err error
	storage.DB, err = sql.Open("sqlite3", ":memory:") // In-memory SQLite
	if err != nil {
		panic("failed to connect to the in-memory database")
	}

	// Create required tables
	storage.DB.Exec(`CREATE TABLE agents (id INTEGER PRIMARY KEY, name TEXT, role TEXT, status TEXT)`)
	storage.DB.Exec(`CREATE TABLE tasks (id INTEGER PRIMARY KEY, agent_id INTEGER, task TEXT)`)
}

func teardownTestDB() {
	storage.DB.Close()
}

func TestMain(m *testing.M) {
	setupTestDB()    // Initialize in-memory database
	code := m.Run()  // Run tests
	teardownTestDB() // Clean up
	os.Exit(code)    // Exit with the test's return code
}

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "healthy")
	}
}

func TestRegisterAgent(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/agents/register", strings.NewReader(`{"name":"Agent007","role":"Spy"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assuming storage.DB is set up with test data
	if assert.NoError(t, RegisterAgent(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Agent registered successfully")
	}
}

func TestSubmitTask(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"agent_id":1,"task":"Analyze data"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, SubmitTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Task submitted successfully")
	}
}

func TestUpdateStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/agents/status", strings.NewReader(`{"agent_id":1,"status":"working"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, UpdateStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Status updated successfully")
	}
}
