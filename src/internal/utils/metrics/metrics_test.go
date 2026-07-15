package metrics_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/metrics"
)

func TestMetrics_RouteTemplateNotRawPath(t *testing.T) {
	// GIVEN — metrics middleware and a parameterized route
	metrics.Reset()
	app := fiber.New()
	app.Use(metrics.Middleware())
	app.Get("/api/v1/pond/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// WHEN — two different raw paths hit the same template
	for _, id := range []string{"1", "99"} {
		req := httptest.NewRequest("GET", "/api/v1/pond/"+id, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	}
	req404 := httptest.NewRequest("GET", "/api/v1/pond/1/extra", nil)
	resp404, err := app.Test(req404)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp404.StatusCode)

	// THEN — one route label for the template; 404 uses unmatched (not raw path)
	body := metrics.Exposition()
	assert.Contains(t, body, `http_requests_total{method="GET",route="/api/v1/pond/:id",status_class="2xx"} 2`)
	assert.Contains(t, body, `http_requests_total{method="GET",route="unmatched",status_class="4xx"} 1`)
	assert.NotContains(t, body, `/api/v1/pond/1`)
	assert.NotContains(t, body, `/api/v1/pond/99`)
	assert.Contains(t, body, "http_request_duration_seconds_count")
}

func TestMetrics_HandlerExposesText(t *testing.T) {
	// GIVEN
	metrics.Reset()
	metrics.IncPanic()
	metrics.IncDBPingFailure()
	app := fiber.New()
	app.Get("/metrics", metrics.Handler())

	// WHEN
	resp, err := app.Test(httptest.NewRequest("GET", "/metrics", nil))
	require.NoError(t, err)
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	text := string(raw)

	// THEN
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.True(t, strings.HasPrefix(resp.Header.Get(fiber.HeaderContentType), "text/plain"))
	assert.Contains(t, text, "http_panics_total 1")
	assert.Contains(t, text, "db_ping_failures_total 1")
}
