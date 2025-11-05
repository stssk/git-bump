package main

import (
	"flag"
	"fmt"

	"github.com/stssk/git-bump/git"
	"github.com/stssk/git-bump/utils"
	"github.com/stssk/git-bump/utils/choice"
	"github.com/stssk/git-bump/utils/operation"
)

func main() {
	var usePreRelease = flag.Bool("pre", false, "Support creating pre release versions")
	var noPush = flag.Bool("no-push", false, "Do not push tag after tagging")
	flag.Parse()

	git.GitInstalled()
	branch := git.CurrentBranch()
	if branch != "main" && branch != "master" {
		question := fmt.Sprintf("Currently on branch %s. Continue?", branch)
		cont := utils.YesNo(question, choice.No)
		if cont != choice.Yes {
			return
		}
	}
	currentTag := git.HeadTagged()
	if currentTag != "" {
		question := fmt.Sprintf("Current branch is tagged with %s. Continue?", currentTag)
		cont := utils.YesNo(question, choice.No)
		if cont != choice.Yes {
			return
		}
	}
	ver := git.GetLastVersion()
	fmt.Printf("Currently on %s\n", ver.String())
	bumpWith := utils.AskForOperation(usePreRelease)
	if bumpWith < 0 {
		return
	}
	switch bumpWith {
	case operation.PreRelease:
		if len(ver.PreRelease) == 0 {
			ver.Patch += 1
		}
		ver.PreRelease = utils.AskForPreReleaseVersion()
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
	if len(ver.Build) > 0 {
		ver.Build = git.GetSha()
	}
	if *noPush {
		git.Tag(ver)
		fmt.Printf("Tagged %s\n", ver.String())
	} else {
		pushAnswer := utils.YesNo(fmt.Sprintf("Tag and push %s?", ver.String()), choice.Yes)
		if pushAnswer != choice.Yes {
			return
		}
		git.Tag(ver)
		git.PushTag(ver)
		fmt.Printf("Pushed %s\n", ver.String())
	}
}
