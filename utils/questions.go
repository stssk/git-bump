// Package utils Contains questions for the user and the typed restuls
package utils

import (
	"errors"
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

func AskForPreReleaseVersion(p Prompter) (string, error) {
	answer, err := p.Input("Pre release version", validatePreRelease)
	return strings.TrimSpace(answer), err
}

func validatePreRelease(s string) error {
	if !ValidPreRelease(strings.TrimSpace(s)) {
		return errors.New("not a valid pre release version")
	}
	return nil
}
