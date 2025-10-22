//go:build !windows

package authz

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/moby/moby/client"
	"github.com/moby/moby/v2/integration/internal/container"
	"github.com/moby/moby/v2/integration/internal/requirement"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/skip"
)

var (
	authzPluginName            = "riyaz/authz-no-volume-plugin"
	authzPluginTag             = "latest"
	authzPluginNameWithTag     = authzPluginName + ":" + authzPluginTag
	authzPluginBadManifestName = "riyaz/authz-plugin-bad-manifest"
	nonexistentAuthzPluginName = "riyaz/nonexistent-authz-plugin"
)

func setupTestV2(t *testing.T) context.Context {
	skip.If(t, testEnv.DaemonInfo.OSType == "windows")
	skip.If(t, !requirement.HasHubConnectivity(t))

	ctx := setupTest(t)
	d.Start(t)
	return ctx
}

func TestAuthZPluginV2AllowNonVolumeRequest(t *testing.T) {
	skip.If(t, testEnv.NotAmd64)
	ctx := setupTestV2(t)

	c := d.NewClientT(t)

	// Install authz plugin
	err := pluginInstallGrantAllPermissions(ctx, c, authzPluginNameWithTag)
	assert.NilError(t, err)
	// start the daemon with the plugin and load busybox, --net=none build fails otherwise
	// because it needs to pull busybox
	d.Restart(t, "--authorization-plugin="+authzPluginNameWithTag)
	d.LoadBusybox(ctx, t)

	// Ensure docker run command and accompanying docker ps are successful
	cID := container.Run(ctx, t, c)

	_, err = c.ContainerInspect(ctx, cID)
	assert.NilError(t, err)
}

func TestAuthZPluginV2Disable(t *testing.T) {
	skip.If(t, testEnv.NotAmd64)
	ctx := setupTestV2(t)

	c := d.NewClientT(t)

	// Install authz plugin
	err := pluginInstallGrantAllPermissions(ctx, c, authzPluginNameWithTag)
	assert.NilError(t, err)

	d.Restart(t, "--authorization-plugin="+authzPluginNameWithTag)
	d.LoadBusybox(ctx, t)

	_, err = c.VolumeCreate(ctx, client.VolumeCreateOptions{Driver: "local"})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))

	// disable the plugin
	_, err = c.PluginDisable(ctx, authzPluginNameWithTag, client.PluginDisableOptions{})
	assert.NilError(t, err)

	// now test to see if the docker api works.
	_, err = c.VolumeCreate(ctx, client.VolumeCreateOptions{Driver: "local"})
	assert.NilError(t, err)
}

func TestAuthZPluginV2RejectVolumeRequests(t *testing.T) {
	skip.If(t, testEnv.NotAmd64)
	ctx := setupTestV2(t)

	c := d.NewClientT(t)

	// Install authz plugin
	err := pluginInstallGrantAllPermissions(ctx, c, authzPluginNameWithTag)
	assert.NilError(t, err)

	// restart the daemon with the plugin
	d.Restart(t, "--authorization-plugin="+authzPluginNameWithTag)

	_, err = c.VolumeCreate(ctx, client.VolumeCreateOptions{Driver: "local"})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))

	_, err = c.VolumeList(ctx, client.VolumeListOptions{})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))

	// The plugin will block the command before it can determine the volume does not exist
	err = c.VolumeRemove(ctx, "test", client.VolumeRemoveOptions{})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))

	_, err = c.VolumeInspect(ctx, "test", client.VolumeInspectOptions{})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))

	_, err = c.VolumesPrune(ctx, client.VolumePruneOptions{})
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, fmt.Sprintf("Error response from daemon: plugin %s failed with error:", authzPluginNameWithTag))
}

func TestAuthZPluginV2BadManifestFailsDaemonStart(t *testing.T) {
	skip.If(t, testEnv.NotAmd64)
	ctx := setupTestV2(t)

	c := d.NewClientT(t)

	// Install authz plugin with bad manifest
	err := pluginInstallGrantAllPermissions(ctx, c, authzPluginBadManifestName)
	assert.NilError(t, err)

	// start the daemon with the plugin, it will error
	err = d.RestartWithError("--authorization-plugin=" + authzPluginBadManifestName)
	assert.Assert(t, err != nil)

	// restarting the daemon without requiring the plugin will succeed
	d.Start(t)
}

func TestAuthZPluginV2NonexistentFailsDaemonStart(t *testing.T) {
	_ = setupTestV2(t)

	// start the daemon with a non-existent authz plugin, it will error
	err := d.RestartWithError("--authorization-plugin=" + nonexistentAuthzPluginName)
	assert.Assert(t, err != nil)

	// restarting the daemon without requiring the plugin will succeed
	d.Start(t)
}

func pluginInstallGrantAllPermissions(ctx context.Context, apiClient client.APIClient, name string) error {
	responseReader, err := apiClient.PluginInstall(ctx, "", client.PluginInstallOptions{
		RemoteRef:            name,
		AcceptAllPermissions: true,
	})
	if err != nil {
		return err
	}
	defer responseReader.Close()
	// we have to read the response out here because the client API
	// actually starts a goroutine which we can only be sure has
	// completed when we get EOF from reading responseBody
	_, err = io.ReadAll(responseReader)
	return err
}
