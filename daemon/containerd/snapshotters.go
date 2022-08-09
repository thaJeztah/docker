package containerd

import (
	"strings"

	"github.com/containerd/containerd"
)

// SnapshotterFromGraphDriver returns the containerd snapshotter name based on
// the supplied graphdriver name. It handles both legacy names and translates
// them into corresponding containerd snapshotter names.
func SnapshotterFromGraphDriver(graphDriver string) string {
	if graphDriver == "" {
		graphDriver = containerd.DefaultSnapshotter
	}

	switch graphDriver {
	case "overlay", "overlay2":
		graphDriver = "overlayfs"
	case "windowsgraph":
		graphDriver = "windows"
	}

	graphDriver = strings.TrimPrefix(graphDriver, "io.containerd.snapshotter.v1.")
	return graphDriver
}
