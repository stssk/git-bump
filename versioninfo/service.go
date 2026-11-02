package versioninfo

import (
	"math"
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
	major, okMajor := parseComponent(matches[2])
	minor, okMinor := parseComponent(matches[3])
	patch, okPatch := parseComponent(matches[4])
	if !okMajor || !okMinor || !okPatch {
		return false, VersionInfo{}
	}
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

// MaxInt is rejected too, so every accepted component can be bumped without overflowing.
func parseComponent(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n == math.MaxInt {
		return 0, false
	}
	return n, true
}
