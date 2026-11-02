package main

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stssk/git-bump/utils"
	"github.com/stssk/git-bump/versioninfo"
)

func TestNextVersion(t *testing.T) {
	v := func(major, minor, patch int, pre string) versioninfo.VersionInfo {
		return versioninfo.VersionInfo{Prefix: "v", Major: major, Minor: minor, Patch: patch, PreRelease: pre}
	}
	tests := []struct {
		name  string
		from  versioninfo.VersionInfo
		pre   bool
		input string
		want  versioninfo.VersionInfo
	}{
		{"patch", v(1, 2, 3, ""), false, "1\n", v(1, 2, 4, "")},
		{"minor", v(1, 2, 3, ""), false, "2\n", v(1, 3, 0, "")},
		{"major", v(1, 2, 3, ""), false, "3\n", v(2, 0, 0, "")},
		{"patch with -pre", v(1, 2, 3, ""), true, "2\n", v(1, 2, 4, "")},
		{"first pre-release", v(1, 2, 3, ""), true, "1\nrc.1\n", v(1, 2, 4, "rc.1")},
		{"next pre-release", v(1, 2, 4, "rc.1"), true, "1\nrc.2\n", v(1, 2, 4, "rc.2")},
		{"release from pre-release", v(1, 2, 4, "rc.2"), true, "2\n", v(1, 2, 5, "")},
		{"invalid pre-release reprompts", v(1, 2, 3, ""), true, "1\nrc 1\nrc1\n", v(1, 2, 4, "rc1")},
	}
	for _, tt := range tests {
		p := utils.NewLinePrompter(strings.NewReader(tt.input), io.Discard)
		got, err := nextVersion(p, tt.from, tt.pre)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, got.String(), tt.want.String())
		}
	}
}

func TestNextVersionStopsWhenInputEnds(t *testing.T) {
	from := versioninfo.VersionInfo{Prefix: "v", Major: 1}
	for _, input := range []string{"", "1\n"} {
		p := utils.NewLinePrompter(strings.NewReader(input), io.Discard)
		if _, err := nextVersion(p, from, true); !errors.Is(err, utils.ErrNoInput) {
			t.Errorf("input %q: err = %v, want ErrNoInput", input, err)
		}
	}
}
