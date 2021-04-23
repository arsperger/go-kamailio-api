package utils

import "github.com/google/uuid"

// GenerateUUID is exported
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

var (
	commit  string
	version string = "0.0.1"
)

// Version returns strng with current version-commit
func Version() string {

	version = "v" + version

	if commit != "" {
		version += "-" + commit[0:12]
	}

	return version
}
