package image

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
