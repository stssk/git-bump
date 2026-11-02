package utils

import (
	"errors"

	"charm.land/huh/v2"
)

var ErrAborted = errors.New("aborted")

type Prompter interface {
	Select(title string, options []string) (int, error)
	Confirm(question string, defaultYes bool) (bool, error)
}

func NewPrompter() Prompter {
	return huhPrompter{}
}

type huhPrompter struct{}

func (huhPrompter) Select(title string, options []string) (int, error) {
	opts := make([]huh.Option[int], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, i)
	}
	var choice int
	err := runField(huh.NewSelect[int]().Title(title).Options(opts...).Value(&choice))
	return choice, err
}

func (huhPrompter) Confirm(question string, defaultYes bool) (bool, error) {
	answer := defaultYes
	err := runField(huh.NewConfirm().Title(question).Affirmative("Yes").Negative("No").Value(&answer))
	return answer, err
}

func runField(f huh.Field) error {
	err := huh.Run(f)
	if errors.Is(err, huh.ErrUserAborted) {
		return ErrAborted
	}
	return err
}
