package archive

import (
	"io"

	"github.com/moby/go-archive"
)

// Generate generates a new archive from the content provided
// as input.
//
// `files` is a sequence of path/content pairs. A new file is
// added to the archive for each pair.
// If the last pair is incomplete, the file is created with an
// empty content. For example:
//
// Generate("foo.txt", "hello world", "emptyfile")
//
// The above call will return an archive with 2 files:
//   - ./foo.txt with content "hello world"
//   - ./empty with empty content
//
// FIXME: stream content instead of buffering
// FIXME: specify permissions and other archive metadata
//
// Deprecated: use [archive.Generate] instead.
func Generate(input ...string) (io.Reader, error) {
	return archive.Generate(input...)
}
