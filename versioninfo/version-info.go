// Package versioninfo contains the version info object that is serialised from the version tags
package versioninfo

import (
	"strconv"
	"strings"
)

type VersionInfo struct {
	Prefix     string
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

func (v VersionInfo) Compare(other VersionInfo) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	if v.PreRelease != other.PreRelease {
		if v.PreRelease == "" && other.PreRelease != "" {
			return 1
		} else if v.PreRelease != "" && other.PreRelease == "" {
			return -1
		} else {
			return comparePreRelease(v.PreRelease, other.PreRelease)
		}
	}

	return 0
}

func comparePreRelease(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		v1Num, v1Err := strconv.Atoi(v1Parts[i])
		v2Num, v2Err := strconv.Atoi(v2Parts[i])

		if v1Err == nil && v2Err == nil {
			if v1Num > v2Num {
				return 1
			} else if v1Num < v2Num {
				return -1
			}
		} else {
			if v1Parts[i] > v2Parts[i] {
				return 1
			} else if v1Parts[i] < v2Parts[i] {
				return -1
			}
		}
	}

	if len(v1Parts) > len(v2Parts) {
		return 1
	} else if len(v1Parts) < len(v2Parts) {
		return -1
	}

	return 0
}

func (v VersionInfo) String() string {
	var sb strings.Builder

	// Add prefix if it exists
	if v.Prefix != "" {
		sb.WriteString(v.Prefix)
	}

	// Add major, minor, and patch versions
	sb.WriteString(strconv.Itoa(v.Major))
	sb.WriteString(".")
	sb.WriteString(strconv.Itoa(v.Minor))
	sb.WriteString(".")
	sb.WriteString(strconv.Itoa(v.Patch))

	// Add pre-release if it exists
	if v.PreRelease != "" {
		sb.WriteString("-")
		sb.WriteString(v.PreRelease)
	}

	// Add build if it exists
	if v.Build != "" {
		sb.WriteString("+")
		sb.WriteString(v.Build)
	}

	return sb.String()
}
