package versioninfo

import (
	"regexp"
	"strconv"
)

func ParseSemver(version string) (bool, VersionInfo) {
	// Regular expression to match a semantic version with optional prefix, pre-release, and build metadata
	semverRegex := `^([a-zA-Z]*)?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`
	re := regexp.MustCompile(semverRegex)
	matches := re.FindStringSubmatch(version)

	if len(matches) == 0 {
		return false, VersionInfo{}
	}

	// Extracting components
	prefix := matches[1]
	major, _ := strconv.Atoi(matches[2])
	minor, _ := strconv.Atoi(matches[3])
	patch, _ := strconv.Atoi(matches[4])
	preRelease := matches[5]
	build := matches[6]

	return true, VersionInfo{
		Prefix:     prefix,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Build:      build,
	}
}
