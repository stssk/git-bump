package main

import (
	"fmt"

	"github.com/stssk/git-bump/git"
	"github.com/stssk/git-bump/utils"
	"github.com/stssk/git-bump/utils/choice"
	"github.com/stssk/git-bump/utils/operation"
)

func main() {
	git.GitInstalled()
	branch := git.CurrentBranch()
	if branch != "main" && branch != "master" {
		question := fmt.Sprintf("Currently on branch %s. Continue?", branch)
		cont := utils.YesNo(question, choice.No)
		if cont != choice.Yes {
			return
		}
	}
	ver := git.GetLastVersion()
	fmt.Printf("Currently on %s\n", ver.String())
	bumpWith := utils.AskForOperation()
	if bumpWith == operation.None {
		return
	}
	if bumpWith == operation.PreRelease {
		if len(ver.PreRelease) == 0 {
			ver.Patch += 1
		}
		ver.PreRelease = utils.AskForPreReleaseVersion()
	} else if bumpWith == operation.Patch {
		ver.Patch += 1
		ver.PreRelease = ""
	} else if bumpWith == operation.Minor {
		ver.Minor += 1
		ver.PreRelease = ""
	} else if bumpWith == operation.Major {
		ver.Major += 1
		ver.PreRelease = ""
	}
	if len(ver.Build) > 0 {
		ver.Build = git.GetSha()
	}
	pushAnswer := utils.YesNo(fmt.Sprintf("Tag and push %s?", ver.String()), choice.Yes)
	if pushAnswer != choice.Yes {
		return
	}
	git.Tag(ver)
	git.PushTag(ver)
	fmt.Printf("Pushed %s", ver.String())
}
