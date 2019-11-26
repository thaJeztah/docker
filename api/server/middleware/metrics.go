package middleware

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/server/httputils"

	"github.com/docker/go-metrics"
)

func init() {
}

type MetricsMiddleware struct {
	apiRequestCounter metrics.LabeledCounter
}

func NewMetricsMiddleware() MetricsMiddleware {
	ns := metrics.NewNamespace("engine", "api", nil)
	arc := ns.NewLabeledCounter("http_requests_total", "The total number of API requests", "api", "action", "apiVersion", "method", "code")
	metrics.Register(ns)

	return MetricsMiddleware{apiRequestCounter: arc}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (m MetricsMiddleware) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		m.apiRequestCounter.WithValues(
			"containers",                      // api
			r.URL.Path,                        // action
			httputils.VersionFromContext(ctx), // apiVersion
			r.Method,                          // method
			"200",                             // code
		).Inc()
		return handler(ctx, w, r, vars)
	}
}
