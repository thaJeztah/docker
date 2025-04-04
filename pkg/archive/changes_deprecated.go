package archive

import (
	"io"

	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/go-archive"
)

// ChangeType represents the change
//
// Deprecated: use [archive.ChangeType] instead.
type ChangeType = archive.ChangeType

const (
	ChangeModify = archive.ChangeModify // ChangeModify represents the modify operation.
	ChangeAdd    = archive.ChangeAdd    // ChangeAdd represents the add operation.
	ChangeDelete = archive.ChangeDelete // ChangeDelete represents the delete operation.
)

// Change represents a change, it wraps the change type and path.
// It describes changes of the files in the path respect to the
// parent layers. The change could be modify, add, delete.
// This is used for layer diff.
//
// Deprecated: use [archive.Change] instead.
type Change = archive.Change

// Changes walks the path rw and determines changes for the files in the path,
// with respect to the parent layers
//
// Deprecated: use [archive.Changes] instead.
func Changes(layers []string, rw string) ([]archive.Change, error) {
	return archive.Changes(layers, rw)
}

// FileInfo describes the information of a file.
//
// Deprecated: use [archive.FileInfo] instead.
type FileInfo = archive.FileInfo

// ChangesDirs compares two directories and generates an array of Change objects describing the changes.
// If oldDir is "", then all files in newDir will be Add-Changes.
//
// Deprecated: use [archive.ChangesDirs] instead.
func ChangesDirs(newDir, oldDir string) ([]archive.Change, error) {
	return archive.ChangesDirs(newDir, oldDir)
}

// ChangesSize calculates the size in bytes of the provided changes, based on newDir.
//
// Deprecated: use [archive.ChangesSize] instead.
func ChangesSize(newDir string, changes []archive.Change) int64 {
	return archive.ChangesSize(newDir, changes)
}

// ExportChanges produces an Archive from the provided changes, relative to dir.
func ExportChanges(dir string, changes []archive.Change, idMap idtools.IdentityMapping) (io.ReadCloser, error) {
	return archive.ExportChanges(dir, changes, idtools.ToUserIdentityMapping(idMap))
}
