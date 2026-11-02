package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrNoInput = errors.New("no answer: input ended")

type linePrompter struct {
	r *bufio.Reader
	w io.Writer
}

func NewLinePrompter(r io.Reader, w io.Writer) Prompter {
	return &linePrompter{r: bufio.NewReader(r), w: w}
}

func (p *linePrompter) readLine() (string, error) {
	line, err := p.r.ReadString('\n')
	if err == io.EOF && line == "" {
		fmt.Fprintln(p.w)
		return "", ErrNoInput
	}
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (p *linePrompter) Select(title string, options []string) (int, error) {
	fmt.Fprintln(p.w, title)
	for i, o := range options {
		fmt.Fprintf(p.w, "  %d) %s\n", i+1, o)
	}
	for {
		fmt.Fprintf(p.w, "Choose 1-%d: ", len(options))
		line, err := p.readLine()
		if err != nil {
			return 0, err
		}
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= len(options) {
			return n - 1, nil
		}
		fmt.Fprintln(p.w, "Not a valid choice")
	}
}

func (p *linePrompter) Confirm(question string, defaultYes bool) (bool, error) {
	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}
	for {
		fmt.Fprintf(p.w, "%s %s ", question, hint)
		line, err := p.readLine()
		if err != nil {
			return false, err
		}
		switch strings.ToLower(line) {
		case "":
			return defaultYes, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Fprintln(p.w, "Please answer y or n")
	}
}

func (p *linePrompter) Input(title string, validate func(string) error) (string, error) {
	for {
		fmt.Fprintf(p.w, "%s: ", title)
		line, err := p.readLine()
		if err != nil {
			return "", err
		}
		if err := validate(line); err != nil {
			fmt.Fprintln(p.w, err)
			continue
		}
		return line, nil
	}
}
