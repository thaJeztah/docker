package archive

import (
	"io"
	"path/filepath"

	"github.com/moby/go-archive"
)

// Errors used or returned by this file.
var (
	ErrNotDirectory      = archive.ErrNotDirectory      // Deprecated: use [archive.ErrNotDirectory] instead.
	ErrDirNotExists      = archive.ErrDirNotExists      // Deprecated: use [archive.ErrDirNotExists] instead.
	ErrCannotCopyDir     = archive.ErrCannotCopyDir     // Deprecated: use [archive.ErrCannotCopyDir] instead.
	ErrInvalidCopySource = archive.ErrInvalidCopySource // Deprecated: use [archive.ErrInvalidCopySource] instead.
)

// PreserveTrailingDotOrSeparator returns the given cleaned path (after
// processing using any utility functions from the path or filepath stdlib
// packages) and appends a trailing `/.` or `/` if its corresponding  original
// path (from before being processed by utility functions from the path or
// filepath stdlib packages) ends with a trailing `/.` or `/`. If the cleaned
// path already ends in a `.` path segment, then another is not added. If the
// clean path already ends in a path separator, then another is not added.
//
// Deprecated: use [archive.PreserveTrailingDotOrSeparator] instead.
func PreserveTrailingDotOrSeparator(cleanedPath string, originalPath string) string {
	return archive.PreserveTrailingDotOrSeparator(cleanedPath, originalPath)
}

// hasTrailingPathSeparator returns whether the given
// path ends with the system's path separator character.
func hasTrailingPathSeparator(path string) bool {
	return len(path) > 0 && path[len(path)-1] == filepath.Separator
}

// specifiesCurrentDir returns whether the given path specifies
// a "current directory", i.e., the last path segment is `.`.
func specifiesCurrentDir(path string) bool {
	return filepath.Base(path) == "."
}

// SplitPathDirEntry splits the given path between its directory name and its
// basename by first cleaning the path but preserves a trailing "." if the
// original path specified the current directory.
//
// Deprecated: use [archive.SplitPathDirEntry] instead.
func SplitPathDirEntry(path string) (dir, base string) {
	return archive.SplitPathDirEntry(path)
}

// TarResource archives the resource described by the given CopyInfo to a Tar
// archive. A non-nil error is returned if sourcePath does not exist or is
// asserted to be a directory but exists as another type of file.
//
// This function acts as a convenient wrapper around TarWithOptions, which
// requires a directory as the source path. TarResource accepts either a
// directory or a file path and correctly sets the Tar options.
//
// Deprecated: use [archive.TarResource] instead.
func TarResource(sourceInfo archive.CopyInfo) (content io.ReadCloser, err error) {
	return archive.TarResource(sourceInfo)
}

// TarResourceRebase is like TarResource but renames the first path element of
// items in the resulting tar archive to match the given rebaseName if not "".
//
// Deprecated: use [archive.TarResourceRebase] instead.
func TarResourceRebase(sourcePath, rebaseName string) (content io.ReadCloser, _ error) {
	return archive.TarResourceRebase(sourcePath, rebaseName)
}

// TarResourceRebaseOpts does not preform the Tar, but instead just creates the rebase
// parameters to be sent to TarWithOptions (the TarOptions struct)
func TarResourceRebaseOpts(sourceBase string, rebaseName string) *TarOptions {
	filter := []string{sourceBase}
	return &TarOptions{
		Compression:      archive.Uncompressed,
		IncludeFiles:     filter,
		IncludeSourceDir: true,
		RebaseNames: map[string]string{
			sourceBase: rebaseName,
		},
	}
}

// CopyInfo holds basic info about the source
// or destination path of a copy operation.
//
// Deprecated: use [archive.CopyInfo] instead.
type CopyInfo = archive.CopyInfo

// CopyInfoSourcePath stats the given path to create a CopyInfo
// struct representing that resource for the source of an archive copy
// operation. The given path should be an absolute local path. A source path
// has all symlinks evaluated that appear before the last path separator ("/"
// on Unix). As it is to be a copy source, the path must exist.
//
// Deprecated: use [archive.CopyInfoSourcePath] instead.
func CopyInfoSourcePath(path string, followLink bool) (archive.CopyInfo, error) {
	return archive.CopyInfoSourcePath(path, followLink)
}

// CopyInfoDestinationPath stats the given path to create a CopyInfo
// struct representing that resource for the destination of an archive copy
// operation. The given path should be an absolute local path.
//
// Deprecated: use [archive.CopyInfoDestinationPath] instead.
func CopyInfoDestinationPath(path string) (info archive.CopyInfo, err error) {
	return archive.CopyInfoDestinationPath(path)
}

// PrepareArchiveCopy prepares the given srcContent archive, which should
// contain the archived resource described by srcInfo, to the destination
// described by dstInfo. Returns the possibly modified content archive along
// with the path to the destination directory which it should be extracted to.
//
// Deprecated: use [archive.PrepareArchiveCopy] instead.
func PrepareArchiveCopy(srcContent io.Reader, srcInfo, dstInfo archive.CopyInfo) (dstDir string, content io.ReadCloser, err error) {
	return archive.PrepareArchiveCopy(srcContent, srcInfo, dstInfo)
}

// RebaseArchiveEntries rewrites the given srcContent archive replacing
// an occurrence of oldBase with newBase at the beginning of entry names.
//
// Deprecated: use [archive.RebaseArchiveEntries] instead.
func RebaseArchiveEntries(srcContent io.Reader, oldBase, newBase string) io.ReadCloser {
	return archive.RebaseArchiveEntries(srcContent, oldBase, newBase)
}

// CopyResource performs an archive copy from the given source path to the
// given destination path. The source path MUST exist and the destination
// path's parent directory must exist.
//
// Deprecated: use [archive.CopyResource] instead.
func CopyResource(srcPath, dstPath string, followLink bool) error {
	return archive.CopyResource(srcPath, dstPath, followLink)
}

// CopyTo handles extracting the given content whose
// entries should be sourced from srcInfo to dstPath.
//
// Deprecated: use [archive.CopyTo] instead.
func CopyTo(content io.Reader, srcInfo archive.CopyInfo, dstPath string) error {
	return archive.CopyTo(content, srcInfo, dstPath)
}

// ResolveHostSourcePath decides real path need to be copied with parameters such as
// whether to follow symbol link or not, if followLink is true, resolvedPath will return
// link target of any symbol link file, else it will only resolve symlink of directory
// but return symbol link file itself without resolving.
//
// Deprecated: use [archive.ResolveHostSourcePath] instead.
func ResolveHostSourcePath(path string, followLink bool) (resolvedPath, rebaseName string, _ error) {
	return archive.ResolveHostSourcePath(path, followLink)
}

// GetRebaseName normalizes and compares path and resolvedPath,
// return completed resolved path and rebased file name
//
// Deprecated: use [archive.GetRebaseName] instead.
func GetRebaseName(path, resolvedPath string) (string, string) {
	return archive.GetRebaseName(path, resolvedPath)
}
