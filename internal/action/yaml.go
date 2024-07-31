package action

import (
	"errors"
	"flag"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	_                  Action = (*yamlEdit)(nil)
	ErrProjectNotfound        = errors.New("project not found")
	envExample                = EnvConfig{Name: "nameExample", Value: "example"}
)

type yamlEdit struct {
	fs       *flag.FlagSet
	filePath *string
	project  *string
	args     *string
}

func NewYamlEdit() Action {
	ye := &yamlEdit{
		fs: flag.NewFlagSet("yaml-edit", flag.ContinueOnError),
	}
	ye.project = ye.fs.String("p", "", "project name")
	ye.filePath = ye.fs.String("f", "", "config file path")
	ye.args = ye.fs.String("args", "", "args for startup")
	return ye
}

// Action implements Action.
func (y *yamlEdit) Action() string {
	return "yaml-edit"
}

// Run implements Action.
func (y *yamlEdit) Run(args []string) *Result {
	if len(args) < 1 {
		return NewFailResult(ErrIncorrectRequest)
	}

	if err := y.fs.Parse(args[1:]); err != nil {
		return NewFailResult(err)
	}

	switch args[0] {
	case "read":
		data, err := y.read()
		if err != nil {
			return NewFailResult(err)
		}
		return NewSuccessResult(data)
	case "modify":
		if err := y.modify(); err != nil {
			return NewFailResult(err)
		}
		return NewSuccessResult(nil)
	default:
		return NewFailResult(ErrIncorrectRequest)
	}
}

// Usege implements Action.
func (y *yamlEdit) Usege() string {
	panic("unimplemented")
}

func (y *yamlEdit) read() (r *RunConfig, err error) {
	return r, y.upset(func(runConfigs *RunConfigs) {
		found := false
		for _, c := range *runConfigs {
			if c.Project == *y.project {
				r = c
				return
			}
		}

		if !found {
			r = &RunConfig{
				Project: *y.project,
				Env:     []EnvConfig{envExample},
			}
			*runConfigs = append(*runConfigs, r)
		}
	})
}

// TODO Support for modifying other fields
func (y *yamlEdit) modify() (err error) {
	// modify args
	return y.upset(func(runConfigs *RunConfigs) {
		found := false
		for _, c := range *runConfigs {
			if c.Project == *y.project {
				c.Args = *y.args
				found = true
			}
		}

		if !found {
			*runConfigs = append(
				*runConfigs,
				&RunConfig{
					Project: *y.project,
					Args:    *y.args,
					Env:     []EnvConfig{envExample},
				},
			)
		}
	})
}

func (y *yamlEdit) upset(fc func(r *RunConfigs)) (err error) {
	f, err := os.OpenFile(*y.filePath, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()

	fb, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var oldData yaml.Node
	if err = yaml.Unmarshal(fb, &oldData); err != nil {
		return err
	}

	var runConfigs RunConfigs
	if err = oldData.Decode(&runConfigs); err != nil {
		return err
	}

	fc(&runConfigs)

	var newData yaml.Node
	if err = newData.Encode(&runConfigs); err != nil {
		return err
	}

	// copy comment
	var cp func(*yaml.Node, *yaml.Node)
	cp = func(o *yaml.Node, n *yaml.Node) {
		if n.Value == o.Value {
			n.HeadComment = o.HeadComment
			n.LineComment = o.LineComment
			n.FootComment = o.FootComment
		}
		for i := range n.Content {
			if i >= len(o.Content) {
				break
			}
			cp(o.Content[i], n.Content[i])
		}
	}
	// oldData.Kind is docComment
	if len(oldData.Content) > 0 {
		cp(oldData.Content[0], &newData)
	}
	nb, err := yaml.Marshal(newData)
	if err != nil {
		return err
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	size, err := f.Write(nb)
	if err != nil {
		return err
	}

	return f.Truncate(int64(size))
}

type RunConfig struct {
	Project string      `yaml:"project" json:"project"`
	Args    string      `yaml:"args,omitempty" json:"args"`
	Env     []EnvConfig `yaml:"env,omitempty" json:"env"`
}

type EnvConfig struct {
	Name  string
	Value string
}

type RunConfigs []*RunConfig
