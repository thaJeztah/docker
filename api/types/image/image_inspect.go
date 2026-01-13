package image

import (
	"time"

	dockerspec "github.com/moby/docker-image-spec/specs-go/v1"
	"github.com/moby/moby/api/types/storage"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// RootFS returns Image's RootFS description including the layer IDs.
type RootFS struct {
	Type   string   `json:",omitempty"`
	Layers []string `json:",omitempty"`
}

// InspectResponse contains response of Engine API:
// GET "/images/{name:.*}/json"
type InspectResponse struct {
	// ID is the content-addressable ID of an image.
	//
	// This identifier is a content-addressable digest calculated from the
	// image's configuration (which includes the digests of layers used by
	// the image).
	//
	// Note that this digest differs from the `RepoDigests` below, which
	// holds digests of image manifests that reference the image.
	ID string `json:"Id"`

	// RepoTags is a list of image names/tags in the local image cache that
	// reference this image.
	//
	// Multiple image tags can refer to the same image, and this list may be
	// empty if no tags reference the image, in which case the image is
	// "untagged", in which case it can still be referenced by its ID.
	RepoTags []string

	// RepoDigests is a list of content-addressable digests of locally available
	// image manifests that the image is referenced from. Multiple manifests can
	// refer to the same image.
	//
	// These digests are usually only available if the image was either pulled
	// from a registry, or if the image was pushed to a registry, which is when
	// the manifest is generated and its digest calculated.
	RepoDigests []string

	// Comment is an optional message that can be set when committing or
	// importing the image. This field is omitted if not set.
	Comment string `json:",omitempty"`

	// Created is the date and time at which the image was created, formatted in
	// RFC 3339 nano-seconds (time.RFC3339Nano).
	//
	// This information is only available if present in the image,
	// and omitted otherwise.
	Created string `json:",omitempty"`

	// Author is the name of the author that was specified when committing the
	// image, or as specified through MAINTAINER (deprecated) in the Dockerfile.
	// This field is omitted if not set.
	Author string `json:",omitempty"`
	Config *dockerspec.DockerOCIImageConfig

	// Architecture is the hardware CPU architecture that the image runs on.
	Architecture string

	// Variant is the CPU architecture variant (presently ARM-only).
	Variant string `json:",omitempty"`

	// OS is the Operating System the image is built to run on.
	Os string

	// OsVersion is the version of the Operating System the image is built to
	// run on (especially for Windows).
	OsVersion string `json:",omitempty"`

	// Size is the total size of the image including all layers it is composed of.
	Size int64

	// GraphDriver holds information about the storage driver used to store the
	// container's and image's filesystem.
	GraphDriver *storage.DriverData `json:"GraphDriver,omitempty"`

	// RootFS contains information about the image's RootFS, including the
	// layer IDs.
	RootFS RootFS

	// Metadata of the image in the local cache.
	//
	// This information is local to the daemon, and not part of the image itself.
	Metadata Metadata

	// Descriptor is the OCI descriptor of the image target.
	// It's only set if the daemon provides a multi-platform image store.
	//
	// WARNING: This is experimental and may change at any time without any backward
	// compatibility.
	Descriptor *ocispec.Descriptor `json:"Descriptor,omitempty"`

	// Manifests is a list of image manifests available in this image. It
	// provides a more detailed view of the platform-specific image manifests or
	// other image-attached data like build attestations.
	//
	// Only available if the daemon provides a multi-platform image store, the client
	// requests manifests AND does not request a specific platform.
	//
	// WARNING: This is experimental and may change at any time without any backward
	// compatibility.
	Manifests []ManifestSummary `json:"Manifests,omitempty"`

	// Identity holds information about the identity and origin of the image.
	// This is trusted information verified by the daemon and cannot be modified
	// by tagging an image to a different name.
	Identity *Identity `json:"Identity,omitempty"`
}

// Identity holds information about the identity and origin of the image.
// This is trusted information verified by the daemon and cannot be modified
// by tagging an image to a different name.
type Identity struct {
	// Signature contains the properties of verified signatures for the image.
	Signature []SignatureIdentity `json:"Signature,omitzero"`
	// Pull contains remote location information if image was created via pull.
	// If image was pulled via mirror, this contains the original repository location.
	Pull []PullIdentity `json:"Pull,omitzero"`
	// Build contains build reference information if image was created via build.
	Build []BuildIdentity `json:"Build,omitzero"`
}

// BuildIdentity contains build reference information if image was created via build.
type BuildIdentity struct {
	// Ref is the identifier for the build request. This reference can be used to
	// look up the build details in BuildKit history API.
	Ref string `json:"Ref,omitempty"`

	// CreatedAt is the time when the build ran.
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
}

// PullIdentity contains remote location information if image was created via pull.
// If image was pulled via mirror, this contains the original repository location.
type PullIdentity struct {
	// Repository is the remote repository location the image was pulled from.
	Repository string `json:"Repository,omitempty"`
}

// SignatureIdentity contains the properties of verified signatures for the image.
type SignatureIdentity struct {
	// Name is a textual description summarizing the type of signature.
	Name string `json:"Name,omitempty"`
	// Timestamps contains a list of verified signed timestamps for the signature.
	Timestamps []SignatureTimestamp `json:"Timestamps,omitzero"`
	// KnownSigner is an identifier for a special signer identity that is known to the implementation.
	KnownSigner KnownSignerIdentity `json:"KnownSigner,omitempty"`
	// DockerReference is the Docker image reference associated with the signature.
	// This is an optional field only present in older hashedrecord signatures.
	DockerReference string `json:"DockerReference,omitempty"`
	// Signer contains information about the signer certificate used to sign the image.
	Signer *SignerIdentity `json:"Signer,omitempty"`
	// SignatureType is the type of signature format. E.g. "bundle-v0.3" or "hashedrecord".
	SignatureType string `json:"SignatureType,omitempty"`

	// Error contains error information if signature verification failed.
	// Other fields will be empty in this case.
	Error string `json:"Error,omitempty"`
	// Warnings contains any warnings that occurred during signature verification.
	// For example, if there was no internet connectivity and cached trust roots were used.
	// Warning does not indicate a failed verification but may point to configuration issues.
	Warnings []string `json:"Warnings,omitzero"`
}

// SignatureTimestamp contains information about a verified signed timestamp for an image signature.
type SignatureTimestamp struct {
	Type      SignatureTimestampType `json:"Type"`
	URI       string                 `json:"URI"`
	Timestamp time.Time              `json:"Timestamp"`
}

// SignatureTimestampType is the type of timestamp used in the signature.
type SignatureTimestampType string

const (
	SignatureTimestampTlog      SignatureTimestampType = "Tlog"
	SignatureTimestampAuthority SignatureTimestampType = "TimestampAuthority"
)

// SignatureType is the type of signature format.
type SignatureType string

const (
	SignatureTypeBundleV03       SignatureType = "bundle-v0.3"
	SignatureTypeSimpleSigningV1 SignatureType = "simplesigning-v1"
)

// KnownSignerIdentity is an identifier for a special signer identity that is known to the implementation.
type KnownSignerIdentity string

const (
	// KnownSignerDHI is the known identity for Docker Hardened Images.
	KnownSignerDHI KnownSignerIdentity = "DHI"
)

// SignerIdentity contains information about the signer certificate used to sign the image.
// This is certificate.Summary with deprecated fields removed and keys in Moby uppercase style.
type SignerIdentity struct {
	CertificateIssuer      string `json:"CertificateIssuer"`
	SubjectAlternativeName string `json:"SubjectAlternativeName"`
	// The OIDC issuer. Should match `iss` claim of ID token or, in the case of
	// a federated login like Dex it should match the issuer URL of the
	// upstream issuer. The issuer is not set the extensions are invalid and
	// will fail to render.
	Issuer string `json:"Issuer,omitempty"` // OID 1.3.6.1.4.1.57264.1.8 and 1.3.6.1.4.1.57264.1.1 (Deprecated)

	// Reference to specific build instructions that are responsible for signing.
	BuildSignerURI string `json:"buildSignerURI,omitempty"` //nolint:tagliatelle // 1.3.6.1.4.1.57264.1.9

	// Immutable reference to the specific version of the build instructions that is responsible for signing.
	BuildSignerDigest string `json:"buildSignerDigest,omitempty"` // 1.3.6.1.4.1.57264.1.10

	// Specifies whether the build took place in platform-hosted cloud infrastructure or customer/self-hosted infrastructure.
	RunnerEnvironment string `json:"runnerEnvironment,omitempty"` // 1.3.6.1.4.1.57264.1.11

	// Source repository URL that the build was based on.
	SourceRepositoryURI string `json:"sourceRepositoryURI,omitempty"` //nolint:tagliatelle  // 1.3.6.1.4.1.57264.1.12

	// Immutable reference to a specific version of the source code that the build was based upon.
	SourceRepositoryDigest string `json:"sourceRepositoryDigest,omitempty"` // 1.3.6.1.4.1.57264.1.13

	// Source Repository Ref that the build run was based upon.
	SourceRepositoryRef string `json:"sourceRepositoryRef,omitempty"` // 1.3.6.1.4.1.57264.1.14

	// Immutable identifier for the source repository the workflow was based upon.
	SourceRepositoryIdentifier string `json:"sourceRepositoryIdentifier,omitempty"` // 1.3.6.1.4.1.57264.1.15

	// Source repository owner URL of the owner of the source repository that the build was based on.
	SourceRepositoryOwnerURI string `json:"sourceRepositoryOwnerURI,omitempty"` //nolint:tagliatelle // 1.3.6.1.4.1.57264.1.16

	// Immutable identifier for the owner of the source repository that the workflow was based upon.
	SourceRepositoryOwnerIdentifier string `json:"sourceRepositoryOwnerIdentifier,omitempty"` // 1.3.6.1.4.1.57264.1.17

	// Build Config URL to the top-level/initiating build instructions.
	BuildConfigURI string `json:"buildConfigURI,omitempty"` //nolint:tagliatelle // 1.3.6.1.4.1.57264.1.18

	// Immutable reference to the specific version of the top-level/initiating build instructions.
	BuildConfigDigest string `json:"buildConfigDigest,omitempty"` // 1.3.6.1.4.1.57264.1.19

	// Event or action that initiated the build.
	BuildTrigger string `json:"buildTrigger,omitempty"` // 1.3.6.1.4.1.57264.1.20

	// Run Invocation URL to uniquely identify the build execution.
	RunInvocationURI string `json:"runInvocationURI,omitempty"` //nolint:tagliatelle // 1.3.6.1.4.1.57264.1.21

	// Source repository visibility at the time of signing the certificate.
	SourceRepositoryVisibilityAtSigning string `json:"sourceRepositoryVisibilityAtSigning,omitempty"` // 1.3.6.1.4.1.57264.1.22
}
