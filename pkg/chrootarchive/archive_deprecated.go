package chrootarchive

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/go-archive/chrootarchive"
)

// NewArchiver returns a new Archiver which uses chrootarchive.Untar
//
// Deprecated: use [chrootarchive.NewArchiver] instead.
func NewArchiver(idMapping idtools.IdentityMapping) *archive.Archiver {
	return &archive.Archiver{
		Untar: func(tarArchive io.Reader, dest string, options *archive.TarOptions) error {
			return chrootarchive.Untar(tarArchive, dest, archive.ToArchiveOpt(options))
		},
		IDMapping: idMapping,
	}
}

// Untar reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive may be compressed with one of the following algorithms:
// identity (uncompressed), gzip, bzip2, xz.
//
// Deprecated: use [chrootarchive.Untar] instead.
func Untar(tarArchive io.Reader, dest string, options *archive.TarOptions) error {
	return chrootarchive.Untar(tarArchive, dest, archive.ToArchiveOpt(options))
}

// UntarWithRoot is the same as `Untar`, but allows you to pass in a root directory
// The root directory is the directory that will be chrooted to.
// `dest` must be a path within `root`, if it is not an error will be returned.
//
// `root` should set to a directory which is not controlled by any potentially
// malicious process.
//
// This should be used to prevent a potential attacker from manipulating `dest`
// such that it would provide access to files outside of `dest` through things
// like symlinks. Normally `ResolveSymlinksInScope` would handle this, however
// sanitizing symlinks in this manner is inherently racey:
// ref: CVE-2018-15664
//
// Deprecated: use [chrootarchive.UntarWithRoot] instead.
func UntarWithRoot(tarArchive io.Reader, dest string, options *archive.TarOptions, root string) error {
	return chrootarchive.UntarWithRoot(tarArchive, dest, archive.ToArchiveOpt(options), root)
}

// UntarUncompressed reads a stream of bytes from `archive`, parses it as a tar archive,
// and unpacks it into the directory at `dest`.
// The archive must be an uncompressed stream.
//
// Deprecated: use [chrootarchive.UntarUncompressed] instead.
func UntarUncompressed(tarArchive io.Reader, dest string, options *archive.TarOptions) error {
	return chrootarchive.UntarUncompressed(tarArchive, dest, archive.ToArchiveOpt(options))
}

// Tar tars the requested path while chrooted to the specified root.
//
// Deprecated: use [chrootarchive.Tar] instead.
func Tar(srcPath string, options *archive.TarOptions, root string) (io.ReadCloser, error) {
	return chrootarchive.Tar(srcPath, archive.ToArchiveOpt(options), root)
}
