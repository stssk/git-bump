package utils

import (
	"errors"
	"os"

	"charm.land/huh/v2"
	"github.com/charmbracelet/x/term"
)

var ErrAborted = errors.New("aborted")

type Prompter interface {
	Select(title string, options []string) (int, error)
	Confirm(question string, defaultYes bool) (bool, error)
	Input(title string, validate func(string) error) (string, error)
}

// huh's own fallback for non-terminals answers with the default when input
// ends, which would accept "Tag and push?" unattended.
func NewPrompter() Prompter {
	if os.Getenv("TERM") == "dumb" || !term.IsTerminal(os.Stdin.Fd()) {
		return NewLinePrompter(os.Stdin, os.Stdout)
	}
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

func (huhPrompter) Input(title string, validate func(string) error) (string, error) {
	var answer string
	err := runField(huh.NewInput().Title(title).Validate(validate).Value(&answer))
	return answer, err
}

func runField(f huh.Field) error {
	err := huh.Run(f)
	if errors.Is(err, huh.ErrUserAborted) {
		return ErrAborted
	}
	return err
}
