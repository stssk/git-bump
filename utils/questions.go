// Package utils Contains questions for the user and the typed restuls
package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/stssk/git-bump/utils/choice"
	"github.com/stssk/git-bump/utils/operation"
	"golang.org/x/term"
)

func YesNo(question string, defaultOption choice.Choice) choice.Choice {
	switch defaultOption {
	case choice.Yes:
		fmt.Printf("%s [Y/n] ", question)
	case choice.No:
		fmt.Printf("%s [y/N] ", question)
	default:
		fmt.Printf("%s [y/n] ", question)
	}
	for {
		b := getYesNoAnswer()
		if b == 13 || b == 10 {
			fmt.Println()
			return defaultOption
		}
		r := rune(b)
		if r == 'Y' || r == 'y' {
			fmt.Println()
			return choice.Yes
		}
		if r == 'N' || r == 'n' {
			fmt.Println()
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

func AskForOperation(p Prompter, usePreRelease bool) (operation.Operation, error) {
	ops := []operation.Operation{operation.Patch, operation.Minor, operation.Major}
	if usePreRelease {
		ops = append([]operation.Operation{operation.PreRelease}, ops...)
	}
	labels := make([]string, len(ops))
	for i, op := range ops {
		labels[i] = operation.Operations[op]
	}
	i, err := p.Select("Bump with", labels)
	if err != nil {
		return 0, err
	}
	return ops[i], nil
}

var preReleaseRegex = regexp.MustCompile(`^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`)

func ValidPreRelease(s string) bool {
	return preReleaseRegex.MatchString(s)
}

func AskForPreReleaseVersion() string {
	preRelease, err := readPreRelease(bufio.NewReader(os.Stdin), os.Stdout)
	if err != nil {
		os.Exit(3)
	}
	return preRelease
}

func readPreRelease(r *bufio.Reader, w io.Writer) (string, error) {
	for {
		fmt.Fprint(w, "Pre release version:")
		line, err := r.ReadString('\n')
		trimmed := strings.TrimSpace(line)
		if ValidPreRelease(trimmed) {
			return trimmed, nil
		}
		if err != nil {
			return "", err
		}
		fmt.Fprintln(w, "Not a valid pre release version")
	}
}
