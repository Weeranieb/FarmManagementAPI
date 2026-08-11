package logging_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/logging"
)

func TestUseHTTP_AccessLogUsesPathWithoutQuery(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	app := fiber.New()
	logging.UseHTTP(app)
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/health?token=secret", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Header.Get(fiber.HeaderXRequestID))

	body, err := io.ReadAll(&buf)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	require.NotEmpty(t, lines)

	var last map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[len(lines)-1]), &last))
	assert.Equal(t, "http_request", last["msg"])
	assert.Equal(t, "/health", last["path"])
	assert.Equal(t, "/health", last["route"])
	assert.Equal(t, "2xx", last["status_class"])
	assert.NotContains(t, last["path"], "token")
	assert.NotEmpty(t, last["request_id"])
}
