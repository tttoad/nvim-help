package utils

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

var ErrModFileNotFount = errors.New("mod file not found")

func GetModPath(curPath string) (modPath string, err error) {
	if curPath, err = filepath.Abs(filepath.Dir(curPath)); err != nil {
		return "", nil
	}

	for curPath != "" {
		fileList, err := os.ReadDir(curPath)
		if err != nil {
			return "", err
		}

		for _, file := range fileList {
			if !file.IsDir() && file.Name() == "go.mod" {
				return curPath, nil
			}
		}
		curPath = filepath.Join(curPath, "../")
	}

	return "", ErrModFileNotFount
}

func GetGoVersion(modFile string) (goVersion string, err error) {
	mf, err := os.ReadFile(modFile)
	if err != nil {
		return "", err
	}

	mfp, err := modfile.ParseLax(modFile, mf, nil)
	if err != nil {
		return "", err
	}

	return mfp.Go.Version, nil
}
