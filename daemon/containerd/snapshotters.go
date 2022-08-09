package containerd

import "github.com/containerd/containerd"

// SnapshotterFromGraphDriver returns the containerd snapshotter name based on
// the supplied graphdriver name. It handles both legacy names and translates
// them into corresponding containerd snapshotter names.
func SnapshotterFromGraphDriver(graphDriver string) string {
	switch graphDriver {
	case "overlay", "overlay2":
		return "overlayfs"
	case "windowsfilter":
		return "windows"
	case "":
		return containerd.DefaultSnapshotter
	default:
		return graphDriver
	}
}
