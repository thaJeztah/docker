package images // import "github.com/docker/docker/daemon/images"

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/dockerversion"
)

var acceptedSearchFilterTags = map[string]bool{
	"is-automated": true,
	"is-official":  true,
	"stars":        true,
}

// SearchRegistryForImages queries the registry for images matching
// term. authConfig is used to login.
//
// TODO: this could be implemented in a registry service instead of the image
// service.
func (i *ImageService) SearchRegistryForImages(ctx context.Context, filtersArgs string, term string, limit int,
	authConfig *types.AuthConfig,
	headers map[string][]string) (*registrytypes.SearchResults, error) {

	searchFilters, err := filters.FromJSON(filtersArgs)
	if err != nil {
		return nil, err
	}
	if err := searchFilters.Validate(acceptedSearchFilterTags); err != nil {
		return nil, err
	}

	var isAutomated, isOfficial bool
	var hasStarFilter = 0
	if searchFilters.Contains("is-automated") {
		switch {
		case searchFilters.UniqueExactMatch("is-automated", "true"):
			isAutomated = true
		case searchFilters.UniqueExactMatch("is-automated", "false"):
			isAutomated = false
		default:
			return nil, invalidFilter{"is-automated", searchFilters.Get("is-automated")}
		}
	}
	if searchFilters.Contains("is-official") {
		switch {
		case searchFilters.UniqueExactMatch("is-official", "true"):
			isOfficial = true
		case searchFilters.UniqueExactMatch("is-official", "false"):
			isOfficial = false
		default:
			return nil, invalidFilter{"is-official", searchFilters.Get("is-official")}
		}
	}
	if searchFilters.Contains("stars") {
		hasStars := searchFilters.Get("stars")
		for _, hasStar := range hasStars {
			iHasStar, err := strconv.Atoi(hasStar)
			if err != nil {
				return nil, invalidFilter{"stars", hasStar}
			}
			if iHasStar > hasStarFilter {
				hasStarFilter = iHasStar
			}
		}
	}

	unfilteredResult, err := i.registryService.Search(ctx, term, limit, authConfig, dockerversion.DockerUserAgent(ctx), headers)
	if err != nil {
		return nil, err
	}

	filteredResults := unfilteredResult.Results[:0]
	for _, result := range unfilteredResult.Results {
		if searchFilters.Contains("is-automated") && isAutomated != result.IsAutomated {
			continue
		}
		if searchFilters.Contains("is-official") && isOfficial != result.IsOfficial {
			continue
		}
		if searchFilters.Contains("stars") && result.StarCount < hasStarFilter {
			continue
		}
		filteredResults = append(filteredResults, result)
	}

	return &registrytypes.SearchResults{
		Query:      unfilteredResult.Query,
		NumResults: len(filteredResults),
		Results:    filteredResults,
	}, nil
}
