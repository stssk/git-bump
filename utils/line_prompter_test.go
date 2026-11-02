package utils

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLinePrompterConfirm(t *testing.T) {
	tests := []struct {
		input      string
		defaultYes bool
		want       bool
	}{
		{"\n", true, true},
		{"\n", false, false},
		{"y\n", false, true},
		{"YES\n", false, true},
		{"n\n", true, false},
		{" No \n", true, false},
		{"maybe\ny\n", false, true},
		{"y", false, true},
	}
	for _, tt := range tests {
		got, err := NewLinePrompter(strings.NewReader(tt.input), io.Discard).Confirm("Continue?", tt.defaultYes)
		if err != nil || got != tt.want {
			t.Errorf("input %q default %v: got %v, %v; want %v", tt.input, tt.defaultYes, got, err, tt.want)
		}
	}
}

func TestLinePrompterConfirmDoesNotDefaultAtEOF(t *testing.T) {
	for _, input := range []string{"", "maybe\n"} {
		got, err := NewLinePrompter(strings.NewReader(input), io.Discard).Confirm("Tag and push?", true)
		if !errors.Is(err, ErrNoInput) || got {
			t.Errorf("input %q: got %v, %v; want false, ErrNoInput", input, got, err)
		}
	}
}

func TestLinePrompterSelect(t *testing.T) {
	var out strings.Builder
	p := NewLinePrompter(strings.NewReader("0\n4\nabc\n99999999999999999999\n2\n"), &out)
	got, err := p.Select("Bump with", []string{"Patch", "Minor", "Major"})
	if err != nil || got != 1 {
		t.Fatalf("got %d, %v; want 1, nil", got, err)
	}
	if n := strings.Count(out.String(), "Not a valid choice"); n != 4 {
		t.Fatalf("rejections = %d, want 4; output:\n%s", n, out.String())
	}
	if !strings.Contains(out.String(), "  2) Minor\n") {
		t.Fatalf("options not listed; output:\n%s", out.String())
	}
}

func TestLinePrompterSelectStopsAtEOF(t *testing.T) {
	_, err := NewLinePrompter(strings.NewReader("9\n"), io.Discard).Select("Bump with", []string{"Patch"})
	if !errors.Is(err, ErrNoInput) {
		t.Fatalf("err = %v, want ErrNoInput", err)
	}
}

func TestLinePrompterInputReprompts(t *testing.T) {
	var out strings.Builder
	p := NewLinePrompter(strings.NewReader("rc 1\n\nrc.1\n"), &out)
	got, err := AskForPreReleaseVersion(p)
	if err != nil || got != "rc.1" {
		t.Fatalf("got %q, %v; want rc.1, nil", got, err)
	}
	if n := strings.Count(out.String(), "not a valid pre release version"); n != 2 {
		t.Fatalf("rejections = %d, want 2; output:\n%s", n, out.String())
	}
}

func TestLinePrompterInputStopsAtEOF(t *testing.T) {
	_, err := AskForPreReleaseVersion(NewLinePrompter(strings.NewReader("bad input\n"), io.Discard))
	if !errors.Is(err, ErrNoInput) {
		t.Fatalf("err = %v, want ErrNoInput", err)
	}
}

func TestLinePrompterSharesInputBetweenPrompts(t *testing.T) {
	p := NewLinePrompter(strings.NewReader("y\n1\nrc.1\nn\n"), io.Discard)
	if ok, err := p.Confirm("Continue?", false); err != nil || !ok {
		t.Fatalf("Confirm: %v, %v", ok, err)
	}
	if i, err := p.Select("Bump with", []string{"a", "b"}); err != nil || i != 0 {
		t.Fatalf("Select: %d, %v", i, err)
	}
	if s, err := AskForPreReleaseVersion(p); err != nil || s != "rc.1" {
		t.Fatalf("Input: %q, %v", s, err)
	}
	if ok, err := p.Confirm("Push?", true); err != nil || ok {
		t.Fatalf("Confirm: %v, %v", ok, err)
	}
}

func TestNewPrompterFallsBackWithoutTerminal(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	stdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = stdin }()

	if _, ok := NewPrompter().(*linePrompter); !ok {
		t.Fatal("piped stdin: want line prompter")
	}
}

func TestNewPrompterFallsBackOnDumbTerminal(t *testing.T) {
	t.Setenv("TERM", "dumb")
	if _, ok := NewPrompter().(*linePrompter); !ok {
		t.Fatal("TERM=dumb: want line prompter")
	}
}
