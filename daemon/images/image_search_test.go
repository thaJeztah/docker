package images // import "github.com/docker/docker/daemon/images"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

type FakeService struct {
	registry.DefaultService

	shouldReturnError bool

	term    string
	limit   int
	results []registrytypes.SearchResult
}

func (s *FakeService) Search(ctx context.Context, term string, limit int, authConfig *types.AuthConfig, userAgent string, headers map[string][]string) (*registrytypes.SearchResults, error) {
	if s.shouldReturnError {
		return nil, errors.New("Search unknown error")
	}
	if s.term != "" && term != s.term {
		return nil, fmt.Errorf("incorrect term; expected %q, got %q", s.term, term)
	}
	if s.limit != 0 && limit != s.limit {
		return nil, fmt.Errorf("incorrect limit; expected %d, got %d", s.limit, limit)
	}
	return &registrytypes.SearchResults{
		Query:      s.term,
		NumResults: len(s.results),
		Results:    s.results,
	}, nil
}

func TestSearchRegistryForImagesErrors(t *testing.T) {
	errorCases := []struct {
		filtersArgs       string
		shouldReturnError bool
		expectedError     string
	}{
		{
			expectedError:     "Search unknown error",
			shouldReturnError: true,
		},
		{
			filtersArgs:   "invalid json",
			expectedError: "invalid character 'i' looking for beginning of value",
		},
		{
			filtersArgs:   `{"type":{"custom":true}}`,
			expectedError: "Invalid filter 'type'",
		},
		{
			filtersArgs:   `{"is-automated":{"invalid":true}}`,
			expectedError: "Invalid filter 'is-automated=[invalid]'",
		},
		{
			filtersArgs:   `{"is-automated":{"true":true,"false":true}}`,
			expectedError: "Invalid filter 'is-automated",
		},
		{
			filtersArgs:   `{"is-official":{"invalid":true}}`,
			expectedError: "Invalid filter 'is-official=[invalid]'",
		},
		{
			filtersArgs:   `{"is-official":{"true":true,"false":true}}`,
			expectedError: "Invalid filter 'is-official",
		},
		{
			filtersArgs:   `{"stars":{"invalid":true}}`,
			expectedError: "Invalid filter 'stars=invalid'",
		},
		{
			filtersArgs:   `{"stars":{"1":true,"invalid":true}}`,
			expectedError: "Invalid filter 'stars=invalid'",
		},
	}
	for index, e := range errorCases {
		daemon := &ImageService{
			registryService: &FakeService{
				shouldReturnError: e.shouldReturnError,
			},
		}
		_, err := daemon.SearchRegistryForImages(context.Background(), e.filtersArgs, "term", 25, nil, map[string][]string{})
		if err == nil {
			t.Errorf("%d: expected an error, got nothing", index)
		}
		if !strings.Contains(err.Error(), e.expectedError) {
			t.Errorf("%d: expected error to contain %s, got %s", index, e.expectedError, err.Error())
		}
	}
}

func TestSearchRegistryForImages(t *testing.T) {
	term := "term"
	successCases := []struct {
		filtersArgs     string
		registryResults []registrytypes.SearchResult
		expectedResults []registrytypes.SearchResult
	}{
		{
			filtersArgs:     "",
			registryResults: []registrytypes.SearchResult{},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: "",
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
		},
		{
			filtersArgs: `{"is-automated":{"true":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: `{"is-automated":{"true":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
		},
		{
			filtersArgs: `{"is-automated":{"false":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: `{"is-automated":{"false":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsAutomated: false,
				},
			},
		},
		{
			filtersArgs: `{"is-official":{"true":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: `{"is-official":{"true":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
		},
		{
			filtersArgs: `{"is-official":{"false":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  true,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: `{"is-official":{"false":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					IsOfficial:  false,
				},
			},
		},
		{
			filtersArgs: `{"stars":{"0":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
		},
		{
			filtersArgs: `{"stars":{"1":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name",
					Description: "description",
					StarCount:   0,
				},
			},
			expectedResults: []registrytypes.SearchResult{},
		},
		{
			filtersArgs: `{"stars":{"1":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
				},
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
				},
			},
		},
		{
			filtersArgs: `{"stars":{"1":true}, "is-official":{"true":true}, "is-automated":{"true":true}}`,
			registryResults: []registrytypes.SearchResult{
				{
					Name:        "name0",
					Description: "description0",
					StarCount:   0,
					IsOfficial:  true,
					IsAutomated: true,
				},
				{
					Name:        "name1",
					Description: "description1",
					StarCount:   1,
					IsOfficial:  false,
					IsAutomated: true,
				},
				{
					Name:        "name2",
					Description: "description2",
					StarCount:   1,
					IsOfficial:  true,
					IsAutomated: false,
				},
				{
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				},
			},
			expectedResults: []registrytypes.SearchResult{
				{
					Name:        "name3",
					Description: "description3",
					StarCount:   2,
					IsOfficial:  true,
					IsAutomated: true,
				},
			},
		},
	}
	for i, tc := range successCases {
		tc := tc
		t.Run(tc.filtersArgs, func(t *testing.T) {
			daemon := &ImageService{
				registryService: &FakeService{
					term:    term,
					limit:   25,
					results: tc.registryResults,
				},
			}
			results, err := daemon.SearchRegistryForImages(context.Background(), tc.filtersArgs, term, 25, nil, map[string][]string{})
			if err != nil {
				t.Fatalf("%d: %v", i, err)
			}
			if results.Query != term {
				t.Errorf("%d: expected Query to be %s, got %s", i, term, results.Query)
			}
			if results.NumResults != len(tc.expectedResults) {
				t.Errorf("%d: expected NumResults to be %d, got %d", i, len(tc.expectedResults), results.NumResults)
			}
			for _, result := range results.Results {
				found := false
				for _, expectedResult := range tc.expectedResults {
					if expectedResult.Name == result.Name &&
						expectedResult.Description == result.Description &&
						expectedResult.IsAutomated == result.IsAutomated &&
						expectedResult.IsOfficial == result.IsOfficial &&
						expectedResult.StarCount == result.StarCount {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%d: expected results %v, got %v", i, tc.expectedResults, results.Results)
				}
			}
		})
	}
}
