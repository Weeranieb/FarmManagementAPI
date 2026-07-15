package handler_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/handler"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"
	"github.com/weeranieb/boonmafarm-backend/src/internal/version"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLive_IncludesVersion(t *testing.T) {
	// GIVEN
	app := fiber.New()
	app.Get("/health", handler.Live)

	// WHEN
	resp, err := app.Test(httptest.NewRequest("GET", "/health", nil))
	require.NoError(t, err)
	body := decodeJSON(t, resp.Body)

	// THEN
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok", body["status"])
	assert.Equal(t, version.Version, body["version"])
}

func TestReady_DBOK(t *testing.T) {
	// GIVEN
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	app := fiber.New()
	app.Get("/ready", handler.Ready(db))

	// WHEN
	resp, err := app.Test(httptest.NewRequest("GET", "/ready", nil))
	require.NoError(t, err)
	body := decodeJSON(t, resp.Body)

	// THEN
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok", body["status"])
	assert.Equal(t, version.Version, body["version"])
}

func TestReady_DBFail_Returns503(t *testing.T) {
	// GIVEN — closed underlying sql.DB so Ping fails
	metrics.Reset()
	db, err := gorm.Open(sqlite.Open("file:ready_fail?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	app := fiber.New()
	app.Get("/ready", handler.Ready(db))

	// WHEN
	resp, err := app.Test(httptest.NewRequest("GET", "/ready", nil))
	require.NoError(t, err)
	body := decodeJSON(t, resp.Body)

	// THEN
	require.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	assert.Equal(t, "unavailable", body["status"])
	assert.Equal(t, "database", body["error"])
	assert.Equal(t, version.Version, body["version"])
	assert.Contains(t, metrics.Exposition(), "db_ping_failures_total 1")
}

func decodeJSON(t *testing.T, r io.Reader) map[string]any {
	t.Helper()
	raw, err := io.ReadAll(r)
	require.NoError(t, err)
	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body))
	return body
}
