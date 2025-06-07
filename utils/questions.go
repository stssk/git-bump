package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stssk/git-bump/utils/choice"
	"github.com/stssk/git-bump/utils/operation"
)

func YesNo(question string, defaultOption choice.Choice) choice.Choice {
	reader := bufio.NewReader(os.Stdin)
	if defaultOption == choice.Yes {
		fmt.Printf("%s [Y/n] ", question)
	} else if defaultOption == choice.No {
		fmt.Printf("%s [y/N] ", question)
	} else {
		fmt.Printf("%s [y/n] ", question)
	}
	r, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln(err)
	}
	if r == 'Y' || r == 'y' {
		return choice.Yes
	}
	if r == 'N' || r == 'n' {
		return choice.No
	}
	return defaultOption
}

func AskForOperation() operation.Operation {
	promptRuns := &survey.Select{
		Message: "Bump with",
		Options: operation.Operations,
	}

	answer := 0
	survey.AskOne(promptRuns, &answer, survey.WithValidator(survey.Required))
	return operation.Operation(answer)
}

func AskForPreReleaseVersion() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pre release version:")
	r, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(3)
	}
	trimmed := strings.TrimSpace(r)
	preReleaseRegex := `[0-9A-Za-z-]`
	re := regexp.MustCompile(preReleaseRegex)
	matches := re.MatchString(trimmed)
	if matches {
		return trimmed
	}
	fmt.Println("Not a valid pre release version")
	return AskForPreReleaseVersion()
}
