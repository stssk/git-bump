package versioninfo

import (
	"fmt"
	"sort"
	"strings"
)

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
