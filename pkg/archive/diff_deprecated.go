package archive

import (
	"io"

	"github.com/moby/go-archive"
)

// UnpackLayer unpack `layer` to a `dest`. The stream `layer` can be
// compressed or uncompressed.
// Returns the size in bytes of the contents of the layer.
//
// Deprecated: use [archive.UnpackLayer] instead.
func UnpackLayer(dest string, layer io.Reader, options *TarOptions) (size int64, err error) {
	return archive.UnpackLayer(dest, layer, ToArchiveOpt(options))
}

// ApplyLayer parses a diff in the standard layer format from `layer`,
// and applies it to the directory `dest`. The stream `layer` can be
// compressed or uncompressed.
// Returns the size in bytes of the contents of the layer.
//
// Deprecated: use [archive.ApplyLayer] instead.
func ApplyLayer(dest string, layer io.Reader) (int64, error) {
	return archive.ApplyLayer(dest, layer)
}

// ApplyUncompressedLayer parses a diff in the standard layer format from
// `layer`, and applies it to the directory `dest`. The stream `layer`
// can only be uncompressed.
// Returns the size in bytes of the contents of the layer.
func ApplyUncompressedLayer(dest string, layer io.Reader, options *TarOptions) (int64, error) {
	return archive.ApplyUncompressedLayer(dest, layer, ToArchiveOpt(options))
}

// IsEmpty checks if the tar archive is empty (doesn't contain any entries).
//
// Deprecated: use [archive.IsEmpty] instead.
func IsEmpty(rd io.Reader) (bool, error) {
	return archive.IsEmpty(rd)
}
