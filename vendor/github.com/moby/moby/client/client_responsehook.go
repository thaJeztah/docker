package client

import (
	"net/http"
)

type responseHookTransport struct {
	base  http.RoundTripper
	hooks []ResponseHook
}

func (t *responseHookTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	for _, h := range t.hooks {
		if h == nil {
			continue
		}
		if err := h(resp); err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
	}

	return resp, nil
}
