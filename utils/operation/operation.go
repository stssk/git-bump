package operation

type Operation int

const (
	None Operation = iota
	PreRelease
	Patch
	Minor
	Major
)

var Operations = []string{"None", "Pre release", "Patch", "Minor", "Major"}
