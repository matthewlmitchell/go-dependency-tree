package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	pt := &pather{
		filePaths: make(map[string]string),
	}

	// Read the first line of go.mod, if it exists, which will contain
	// the name of the current go project's module
	var err error
	pt.moduleName, err = readAndParseFileString("go.mod", 10, parseGoMod)
	if err != nil {
		panic(err)
	}

	// Determine the current working directory of the executable
	pt.packageRoot, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// This populates pather.filePaths with the file path to main.go, if found
	err = filepath.WalkDir(pt.packageRoot, pt.fileIsMainGo)
	if err != ErrFoundMain && err != nil {
		panic(err)
	}

	_, err = readAndParseFileArray(pt.filePaths["main.go"], 50, parseGoDependencies)
	if err != nil {
		panic(err)
	}

	// pkgList maps a filename to a []string of the form:
	//   fileName: [filePackageName, importedPackages . . .]
	pkgList := make(map[string][]string, len(pt.filePaths))
	for name, path := range pt.filePaths {
		pkgList[name], err = readAndParseFileArray(path, 500, parseGoDependencies)
		if err != nil {
			fmt.Errorf("err: %s", path)
			panic(err)
		}
	}

	for name, elem := range pkgList {
		fmt.Printf("%s: %q \n", name, elem)
	}

}
