package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	cerrdefs "github.com/containerd/errdefs"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestPluginSetError(t *testing.T) {
	client, err := NewClientWithOpts(WithMockClient(errorMock(http.StatusInternalServerError, "Server error")))
	assert.NilError(t, err)

	_, err = client.PluginSet(context.Background(), "plugin_name", PluginSetOptions{})
	assert.Check(t, is.ErrorType(err, cerrdefs.IsInternal))

	_, err = client.PluginSet(context.Background(), "", PluginSetOptions{})
	assert.Check(t, is.ErrorType(err, cerrdefs.IsInvalidArgument))
	assert.Check(t, is.ErrorContains(err, "value is empty"))

	_, err = client.PluginSet(context.Background(), "    ", PluginSetOptions{})
	assert.Check(t, is.ErrorType(err, cerrdefs.IsInvalidArgument))
	assert.Check(t, is.ErrorContains(err, "value is empty"))
}

func TestPluginSet(t *testing.T) {
	const expectedURL = "/plugins/plugin_name/set"

	client, err := NewClientWithOpts(WithMockClient(func(req *http.Request) (*http.Response, error) {
		if err := assertRequest(req, http.MethodPost, expectedURL); err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}, nil
	}))
	assert.NilError(t, err)

	_, err = client.PluginSet(context.Background(), "plugin_name", PluginSetOptions{Args: []string{"arg1"}})
	assert.NilError(t, err)
}
