package containerd

import (
	"testing"

	"github.com/containerd/containerd"
	"gotest.tools/v3/assert"
)

func TestSnapshotterFromGraphDriver(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "empty defaults to containerd default",
			input:    "",
			expected: containerd.DefaultSnapshotter,
		},
		{
			desc:     "overlay -> overlayfs",
			input:    "overlay",
			expected: "overlayfs",
		},
		{
			desc:     "overlay2 -> overlayfs",
			input:    "overlay2",
			expected: "overlayfs",
		},
		{
			desc:     "windowsfilter -> windows",
			input:    "windowsfilter",
			expected: "windows",
		},
		{
			desc:     "containerd overlayfs",
			input:    "overlayfs",
			expected: "overlayfs",
		},
		{
			desc:     "containerd zfs",
			input:    "zfs",
			expected: "zfs",
		},
		{
			desc:     "unknown is unchanged",
			input:    "somefuturesnapshotter",
			expected: "somefuturesnapshotter",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, SnapshotterFromGraphDriver(tc.input), tc.expected)
		})
	}
}
