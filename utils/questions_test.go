package utils

import (
	"errors"
	"slices"
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

func TestAskForPreReleaseVersionValidates(t *testing.T) {
	p := &fakePrompter{inputs: []string{"rc 1", "x:refs/heads/main", "", " rc.1 "}}
	got, err := AskForPreReleaseVersion(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "rc.1" {
		t.Fatalf("got %q, want rc.1", got)
	}
	if want := []string{"rc 1", "x:refs/heads/main", ""}; !slices.Equal(p.rejected, want) {
		t.Fatalf("rejected %q, want %q", p.rejected, want)
	}
}

func TestAskForPreReleaseVersionReturnsAbort(t *testing.T) {
	_, err := AskForPreReleaseVersion(&fakePrompter{err: ErrAborted})
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("err = %v, want ErrAborted", err)
	}
}
