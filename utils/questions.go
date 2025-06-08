package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stssk/git-bump/utils/choice"
	"github.com/stssk/git-bump/utils/operation"
	"golang.org/x/term"
)

func YesNo(question string, defaultOption choice.Choice) choice.Choice {
	if defaultOption == choice.Yes {
		fmt.Printf("%s [Y/n] ", question)
	} else if defaultOption == choice.No {
		fmt.Printf("%s [y/N] ", question)
	} else {
		fmt.Printf("%s [y/n] ", question)
	}
	for {
		b := getYesNoAnswer()
		if b == 13 || b == 10 {
			return defaultOption
		}
		r := rune(b)
		if r == 'Y' || r == 'y' {
			return choice.Yes
		}
		if r == 'N' || r == 'n' {
			return choice.No
		}
	}
}

func getYesNoAnswer() byte {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic("Failed to initialise terminal")
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		panic("Could not read keyboard")
	}
	if b[0] == 3 {
		os.Exit(0)
	}
	return b[0]
}

func AskForOperation() operation.Operation {
	promptRuns := &survey.Select{
		Message: "Bump with",
		Options: operation.Operations,
	}

	answer := -1
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
