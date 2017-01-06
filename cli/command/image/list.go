package image

import (
	"fmt"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/opts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types/filters"
)

type imagesOptions struct {
	matchName string

	quiet       bool
	all         bool
	noTrunc     bool
	showDigests bool
	format      string
	filter      opts.FilterOpt
}

// NewImagesCommand creates a new `docker images` command
func NewImagesCommand(dockerCli *command.DockerCli) *cobra.Command {
	opts := imagesOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:   "images [OPTIONS] [REPOSITORY[:TAG]]",
		Short: "List images",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.matchName = args[0]
			}
			return runImages(dockerCli, opts)
		},
	}

	flags := cmd.Flags()

	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only show numeric IDs")
	flags.BoolVarP(&opts.all, "all", "a", false, "Show all images (default hides intermediate images)")
	flags.BoolVar(&opts.noTrunc, "no-trunc", false, "Don't truncate output")
	flags.BoolVar(&opts.showDigests, "digests", false, "Show digests")
	flags.StringVar(&opts.format, "format", "", "Pretty-print images using a Go template")
	flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")

	return cmd
}

func newListCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := *NewImagesCommand(dockerCli)
	cmd.Aliases = []string{"images", "list"}
	cmd.Use = "ls [OPTIONS] [REPOSITORY[:TAG]]"
	return &cmd
}

func runImages(dockerCli *command.DockerCli, opts imagesOptions) error {
	ctx := context.Background()

	filters := opts.filter.Value()
	if opts.matchName != "" {
		filters.Add("reference", opts.matchName)
	}

	options := types.ImageListOptions{
		All:     opts.all,
		Filters: filters,
	}

	images, err := dockerCli.Client().ImageList(ctx, options)
	if err != nil {
		return err
	}

	format := opts.format
	if len(format) == 0 {
		if len(dockerCli.ConfigFile().ImagesFormat) > 0 && !opts.quiet {
			format = dockerCli.ConfigFile().ImagesFormat
		} else {
			format = formatter.TableFormatKey
		}
	}

	// Images returned by the API have all references (tags, digests) included.
	// When _displaying_ those images, each reference is printed on a new line,
	// so what is seen by the user as "image" is actually an image _reference_.
	//
	// When filtering by reference ("repo:tag", or "repo@digest"), the API matches
	// images on those criteria, but always returns all references that are present
	// for an image.
	//
	// The code below filters
	// returns
	//

	// problem here is that all _tags_ that match should be printed
	// and the RepoDigest should only be printed if it matches that tag/image
	// if searching by _digest_, only images with that digest should be shown
	// but we should preserve the tag (and not print <none>)
	// if the _name_ of the image matches, the digest should be shown if the name in the digest matches a name
	// if the digest _itself_ was matched, it should be shown separately
	if filters.Include("reference") {
		for i, img := range images {
			tagsByName := map[string][]string{}
			tagMatches := map[string][]string{}
			tags := []string{}
			digests := []string{}

			for _, refString := range img.RepoTags {
				ref, err := reference.ParseNamed(refString)
				if err != nil {
					continue
				}
				if _, ok := ref.(reference.NamedTagged); !ok {
					continue
				}
				tagsByName[ref.Name()] = append(tagsByName[ref.Name()], ref.String())

				for _, pattern := range filters.Get("reference") {
					found, matchErr := reference.Match(pattern, ref)
					if matchErr != nil || !found {
						continue
					}
					// Matched by name/tag
					tags = append(tags, ref.String())
					tagMatches[ref.Name()] = append(tagMatches[ref.Name()], ref.String())
				}
			}
			for _, refString := range img.RepoDigests {
				ref, err := reference.ParseNamed(refString)
				if err != nil {
					continue
				}
				if _, ok := ref.(reference.Canonical); !ok {
					continue
				}
				if _, ok := tagMatches[ref.Name()]; ok {
					// Reference by tag found. Store the corresponding digest for presentation
					digests = append(digests, ref.String())
					continue
				}

				for _, pattern := range filters.Get("reference") {
					found, matchErr := reference.Match(pattern, ref)
					if matchErr != nil || !found {
						continue
					}
					// Matched by digest
					digests = append(digests, ref.String())

					// Also append any tag associated with this digest (but only for the same image name)
					tags = append(tags, tagsByName[ref.Name()]...)
				}
			}

			// Now, include in the image;
			// - any tag that was matched
			// - any digest that was matched
			// - anny tag for images that were matched by digest
			// - anny digest that matches an image-name that was matched
			// don't worry about duplicates; they are de-duplicated during printing
			img.RepoTags = nil
			img.RepoDigests = nil
			img.RepoTags = append(img.RepoTags, tags...)
			img.RepoDigests = append(img.RepoDigests, digests...)
			images[i] = img
		}
	}


	fmt.Println("IMAGES", images)


	imageCtx := formatter.ImageContext{
		Context: formatter.Context{
			Output: dockerCli.Out(),
			Format: formatter.NewImageFormat(format, opts.quiet, opts.showDigests),
			Trunc:  !opts.noTrunc,
		},
		Digest: opts.showDigests,
	}
	return formatter.ImageWrite(imageCtx, images)
}

// getFilteredImages materializes the specified images
//
// Images returned by the API have all references (tags, digests) included.
// When _displaying_ those images, each reference is printed on a new line,
// so what is seen by the user as "image" is actually an image _reference_.
//
// When filtering by reference ("repo:tag", or "repo@digest"), the API matches
// images on those criteria, but always returns all references that are present
// for an image.
//
// The code below filters
// returns
//

// problem here is that all _tags_ that match should be printed
// and the RepoDigest should only be printed if it matches that tag/image
// if searching by _digest_, only images with that digest should be shown
// but we should preserve the tag (and not print <none>)
// if the _name_ of the image matches, the digest should be shown if the name in the digest matches a name
// if the digest _itself_ was matched, it should be shown separately

func getFilteredImages(images []types.ImageSummary, references []string) ([]types.imageContext, error) {

}

func materializeImage(img types.ImageSummary) ([]types.imageContext) {
	tagsByName := map[string][]string{}
	digestsByName := map[string][]string{}
	images := []types.imageContext{}

	for _, refString := range append(img.RepoTags, img.RepoDigests...) {
		ref, err := reference.ParseNamed(refString)
		if err != nil {
			continue
		}
		if _, ok := ref.(reference.NamedTagged); ok {
			tagsByName[ref.Name()] = append(tagsByName[ref.Name()], ref.String())
			continue
		}
		if _, ok := ref.(reference.Canonical); !ok {
			digestsByName[ref.Name()] = append(digestsByName[ref.Name()], ref.String())
			continue
		}
	}

	for reponame, tags := range tagsByName {
		images = append(images, &imageContext{
			trunc:  ctx.Trunc,
			i:      image,
			repo:   repo,
			tag:    tag,
			digest: "<none>",
		})

	}

}
