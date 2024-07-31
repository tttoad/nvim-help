package action

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

var (
	_                  Action = (*modPath)(nil)
	ErrModFileNotFount        = errors.New("mod file not found")
)

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
	curPath, err := filepath.Abs(filepath.Dir(*m.CurPath))
	if err != nil {
		return NewFailResult(err)
	}

	for curPath != "" {
		fileList, err := os.ReadDir(curPath)
		if err != nil {
			return NewFailResult(err)
		}

		for _, file := range fileList {
			if !file.IsDir() && file.Name() == "go.mod" {
				return NewSuccessResult(&modPathData{
					Path: curPath,
				})
			}
		}
		curPath = filepath.Join(curPath, "../")
	}

	return NewFailResult(ErrModFileNotFount)
}

func (m *modPath) Action() string {
	return "mod-path"
}
