// Package utils Contains questions for the user and the typed restuls
package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/stssk/git-bump/utils/operation"
)

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
