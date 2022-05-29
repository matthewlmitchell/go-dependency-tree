package main

import (
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

	// pkgList maps a filename to a [][]string of the form:
	//   filePackageName: [[fileName, importedPackages . . .] [fileName2, . . .]]
	pkgList := make(map[string][][]string, len(pt.filePaths))
	for name, path := range pt.filePaths {
		deps, err := readAndParseFileArray(path, 200, parseGoDependencies)
		name, deps[0] = deps[0], name
		pkgList[name] = append(pkgList[name], deps)

		if err != nil {
			panic(err)
		}
	}

	printToGraph(pkgList)
}
