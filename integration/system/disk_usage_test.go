package system // import "github.com/docker/docker/integration/system"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/integration/internal/container"
	"github.com/docker/docker/testutil/request"
	"gotest.tools/v3/assert"
)

func TestDiskUsage(t *testing.T) {
	defer setupTest(t)()

	ctx := context.Background()
	client := testEnv.APIClient()

	assertGet := func(t *testing.T, expected types.DiskUsage, params ...string) {
		t.Helper()

		ep := "/system/df"
		if len(params) > 0 {
			ep = fmt.Sprintf("%s?%s", ep, strings.Join(params, "&"))
		} else {
			// When no query parameters are specified, the return values should be identical.
			du, err := client.DiskUsage(ctx)
			assert.NilError(t, err)
			assert.DeepEqual(t, du, expected)
		}
		res, body, err := request.Get(ep, request.JSON)
		assert.NilError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		buf, err := request.ReadBody(body)
		assert.NilError(t, err)
		var du types.DiskUsage
		err = json.Unmarshal(buf, &du)
		assert.NilError(t, err)
		assert.DeepEqual(t, du, expected)

	}

	baseDU, err := client.DiskUsage(ctx)
	assert.NilError(t, err)
	assert.Assert(t, baseDU.LayersSize > 0)
	assert.Assert(t, len(baseDU.Images) > 0)
	assert.DeepEqual(t, baseDU, types.DiskUsage{
		LayersSize: baseDU.LayersSize,
		Images:     baseDU.Images,
		Containers: []*types.Container{},
		Volumes:    []*types.Volume{},
		BuildCache: []*types.BuildCache{},
	})

	assertGet(t, baseDU)
	assertGet(t, baseDU,
		"types=build-cache",
		"types=container",
		"types=image",
		"types=volume",
	)
	assertGet(t, types.DiskUsage{
		LayersSize: baseDU.LayersSize,
		Images:     baseDU.Images,
	},
		"types=image",
	)
	assertGet(t, types.DiskUsage{
		BuildCache: []*types.BuildCache{},
	},
		"types=build-cache",
	)
	assertGet(t, types.DiskUsage{
		Containers: []*types.Container{},
		Volumes:    []*types.Volume{},
	},
		"types=container",
		"types=volume",
	)

	cID := container.Run(ctx, t, client)

	du, err := client.DiskUsage(ctx)
	assert.NilError(t, err)
	assert.Assert(t, len(du.Containers) == 1 && du.Containers[0].ID == cID)
	assert.DeepEqual(t, du, types.DiskUsage{
		LayersSize: baseDU.LayersSize,
		Images:     du.Images, // Container counter for one of the images should have increased.
		Containers: du.Containers,
		Volumes:    []*types.Volume{},
		BuildCache: []*types.BuildCache{},
	})

	assertGet(t, du)
	assertGet(t, du,
		"types=build-cache",
		"types=container",
		"types=image",
		"types=volume",
	)
	assertGet(t, types.DiskUsage{
		LayersSize: du.LayersSize,
		Images:     du.Images,
	},
		"types=image",
	)
	assertGet(t, types.DiskUsage{
		LayersSize: du.LayersSize,
		Images:     du.Images,
		Volumes:    []*types.Volume{},
	},
		"types=image",
		"types=volume",
	)
	assertGet(t, types.DiskUsage{
		Containers: du.Containers,
	},
		"types=container",
	)
}
