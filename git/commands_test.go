package git

import (
	"os/exec"
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
	args := tagArgs(versioninfo.VersionInfo{Major: -1})
	sep := -1
	for i, a := range args {
		if a == "--" {
			sep = i
		}
	}
	if sep < 0 || args[len(args)-1] != "-1.0.0" || sep > len(args)-2 {
		t.Fatalf("tag name not placed after \"--\": %q", args)
	}
}

func TestTagCreatesTag(t *testing.T) {
	newRepo(t)
	Tag(versioninfo.VersionInfo{Prefix: "v", Major: 1, Minor: 2, Patch: 3})
	if got := run(t, "tag", "--list"); got != "v1.2.3" {
		t.Fatalf("tags = %q, want v1.2.3", got)
	}
}
