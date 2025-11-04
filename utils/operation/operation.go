// Package operation for the user selection of what to bump
package operation

type Operation int

const (
	PreRelease Operation = iota
	Patch
	Minor
	Major
)

var Operations = []string{"Pre release", "Patch", "Minor", "Major"}
