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

type TagKind int

const (
	Lightweight TagKind = iota
	Annotated
	Signed
)

// "--" stops a tag name starting with "-" from being read as an option.
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

func configValue(key string) string {
	output, err := exec.Command("git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func remotes() ([]string, error) {
	output, err := exec.Command("git", "remote").Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(output)), nil
}

func PushRemote(branch, override string) (string, error) {
	configured, err := remotes()
	if err != nil {
		return "", fmt.Errorf("could not list remotes: %w", err)
	}
	if len(configured) == 0 {
		return "", fmt.Errorf("no remotes configured")
	}

	candidate := override
	if candidate == "" && branch != "" {
		candidate = configValue("branch." + branch + ".pushRemote")
	}
	if candidate == "" {
		candidate = configValue("remote.pushDefault")
	}
	if candidate == "" && branch != "" {
		candidate = configValue("branch." + branch + ".remote")
	}
	if candidate == "" {
		if len(configured) != 1 {
			return "", fmt.Errorf("several remotes configured (%s); choose one with -remote", strings.Join(configured, ", "))
		}
		candidate = configured[0]
	}
	for _, r := range configured {
		if r == candidate {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%q is not a configured remote (have: %s)", candidate, strings.Join(configured, ", "))
}

func pushArgs(remote string, version versioninfo.VersionInfo) []string {
	return []string{"push", "--", remote, version.String()}
}

func PushTag(version versioninfo.VersionInfo, remote string) {
	cmd := exec.Command("git", pushArgs(remote, version)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Could not push %s to %s: %s\n%s\n", version.String(), remote, err, strings.TrimSpace(string(output)))
		os.Exit(4)
	}
}
