package chrootarchive

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/moby/go-archive/chrootarchive"
)

// ApplyLayer parses a diff in the standard layer format from `layer`,
// and applies it to the directory `dest`. The stream `layer` can only be
// uncompressed.
// Returns the size in bytes of the contents of the layer.
//
// Deprecated: use [chrootarchive.ApplyLayer] insteead.
func ApplyLayer(dest string, layer io.Reader) (size int64, err error) {
	return chrootarchive.ApplyLayer(dest, layer)
}

// ApplyUncompressedLayer parses a diff in the standard layer format from
// `layer`, and applies it to the directory `dest`. The stream `layer`
// can only be uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyUncompressedLayer(dest string, layer io.Reader, options *archive.TarOptions) (int64, error) {
	return chrootarchive.ApplyUncompressedLayer(dest, layer, archive.ToArchiveOpt(options))
}
