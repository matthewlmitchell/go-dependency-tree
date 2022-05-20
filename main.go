package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrNoGoModFound      = errors.New("no go.mod found, is this a project?")
	ErrNoModuleNameFound = errors.New("no module name found")
	ErrNoImportsFound    = errors.New("no imports found")
)

func readFileToString(path string, lineReadLimit int32) (*string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileStats.Size() > 10_000_000 {
		log.Printf("WARNING: File %s is really large (size = %d), reading first %d lines", path, fileStats.Size(), lineReadLimit)
	}

	scanner := bufio.NewScanner(file)

	var output string
	var i int32 = 0
	for scanner.Scan() {
		output += fmt.Sprintf("%s\n", scanner.Text())
		if i > lineReadLimit {
			break
		}

		i++

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &output, nil
}

var mainGoPath string = ""
var ErrFoundMain = errors.New("found main.go")

func main() {

	/*
		pkgNames, err := readAndParseFileArray("main.go", 100, parseImports)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%q\n", pkgNames)
	*/

	// Read the first line of go.mod, if it exists, which will contain
	// the name of the current go project's module
	packageRoot, err := readAndParseFileString("go.mod", 10, parsePackageRoot)
	if err != nil {
		panic(err)
	}

	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var fileIsMainGo = func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == "main.go" {
			mainGoPath = path
			return ErrFoundMain
		}

		return nil
	}

	err = filepath.WalkDir(currentDirectory, fileIsMainGo)
	if err != ErrFoundMain && err != nil {
		panic(err)
	}

	fmt.Println(mainGoPath)

	// path, err := findProjectMain()
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(packageRoot)

}
