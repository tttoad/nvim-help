package utils

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type GoEnv struct {
	GOPATH     string `json:"GOPATH"`
	GOROOT     string `json:"GOROOT"`
	GOOS       string `json:"GOOS"`
	GOARCH     string `json:"GOARCH"`
	GOCACHE    string `json:"GOCACHE"`
	GOMODCACHE string `json:"GOMODCACHE"`
}

func GetGoEnv() (*GoEnv, error) {
	var out bytes.Buffer

	cmd := exec.Command("go", "env", "-json")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	env := &GoEnv{}
	if err := json.Unmarshal(out.Bytes(), env); err != nil {
		return nil, err
	}

	return env, nil
}
