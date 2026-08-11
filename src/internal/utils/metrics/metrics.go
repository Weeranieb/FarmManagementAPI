// Package metrics provides a tiny in-process Prometheus text exporter for HTTP RED.
//
// Scrape GET /metrics (local). On Vercel, process memory resets on cold start —
// prefer log-derived RED from http_request lines (route, status_class, latency_ms).
//
// Labels use Fiber route templates (e.g. /api/v1/pond/:id), never raw paths.
package metrics

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// ponytail: global mutex is fine at our QPS; upgrade to sync/atomic maps if contention shows up.
var (
	mu sync.Mutex

	requestsTotal = map[string]uint64{} // key: method|route|status_class
	panicsTotal   uint64
	dbPingFail    uint64

	// Histogram: cumulative counts per bucket + sum + count, keyed by method|route.
	durationBuckets = map[string]*hist{}
)

var bucketBounds = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

type hist struct {
	counts []uint64 // len = len(bucketBounds)+1 (+Inf)
	sum    float64
	count  uint64
}

func newHist() *hist {
	return &hist{counts: make([]uint64, len(bucketBounds)+1)}
}

// Middleware records http_requests_total and http_request_duration_seconds after the handler.
func Middleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()
		if err != nil {
			var fe *fiber.Error
			if errors.As(err, &fe) && fe.Code > 0 {
				status = fe.Code
			}
		}
		observe(c.Method(), RouteLabel(c, err), status, time.Since(start).Seconds())
		return err
	}
}

// Handler serves Prometheus text exposition format.
func Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/plain; version=0.0.4; charset=utf-8")
		return c.SendString(Exposition())
	}
}

// IncPanic increments http_panics_total (call from recover stack handler).
func IncPanic() {
	mu.Lock()
	panicsTotal++
	mu.Unlock()
}

// IncDBPingFailure increments db_ping_failures_total (ready checks).
func IncDBPingFailure() {
	mu.Lock()
	dbPingFail++
	mu.Unlock()
}

// Reset clears all series (tests only).
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	requestsTotal = map[string]uint64{}
	panicsTotal = 0
	dbPingFail = 0
	durationBuckets = map[string]*hist{}
}

// Exposition returns the current Prometheus text payload.
func Exposition() string {
	mu.Lock()
	defer mu.Unlock()

	var b strings.Builder
	b.WriteString("# HELP http_requests_total Total HTTP requests by method, route template, and status class.\n")
	b.WriteString("# TYPE http_requests_total counter\n")
	for _, k := range sortedKeys(requestsTotal) {
		method, route, class := split3(k)
		fmt.Fprintf(&b, "http_requests_total{method=%q,route=%q,status_class=%q} %d\n",
			method, route, class, requestsTotal[k])
	}

	b.WriteString("# HELP http_request_duration_seconds HTTP request latency in seconds by method and route template.\n")
	b.WriteString("# TYPE http_request_duration_seconds histogram\n")
	for _, k := range sortedHistKeys(durationBuckets) {
		method, route := split2(k)
		h := durationBuckets[k]
		var cum uint64
		for i, bound := range bucketBounds {
			cum += h.counts[i]
			fmt.Fprintf(&b, "http_request_duration_seconds_bucket{method=%q,route=%q,le=%q} %d\n",
				method, route, formatLE(bound), cum)
		}
		cum += h.counts[len(bucketBounds)]
		fmt.Fprintf(&b, "http_request_duration_seconds_bucket{method=%q,route=%q,le=\"+Inf\"} %d\n",
			method, route, cum)
		fmt.Fprintf(&b, "http_request_duration_seconds_sum{method=%q,route=%q} %s\n",
			method, route, strconv.FormatFloat(h.sum, 'f', -1, 64))
		fmt.Fprintf(&b, "http_request_duration_seconds_count{method=%q,route=%q} %d\n",
			method, route, h.count)
	}

	b.WriteString("# HELP http_panics_total Total panics recovered by the HTTP stack.\n")
	b.WriteString("# TYPE http_panics_total counter\n")
	fmt.Fprintf(&b, "http_panics_total %d\n", panicsTotal)

	b.WriteString("# HELP db_ping_failures_total Database ping failures from readiness checks.\n")
	b.WriteString("# TYPE db_ping_failures_total counter\n")
	fmt.Fprintf(&b, "db_ping_failures_total %d\n", dbPingFail)

	return b.String()
}

func observe(method, route string, status int, seconds float64) {
	class := StatusClass(status)
	key := method + "|" + route + "|" + class
	hkey := method + "|" + route

	mu.Lock()
	defer mu.Unlock()
	requestsTotal[key]++
	h := durationBuckets[hkey]
	if h == nil {
		h = newHist()
		durationBuckets[hkey] = h
	}
	h.sum += seconds
	h.count++
	placed := false
	for i, bound := range bucketBounds {
		if seconds <= bound {
			h.counts[i]++
			placed = true
			break
		}
	}
	if !placed {
		h.counts[len(bucketBounds)]++
	}
}

// RouteLabel returns the Fiber route template, or "unmatched" for 404s.
// Never use c.Path() for labels — raw paths explode cardinality.
func RouteLabel(c fiber.Ctx, err error) string {
	var fe *fiber.Error
	if errors.As(err, &fe) && fe.Code == fiber.StatusNotFound {
		return "unmatched"
	}
	if c == nil {
		return "unknown"
	}
	if r := c.Route(); r != nil && r.Path != "" {
		return r.Path
	}
	return "unmatched"
}

// StatusClass maps an HTTP status to 1xx|2xx|3xx|4xx|5xx.
func StatusClass(status int) string {
	switch {
	case status >= 500:
		return "5xx"
	case status >= 400:
		return "4xx"
	case status >= 300:
		return "3xx"
	case status >= 200:
		return "2xx"
	default:
		return "1xx"
	}
}

func formatLE(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func sortedKeys(m map[string]uint64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedHistKeys(m map[string]*hist) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func split3(k string) (method, route, class string) {
	parts := strings.SplitN(k, "|", 3)
	if len(parts) != 3 {
		return "?", "?", "?"
	}
	return parts[0], parts[1], parts[2]
}

func split2(k string) (method, route string) {
	parts := strings.SplitN(k, "|", 2)
	if len(parts) != 2 {
		return "?", "?"
	}
	return parts[0], parts[1]
}
