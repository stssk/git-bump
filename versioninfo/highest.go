package versioninfo

import (
	"fmt"
	"sort"
	"strings"
)

// Highest returns the highest semantic version among tags.
//
// If prefixSet is true, only tags with exactly that prefix are considered.
// Otherwise all version tags must share a single prefix; if they don't, an
// error is returned rather than guessing, since a stray tag with a different
// prefix (e.g. "x99.0.0" next to "v1.4.2") would otherwise decide the next
// version.
func Highest(tags []string, prefix string, prefixSet bool) (VersionInfo, bool, error) {
	var highest *VersionInfo
	prefixes := map[string]bool{}
	for _, tag := range tags {
		valid, info := ParseSemver(strings.TrimSpace(tag))
		if !valid {
			continue
		}
		if prefixSet && info.Prefix != prefix {
			continue
		}
		prefixes[info.Prefix] = true
		if highest == nil || info.Compare(*highest) > 0 {
			highest = &info
		}
	}
	if len(prefixes) > 1 {
		found := make([]string, 0, len(prefixes))
		for p := range prefixes {
			found = append(found, fmt.Sprintf("%q", p))
		}
		sort.Strings(found)
		return VersionInfo{}, false, fmt.Errorf("version tags use different prefixes (%s); choose one with -prefix", strings.Join(found, ", "))
	}
	if highest == nil {
		return VersionInfo{}, false, nil
	}
	return *highest, true, nil
}
