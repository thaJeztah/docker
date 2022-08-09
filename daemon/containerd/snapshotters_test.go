package containerd

import (
	"testing"

	"github.com/containerd/containerd"
)

func TestSnapshotterFromGraphDriver(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output string
	}{
		{
			desc:   "empty defaults to containerd default",
			input:  "",
			output: containerd.DefaultSnapshotter,
		},
		{
			desc:   "overlay -> overlayfs",
			input:  "overlay",
			output: "overlayfs",
		},
		{
			desc:   "overlay2 -> overlayfs",
			input:  "overlay2",
			output: "overlayfs",
		},
		{
			desc:   "windowsfilter -> windows",
			input:  "windowsfilter",
			output: "windows",
		},
		{
			desc:   "containerd overlayfs",
			input:  "overlayfs",
			output: "overlayfs",
		},
		{
			desc:   "containerd zfs",
			input:  "zfs",
			output: "zfs",
		},
		{
			desc:   "unknown is unchanged",
			input:  "somefuturesnapshotter",
			output: "somefuturesnapshotter",
		},
	}
	for _, tC := range testCases {
		want := tC.output
		t.Run(tC.desc, func(t *testing.T) {
			got := SnapshotterFromGraphDriver(tC.input)
			if want != got {
				t.Errorf("Expected sanitizeGraphDriver to return %q, got %q", want, got)
			}
		})
	}
}
