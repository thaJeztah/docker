package archive

import (
	"github.com/docker/docker/pkg/idtools"
	"github.com/moby/go-archive"
)

func ToArchiveOpt(options *TarOptions) *archive.TarOptions {
	if options == nil {
		return nil
	}

	var chownOpts *archive.ChownOpts
	if options.ChownOpts != nil {
		chownOpts = &archive.ChownOpts{
			UID: options.ChownOpts.UID,
			GID: options.ChownOpts.GID,
		}
	}

	return &archive.TarOptions{
		IncludeFiles:         options.IncludeFiles,
		ExcludePatterns:      options.ExcludePatterns,
		Compression:          options.Compression,
		NoLchown:             options.NoLchown,
		IDMap:                idtools.ToUserIdentityMapping(options.IDMap),
		ChownOpts:            chownOpts,
		IncludeSourceDir:     options.IncludeSourceDir,
		WhiteoutFormat:       options.WhiteoutFormat,
		NoOverwriteDirNonDir: options.NoOverwriteDirNonDir,
		RebaseNames:          options.RebaseNames,
		InUserNS:             options.InUserNS,
		BestEffortXattrs:     options.BestEffortXattrs,
	}
}
