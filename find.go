package main

import (
	"errors"
	"os"
	"strings"
)

var (
	ErrFoundMain = errors.New("found main.go")
)

type pather struct {
	filePaths   map[string]string
	moduleName  string
	packageRoot string
	walkDirFunc func(string, os.DirEntry, error) error
}

func (p *pather) fileIsMainGo(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if strings.HasSuffix(d.Name(), ".go") && !d.IsDir() {
		p.filePaths[d.Name()] = path
	}

	return nil
}
