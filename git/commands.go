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

func GetLastVersion() versioninfo.VersionInfo {
	cmd := exec.Command("git", "tag", "--list")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not fetch version tags")
		os.Exit(2)
	}
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")

	var highestVersion *versioninfo.VersionInfo
	for _, tag := range tags {
		valid, info := versioninfo.ParseSemver(tag)
		if !valid {
			continue
		}
		if highestVersion == nil {
			highestVersion = &info
		} else {
			if info.Compare(*highestVersion) > 0 {
				highestVersion = &info
			}
		}
	}
	return *highestVersion
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

func Tag(version versioninfo.VersionInfo) {
	cmd := exec.Command("git", "tag", version.String())
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not tag version")
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
