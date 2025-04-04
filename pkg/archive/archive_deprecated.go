// Package archive provides helper functions for dealing with archive files.
package archive

import (
	"archive/tar"
	"io"
	"os"

	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/go-archive"
)

// ImpliedDirectoryMode represents the mode (Unix permissions) applied to directories that are implied by files in a
// tar, but that do not have their own header entry.
//
// The permissions mask is stored in a constant instead of locally to ensure that magic numbers do not
// proliferate in the codebase. The default value 0755 has been selected based on the default umask of 0022, and
// a convention of mkdir(1) calling mkdir(2) with permissions of 0777, resulting in a final value of 0755.
//
// This value is currently implementation-defined, and not captured in any cross-runtime specification. Thus, it is
// subject to change in Moby at any time -- image authors who require consistent or known directory permissions
// should explicitly control them by ensuring that header entries exist for any applicable path.
//
// Deprecated: use [archive.ImpliedDirectoryMode] instead.
const ImpliedDirectoryMode = archive.ImpliedDirectoryMode

type (
	// Compression is the state represents if compressed or not.
	//
	// Deprecated: use [archive.Compression] instead.
	Compression = archive.Compression
	// WhiteoutFormat is the format of whiteouts unpacked
	//
	// Deprecated: use [archive.WhiteoutFormat] instead.
	WhiteoutFormat = archive.WhiteoutFormat

	// TarOptions wraps the tar options.
	TarOptions struct {
		IncludeFiles     []string
		ExcludePatterns  []string
		Compression      archive.Compression
		NoLchown         bool
		IDMap            idtools.IdentityMapping
		ChownOpts        *idtools.Identity
		IncludeSourceDir bool
		// WhiteoutFormat is the expected on disk format for whiteout files.
		// This format will be converted to the standard format on pack
		// and from the standard format on unpack.
		WhiteoutFormat archive.WhiteoutFormat
		// When unpacking, specifies whether overwriting a directory with a
		// non-directory is allowed and vice versa.
		NoOverwriteDirNonDir bool
		// For each include when creating an archive, the included name will be
		// replaced with the matching name from this map.
		RebaseNames map[string]string
		InUserNS    bool
		// Allow unpacking to succeed in spite of failures to set extended
		// attributes on the unpacked files due to the destination filesystem
		// not supporting them or a lack of permissions. Extended attributes
		// were probably in the archive for a reason, so set this option at
		// your own peril.
		BestEffortXattrs bool
	}
)

// Archiver implements the Archiver interface and allows the reuse of most utility functions of
// this package with a pluggable Untar function. Also, to facilitate the passing of specific id
// mappings for untar, an Archiver can be created with maps which will then be passed to Untar operations.
//
// Deprecated: use [archive.Archiver] instead.
type Archiver struct {
	Untar     func(io.Reader, string, *TarOptions) error
	IDMapping idtools.IdentityMapping
}

// NewDefaultArchiver returns a new Archiver without any IdentityMapping
//
// Deprecated: use [archive.NewDefaultArchiver] instead.
func NewDefaultArchiver() *Archiver {
	return &Archiver{Untar: func(r io.Reader, path string, tarOptions *TarOptions) error {
		return archive.Untar(r, path, ToArchiveOpt(tarOptions))
	}}
}

const (
	Uncompressed = archive.Uncompressed // Uncompressed represents the uncompressed.
	Bzip2        = archive.Bzip2        // Bzip2 is bzip2 compression algorithm.
	Gzip         = archive.Gzip         // Gzip is gzip compression algorithm.
	Xz           = archive.Xz           // Xz is xz compression algorithm.
	Zstd         = archive.Zstd         // Zstd is zstd compression algorithm.
)

const (
	AUFSWhiteoutFormat    = archive.AUFSWhiteoutFormat    // AUFSWhiteoutFormat is the default format for whiteouts
	OverlayWhiteoutFormat = archive.OverlayWhiteoutFormat // OverlayWhiteoutFormat formats whiteout according to the overlay standard.
)

// IsArchivePath checks if the (possibly compressed) file at the given path
// starts with a tar file header.
//
// Deprecated: use [archive.IsArchivePath] instead.
func IsArchivePath(path string) bool {
	return archive.IsArchivePath(path)
}

// DetectCompression detects the compression algorithm of the source.
//
// Deprecated: use [archive.DetectCompression] instead.
func DetectCompression(source []byte) archive.Compression {
	return archive.DetectCompression(source)
}

// DecompressStream decompresses the archive and returns a ReaderCloser with the decompressed archive.
//
// Deprecated: use [archive.DecompressStream] instead.
func DecompressStream(arch io.Reader) (io.ReadCloser, error) {
	return archive.DecompressStream(arch)
}

// CompressStream compresses the dest with specified compression algorithm.
//
// Deprecated: use [archive.CompressStream] instead.
func CompressStream(dest io.Writer, compression archive.Compression) (io.WriteCloser, error) {
	return archive.CompressStream(dest, compression)
}

// TarModifierFunc is a function that can be passed to ReplaceFileTarWrapper to
// modify the contents or header of an entry in the archive. If the file already
// exists in the archive the TarModifierFunc will be called with the Header and
// a reader which will return the files content. If the file does not exist both
// header and content will be nil.
//
// Deprecated: use [archive.TarModifierFunc] instead.
type TarModifierFunc = archive.TarModifierFunc

// ReplaceFileTarWrapper converts inputTarStream to a new tar stream. Files in the
// tar stream are modified if they match any of the keys in mods.
//
// Deprecated: use [archive.ReplaceFileTarWrapper] instead.
func ReplaceFileTarWrapper(inputTarStream io.ReadCloser, mods map[string]archive.TarModifierFunc) io.ReadCloser {
	return archive.ReplaceFileTarWrapper(inputTarStream, mods)
}

// FileInfoHeaderNoLookups creates a partially-populated tar.Header from fi.
//
// Compared to the archive/tar.FileInfoHeader function, this function is safe to
// call from a chrooted process as it does not populate fields which would
// require operating system lookups. It behaves identically to
// tar.FileInfoHeader when fi is a FileInfo value returned from
// tar.Header.FileInfo().
//
// When fi is a FileInfo for a native file, such as returned from os.Stat() and
// os.Lstat(), the returned Header value differs from one returned from
// tar.FileInfoHeader in the following ways. The Uname and Gname fields are not
// set as OS lookups would be required to populate them. The AccessTime and
// ChangeTime fields are not currently set (not yet implemented) although that
// is subject to change. Callers which require the AccessTime or ChangeTime
// fields to be zeroed should explicitly zero them out in the returned Header
// value to avoid any compatibility issues in the future.
//
// Deprecated: use [archive.FileInfoHeaderNoLookups] instead.
func FileInfoHeaderNoLookups(fi os.FileInfo, link string) (*tar.Header, error) {
	return archive.FileInfoHeaderNoLookups(fi, link)
}

// FileInfoHeader creates a populated Header from fi.
//
// Compared to the archive/tar package, this function fills in less information
// but is safe to call from a chrooted process. The AccessTime and ChangeTime
// fields are not set in the returned header, ModTime is truncated to one-second
// precision, and the Uname and Gname fields are only set when fi is a FileInfo
// value returned from tar.Header.FileInfo().
//
// Deprecated: use [archive.FileInfoHeader] instead.
func FileInfoHeader(name string, fi os.FileInfo, link string) (*tar.Header, error) {
	return archive.FileInfoHeader(name, fi, link)
}

// ReadSecurityXattrToTarHeader reads security.capability xattr from filesystem
// to a tar header
//
// Deprecated: use [archive.ReadSecurityXattrToTarHeader] instead.
func ReadSecurityXattrToTarHeader(path string, hdr *tar.Header) error {
	return archive.ReadSecurityXattrToTarHeader(path, hdr)
}

// Tar creates an archive from the directory at `path`, and returns it as a
// stream of bytes.
//
// Deprecated: use [archive.Tar] instead.
func Tar(path string, compression archive.Compression) (io.ReadCloser, error) {
	return archive.TarWithOptions(path, &archive.TarOptions{Compression: compression})
}

// TarWithOptions creates an archive from the directory at `path`, only including files whose relative
// paths are included in `options.IncludeFiles` (if non-nil) or not in `options.ExcludePatterns`.
//
// Deprecated: use [archive.TarWithOptions] instead.
func TarWithOptions(srcPath string, options *TarOptions) (io.ReadCloser, error) {
	return archive.TarWithOptions(srcPath, ToArchiveOpt(options))
}

// Tarballer is a lower-level interface to TarWithOptions which gives the caller
// control over which goroutine the archiving operation executes on.
//
// Deprecated: use [archive.Tarballer] instead.
type Tarballer = archive.Tarballer

// NewTarballer constructs a new tarballer. The arguments are the same as for
// TarWithOptions.
//
// Deprecated: use [archive.Tarballer] instead.
func NewTarballer(srcPath string, options *TarOptions) (*archive.Tarballer, error) {
	return archive.NewTarballer(srcPath, ToArchiveOpt(options))
}

// Unpack unpacks the decompressedArchive to dest with options.
//
// Deprecated: use [archive.Unpack] instead.
func Unpack(decompressedArchive io.Reader, dest string, options *TarOptions) error {
	return archive.Unpack(decompressedArchive, dest, ToArchiveOpt(options))
}

// Untar reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive may be compressed with one of the following algorithms:
// identity (uncompressed), gzip, bzip2, xz.
//
// FIXME: specify behavior when target path exists vs. doesn't exist.
//
// Deprecated: use [archive.Untar] instead.
func Untar(tarArchive io.Reader, dest string, options *TarOptions) error {
	return archive.Untar(tarArchive, dest, ToArchiveOpt(options))
}

// UntarUncompressed reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive must be an uncompressed stream.
//
// Deprecated: use [archive.UntarUncompressed] instead.
func UntarUncompressed(tarArchive io.Reader, dest string, options *TarOptions) error {
	return archive.UntarUncompressed(tarArchive, dest, ToArchiveOpt(options))
}

// TarUntar is a convenience function which calls Tar and Untar, with the output of one piped into the other.
// If either Tar or Untar fails, TarUntar aborts and returns the error.
func (archiver *Archiver) TarUntar(src, dst string) error {
	return (&archive.Archiver{
		Untar: func(reader io.Reader, s string, options *archive.TarOptions) error {
			return archiver.Untar(reader, s, &TarOptions{
				IDMap: archiver.IDMapping,
			})
		},
		IDMapping: idtools.ToUserIdentityMapping(archiver.IDMapping),
	}).TarUntar(src, dst)
}

// UntarPath untar a file from path to a destination, src is the source tar file path.
func (archiver *Archiver) UntarPath(src, dst string) error {
	return (&archive.Archiver{
		Untar: func(reader io.Reader, s string, options *archive.TarOptions) error {
			return archiver.Untar(reader, s, &TarOptions{
				IDMap: archiver.IDMapping,
			})
		},
		IDMapping: idtools.ToUserIdentityMapping(archiver.IDMapping),
	}).UntarPath(src, dst)
}

// CopyWithTar creates a tar archive of filesystem path `src`, and
// unpacks it at filesystem path `dst`.
// The archive is streamed directly with fixed buffering and no
// intermediary disk IO.
func (archiver *Archiver) CopyWithTar(src, dst string) error {
	return (&archive.Archiver{
		Untar: func(reader io.Reader, s string, options *archive.TarOptions) error {
			return archiver.Untar(reader, s, nil)
		},
		IDMapping: idtools.ToUserIdentityMapping(archiver.IDMapping),
	}).CopyWithTar(src, dst)
}

// CopyFileWithTar emulates the behavior of the 'cp' command-line
// for a single file. It copies a regular file from path `src` to
// path `dst`, and preserves all its metadata.
func (archiver *Archiver) CopyFileWithTar(src, dst string) (err error) {
	return (&archive.Archiver{
		Untar: func(reader io.Reader, s string, options *archive.TarOptions) error {
			return archiver.Untar(reader, s, nil)
		},
		IDMapping: idtools.ToUserIdentityMapping(archiver.IDMapping),
	}).CopyFileWithTar(src, dst)
}

// IdentityMapping returns the IdentityMapping of the archiver.
func (archiver *Archiver) IdentityMapping() idtools.IdentityMapping {
	return archiver.IDMapping
}
