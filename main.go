package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/stssk/git-bump/git"
	"github.com/stssk/git-bump/utils"
	"github.com/stssk/git-bump/utils/operation"
	"github.com/stssk/git-bump/versioninfo"
)

func main() {
	var usePreRelease = flag.Bool("pre", false, "Support creating pre release versions")
	var noPush = flag.Bool("no-push", false, "Do not push tag after tagging")
	var annotate = flag.Bool("annotate", false, "Create an annotated tag")
	var sign = flag.Bool("sign", false, "Create a signed tag (implies -annotate); uses git's signing configuration")
	var remoteFlag = flag.String("remote", "", "Remote to push the tag to (default: the current branch's push remote)")
	var prefix = flag.String("prefix", "", "Only consider version tags with this prefix, e.g. \"v\" (use -prefix= for no prefix)")
	flag.Parse()
	prefixSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "prefix" {
			prefixSet = true
		}
	})

	git.GitInstalled()
	prompter := utils.NewPrompter()
	branch := git.CurrentBranch()
	if branch != "main" && branch != "master" {
		question := fmt.Sprintf("Currently on branch %s. Continue?", branch)
		cont, err := prompter.Confirm(question, false)
		exitOnPromptError(err)
		if !cont {
			return
		}
	}
	currentTag := git.HeadTagged()
	if currentTag != "" {
		question := fmt.Sprintf("Current branch is tagged with %s. Continue?", currentTag)
		cont, err := prompter.Confirm(question, false)
		exitOnPromptError(err)
		if !cont {
			return
		}
	}
	ver, hasTags := git.GetLastVersion(*prefix, prefixSet)
	if !hasTags {
		fmt.Println("is not tagged")
		ver = versioninfo.VersionInfo{
			Prefix: *prefix,
			Major:  0,
			Minor:  0,
			Patch:  0,
		}
	}
	fmt.Printf("Currently on %s\n", ver.String())
	ver, err := nextVersion(prompter, ver, *usePreRelease)
	exitOnPromptError(err)
	if len(ver.Build) > 0 {
		ver.Build = git.GetSha()
	}
	tagKind := git.Lightweight
	if *sign {
		tagKind = git.Signed
	} else if *annotate {
		tagKind = git.Annotated
	}
	if *noPush {
		git.Tag(ver, tagKind)
		fmt.Printf("Tagged %s\n", ver.String())
	} else {
		remote, err := git.PushRemote(branch, *remoteFlag)
		if err != nil {
			fmt.Printf("Could not determine remote: %s\n", err)
			os.Exit(3)
		}
		push, err := prompter.Confirm(fmt.Sprintf("Tag and push %s to %s?", ver.String(), remote), true)
		exitOnPromptError(err)
		if !push {
			return
		}
		git.Tag(ver, tagKind)
		git.PushTag(ver, remote)
		fmt.Printf("Pushed %s to %s\n", ver.String(), remote)
	}
}

func exitOnPromptError(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, utils.ErrAborted) {
		os.Exit(130)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func nextVersion(p utils.Prompter, ver versioninfo.VersionInfo, allowPreRelease bool) (versioninfo.VersionInfo, error) {
	bumpWith, err := utils.AskForOperation(p, allowPreRelease)
	if err != nil {
		return ver, err
	}
	switch bumpWith {
	case operation.PreRelease:
		preRelease, err := utils.AskForPreReleaseVersion(p)
		if err != nil {
			return ver, err
		}
		if len(ver.PreRelease) == 0 {
			ver.Patch += 1
		}
		ver.PreRelease = preRelease
	case operation.Patch:
		ver.Patch += 1
		ver.PreRelease = ""
	case operation.Minor:
		ver.Minor += 1
		ver.Patch = 0
		ver.PreRelease = ""
	case operation.Major:
		ver.Major += 1
		ver.Minor = 0
		ver.Patch = 0
		ver.PreRelease = ""
	}
	return ver, nil
}
