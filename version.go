package mcp

// Version information
const (
	// Version is the current version of the gorest-mcp plugin
	Version = "0.1.0"

	// VersionMajor is the major version number
	VersionMajor = 0

	// VersionMinor is the minor version number
	VersionMinor = 1

	// VersionPatch is the patch version number
	VersionPatch = 0
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}
