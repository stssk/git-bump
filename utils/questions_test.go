package utils

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/stssk/git-bump/versioninfo"
)

func TestValidPreRelease(t *testing.T) {
	valid := []string{"rc", "rc.1", "alpha-2", "0", "beta.11.x-y", "RC1"}
	invalid := []string{
		"", " ", "rc 1", "rc..1", ".rc", "rc.", "x:refs/heads/main",
		"a/../b", "rc+build", "rc_1", "rc\t1", "æøå", "rc\n1", "-rc;rm -rf /",
	}
	for _, s := range valid {
		if !ValidPreRelease(s) {
			t.Errorf("ValidPreRelease(%q) = false, want true", s)
		}
	}
	for _, s := range invalid {
		if ValidPreRelease(s) {
			t.Errorf("ValidPreRelease(%q) = true, want false", s)
		}
	}
}

func TestValidPreReleaseRoundTrips(t *testing.T) {
	for _, s := range []string{"rc", "rc.1", "alpha-2", "beta.11.x-y"} {
		v := versioninfo.VersionInfo{Prefix: "v", Major: 1, PreRelease: s}
		ok, parsed := versioninfo.ParseSemver(v.String())
		if !ok || parsed.PreRelease != s {
			t.Errorf("%q: ParseSemver(%q) = %v, %+v", s, v.String(), ok, parsed)
		}
	}
}

func TestReadPreReleaseRepromptsUntilValid(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("rc 1\nx:refs/heads/main\n\nrc.1\n"))
	var out strings.Builder
	got, err := readPreRelease(in, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "rc.1" {
		t.Fatalf("got %q, want rc.1", got)
	}
	if n := strings.Count(out.String(), "Not a valid pre release version"); n != 3 {
		t.Fatalf("rejections = %d, want 3; output:\n%s", n, out.String())
	}
}

func TestReadPreReleaseAcceptsLastLineWithoutNewline(t *testing.T) {
	got, err := readPreRelease(bufio.NewReader(strings.NewReader("rc.2")), io.Discard)
	if err != nil || got != "rc.2" {
		t.Fatalf("got %q, %v; want rc.2, nil", got, err)
	}
}

func TestReadPreReleaseStopsAtEOF(t *testing.T) {
	_, err := readPreRelease(bufio.NewReader(strings.NewReader("bad input\n")), io.Discard)
	if err != io.EOF {
		t.Fatalf("err = %v, want io.EOF", err)
	}
}
