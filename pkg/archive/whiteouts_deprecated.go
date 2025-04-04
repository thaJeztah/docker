package archive

import "github.com/moby/go-archive"

// Whiteouts are files with a special meaning for the layered filesystem.
// Docker uses AUFS whiteout files inside exported archives. In other
// filesystems these files are generated/handled on tar creation/extraction.

// WhiteoutPrefix prefix means file is a whiteout. If this is followed by a
// filename this means that file has been removed from the base layer.
//
// Deprecated: use [archive.WhiteoutPrefix] instead.
const WhiteoutPrefix = archive.WhiteoutPrefix

// WhiteoutMetaPrefix prefix means whiteout has a special meaning and is not
// for removing an actual file. Normally these files are excluded from exported
// archives.
//
// Deprecated: use [archive.WhiteoutMetaPrefix] instead.
const WhiteoutMetaPrefix = archive.WhiteoutMetaPrefix

// WhiteoutLinkDir is a directory AUFS uses for storing hardlink links to other
// layers. Normally these should not go into exported archives and all changed
// hardlinks should be copied to the top layer.
//
// Deprecated: use [archive.WhiteoutLinkDir] instead.
const WhiteoutLinkDir = archive.WhiteoutLinkDir

// WhiteoutOpaqueDir file means directory has been made opaque - meaning
// readdir calls to this directory do not follow to lower layers.
//
// Deprecated: use [archive.WhiteoutOpaqueDir] instead.
const WhiteoutOpaqueDir = archive.WhiteoutOpaqueDir
