package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

var ErrActionUnregistered = errors.New("action not registered")

type executor struct {
	actions map[string]Action
}

func NewExector() *executor {
	return &executor{actions: make(map[string]Action)}
}

func (e *executor) Run(action string, args []string) {
	var res *Result
	if ac, ok := e.actions[action]; ok {
		res = ac.Run(args)
	} else {
		res = NewFailResult(ErrActionUnregistered)
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	fmt.Printf("%s", unsafe.String(&resBytes[0], len(resBytes)))
}

func (e *executor) Register(actions ...Action) {
	for _, a := range actions {
		e.actions[a.Action()] = a
	}
}
