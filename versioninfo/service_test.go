package versioninfo

import (
	"math"
	"strconv"
	"testing"
)

func TestParseSemverValid(t *testing.T) {
	tests := []struct {
		in   string
		want VersionInfo
	}{
		{"1.2.3", VersionInfo{Major: 1, Minor: 2, Patch: 3}},
		{"v1.2.3", VersionInfo{Prefix: "v", Major: 1, Minor: 2, Patch: 3}},
		{"v1.2.3-rc.1+abc", VersionInfo{Prefix: "v", Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1", Build: "abc"}},
		{strconv.Itoa(math.MaxInt-1) + ".0.0", VersionInfo{Major: math.MaxInt - 1}},
	}
	for _, tt := range tests {
		ok, got := ParseSemver(tt.in)
		if !ok {
			t.Errorf("ParseSemver(%q) rejected, want accepted", tt.in)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseSemver(%q) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestParseSemverRejectsOverflow(t *testing.T) {
	max := strconv.Itoa(math.MaxInt)
	for _, in := range []string{
		"99999999999999999999.0.0",
		"0.99999999999999999999.0",
		"0.0.99999999999999999999",
		max + ".0.0",
		"0." + max + ".0",
		"v0.0." + max,
	} {
		if ok, v := ParseSemver(in); ok {
			t.Errorf("ParseSemver(%q) accepted as %+v, want rejected", in, v)
		}
	}
}

func TestAcceptedVersionsBumpWithoutGoingNegative(t *testing.T) {
	_, v := ParseSemver(strconv.Itoa(math.MaxInt-1) + "." + strconv.Itoa(math.MaxInt-1) + "." + strconv.Itoa(math.MaxInt-1))
	v.Major++
	v.Minor++
	v.Patch++
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		t.Fatalf("bumped version went negative: %+v", v)
	}
	if s := v.String(); s[0] == '-' {
		t.Fatalf("bumped version string starts with '-': %q", s)
	}
}
