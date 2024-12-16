package action

import (
	"flag"
	"nvim-help/internal/utils"
)

var _ Action = (*modPath)(nil)

type modPath struct {
	CurPath *string
	fs      *flag.FlagSet
}
type modPathData struct {
	Path string `json:"path"`
}

// Usege implements Action.
func (m *modPath) Usege() string {
	panic("unimplemented")
}

func NewModPath() Action {
	mp := &modPath{fs: flag.NewFlagSet("mod-path", flag.ContinueOnError)}
	mp.CurPath = mp.fs.String("path", "", "Current file directory")
	return mp
}

// Run implements Action.
func (m *modPath) Run(args []string) *Result {
	m.fs.Parse(args)

	modPath, err := utils.GetModPath(*m.CurPath)
	if err != nil {
		return NewFailResult(err)
	}

	return NewSuccessResult(&modPathData{
		Path: modPath,
	})
}

func (m *modPath) Action() string {
	return "mod-path"
}
