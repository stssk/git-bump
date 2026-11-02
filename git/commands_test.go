package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stssk/git-bump/versioninfo"
)

// newRepo creates an empty git repository with one commit in a temp dir and
// makes it the working directory for the rest of the test.
func newRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	run(t, "init", "-q", "-b", "main")
	run(t, "config", "user.email", "test@example.com")
	run(t, "config", "user.name", "Test")
	run(t, "config", "commit.gpgsign", "false")
	run(t, "config", "tag.gpgsign", "false")
	run(t, "commit", "-q", "--allow-empty", "-m", "initial")
	return dir
}

func run(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestTagArgsSeparateNameFromOptions(t *testing.T) {
	// A version that would start with "-" must still come after "--".
	v := versioninfo.VersionInfo{Major: -1}
	for _, kind := range []TagKind{Lightweight, Annotated, Signed} {
		args := tagArgs(v, kind)
		sep := slices.Index(args, "--")
		if sep < 0 || sep != len(args)-2 || args[len(args)-1] != "-1.0.0" {
			t.Errorf("kind %d: tag name not placed right after \"--\": %q", kind, args)
		}
	}
}

func TestTagArgsKinds(t *testing.T) {
	v := versioninfo.VersionInfo{Prefix: "v", Major: 1}
	tests := map[TagKind][]string{
		Lightweight: {"tag", "--", "v1.0.0"},
		Annotated:   {"tag", "-a", "-m", "v1.0.0", "--", "v1.0.0"},
		Signed:      {"tag", "-s", "-m", "v1.0.0", "--", "v1.0.0"},
	}
	for kind, want := range tests {
		if got := tagArgs(v, kind); !slices.Equal(got, want) {
			t.Errorf("kind %d: got %q, want %q", kind, got, want)
		}
	}
}

func TestTagCreatesLightweightTag(t *testing.T) {
	newRepo(t)
	Tag(versioninfo.VersionInfo{Prefix: "v", Major: 1, Minor: 2, Patch: 3}, Lightweight)
	if got := run(t, "tag", "--list"); got != "v1.2.3" {
		t.Fatalf("tags = %q, want v1.2.3", got)
	}
	if got := run(t, "cat-file", "-t", "v1.2.3"); got != "commit" {
		t.Fatalf("object type = %q, want commit (lightweight)", got)
	}
}

func TestTagCreatesAnnotatedTag(t *testing.T) {
	newRepo(t)
	Tag(versioninfo.VersionInfo{Prefix: "v", Major: 2}, Annotated)
	if got := run(t, "cat-file", "-t", "v2.0.0"); got != "tag" {
		t.Fatalf("object type = %q, want tag (annotated)", got)
	}
	if got := run(t, "tag", "-l", "--format=%(contents:subject)", "v2.0.0"); got != "v2.0.0" {
		t.Fatalf("tag message = %q, want v2.0.0", got)
	}
}

func TestTagCreatesSignedTag(t *testing.T) {
	if _, err := exec.LookPath("ssh-keygen"); err != nil {
		t.Skip("ssh-keygen not available")
	}
	dir := newRepo(t)
	key := filepath.Join(t.TempDir(), "key")
	if out, err := exec.Command("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-C", "test", "-f", key).CombinedOutput(); err != nil {
		t.Fatalf("ssh-keygen: %v\n%s", err, out)
	}
	pub, err := os.ReadFile(key + ".pub")
	if err != nil {
		t.Fatal(err)
	}
	allowed := filepath.Join(dir, ".git", "allowed_signers")
	if err := os.WriteFile(allowed, []byte("test@example.com "+string(pub)), 0o600); err != nil {
		t.Fatal(err)
	}
	run(t, "config", "gpg.format", "ssh")
	run(t, "config", "user.signingkey", key)
	run(t, "config", "gpg.ssh.allowedSignersFile", allowed)

	Tag(versioninfo.VersionInfo{Prefix: "v", Major: 3}, Signed)
	run(t, "tag", "-v", "v3.0.0") // fails the test if the signature does not verify
}

func TestGetLastVersionIgnoresTagsOutsideHistory(t *testing.T) {
	newRepo(t)
	run(t, "tag", "v1.0.0")

	// A tag on a commit that is not part of HEAD's history, as when tags
	// are fetched from a fork.
	run(t, "checkout", "-q", "-b", "fork")
	run(t, "commit", "-q", "--allow-empty", "-m", "fork commit")
	run(t, "tag", "v99.0.0")
	run(t, "checkout", "-q", "main")

	got, found := GetLastVersion("", false)
	if !found || got.String() != "v1.0.0" {
		t.Fatalf("GetLastVersion() = %s, %v; want v1.0.0, true", got.String(), found)
	}
}

func TestGetLastVersionWithPrefix(t *testing.T) {
	newRepo(t)
	run(t, "tag", "v1.4.2")
	run(t, "tag", "x99.0.0")

	got, found := GetLastVersion("v", true)
	if !found || got.String() != "v1.4.2" {
		t.Fatalf("GetLastVersion(\"v\") = %s, %v; want v1.4.2, true", got.String(), found)
	}
}
