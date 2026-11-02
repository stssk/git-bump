// Package git contains any git command used by the program and execution of the git commands
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stssk/git-bump/versioninfo"
)

func GitInstalled() string {
	cmd := exec.Command("git", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not find git. Please ensure git is in PATH")
		os.Exit(1)
	}
	return string(output)
}

func HeadTagged() string {
	cmd := exec.Command("git", "tag", "--points-at", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not check tags")
		os.Exit(5)
	}
	return strings.TrimSpace(string(output))
}

func CurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not find git. Please ensure git is in PATH")
		os.Exit(1)
	}
	return strings.TrimSpace(string(output))
}

// GetLastVersion returns the highest version tag reachable from HEAD. Tags on
// commits outside the current history (e.g. fetched from a fork) are ignored.
// See versioninfo.Highest for how prefix and prefixSet are used.
func GetLastVersion(prefix string, prefixSet bool) (versioninfo.VersionInfo, bool) {
	cmd := exec.Command("git", "tag", "--list", "--merged", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not fetch version tags")
		os.Exit(2)
	}
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	version, found, err := versioninfo.Highest(tags, prefix, prefixSet)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return version, found
}

func GetSha() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not fetch version tags")
		os.Exit(2)
	}
	return strings.TrimSpace(string(output))
}

// TagKind selects what kind of tag Tag creates.
type TagKind int

const (
	Lightweight TagKind = iota
	Annotated
	Signed
)

// tagArgs builds the arguments for creating a tag. The "--" ensures the tag
// name is never interpreted as an option, even if it starts with "-".
func tagArgs(version versioninfo.VersionInfo, kind TagKind) []string {
	name := version.String()
	args := []string{"tag"}
	switch kind {
	case Annotated:
		args = append(args, "-a", "-m", name)
	case Signed:
		args = append(args, "-s", "-m", name)
	}
	return append(args, "--", name)
}

func Tag(version versioninfo.VersionInfo, kind TagKind) {
	cmd := exec.Command("git", tagArgs(version, kind)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Could not tag version: %s\n", strings.TrimSpace(string(output)))
		os.Exit(3)
	}
}

func PushTag(version versioninfo.VersionInfo) {
	cmd := exec.Command("git", "remote")
	remote, err := cmd.Output()
	if err != nil || len(remote) == 0 {
		fmt.Println("Could not get remote name")
		os.Exit(3)
	}
	remoteName := strings.TrimSpace(string(remote))
	cmd = exec.Command("git", "push", remoteName, version.String())
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Could not push %s to %s: %s\n", version.String(), remoteName, err)
		os.Exit(4)
	}
}
