package utils

import (
	"errors"
	"testing"

	"github.com/stssk/git-bump/utils/operation"
)

type fakePrompter struct {
	selects  []int
	confirms []bool
	err      error
	shown    [][]string
}

func (f *fakePrompter) Select(_ string, options []string) (int, error) {
	f.shown = append(f.shown, options)
	if f.err != nil {
		return 0, f.err
	}
	i := f.selects[0]
	f.selects = f.selects[1:]
	return i, nil
}

func (f *fakePrompter) Confirm(string, bool) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	b := f.confirms[0]
	f.confirms = f.confirms[1:]
	return b, nil
}

func TestAskForOperation(t *testing.T) {
	tests := []struct {
		pre     bool
		choice  int
		want    operation.Operation
		options int
	}{
		{false, 0, operation.Patch, 3},
		{false, 2, operation.Major, 3},
		{true, 0, operation.PreRelease, 4},
		{true, 1, operation.Patch, 4},
		{true, 3, operation.Major, 4},
	}
	for _, tt := range tests {
		p := &fakePrompter{selects: []int{tt.choice}}
		got, err := AskForOperation(p, tt.pre)
		if err != nil {
			t.Fatalf("pre=%v choice=%d: %v", tt.pre, tt.choice, err)
		}
		if got != tt.want {
			t.Errorf("pre=%v choice=%d: got %v, want %v", tt.pre, tt.choice, got, tt.want)
		}
		if len(p.shown[0]) != tt.options {
			t.Errorf("pre=%v: options = %v, want %d of them", tt.pre, p.shown[0], tt.options)
		}
		if p.shown[0][tt.choice] != operation.Operations[tt.want] {
			t.Errorf("pre=%v: option %d labelled %q, want %q", tt.pre, tt.choice, p.shown[0][tt.choice], operation.Operations[tt.want])
		}
	}
}

func TestAskForOperationReturnsAbort(t *testing.T) {
	_, err := AskForOperation(&fakePrompter{err: ErrAborted}, false)
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("err = %v, want ErrAborted", err)
	}
}
