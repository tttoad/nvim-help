package action

const v version = "0.0.2"

var _ Action = (version)("ss")

type version string

func NewVersionAction() Action {
	return v
}

// Action implements Action.
func (v version) Action() string {
	return "version"
}

// Run implements Action.
func (v version) Run([]string) *Result {
	return NewSuccessResult(v)
}

// Usege implements Action.
func (v version) Usege() string {
	panic("unimplemented")
}
