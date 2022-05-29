package main

import (
	"fmt"
)

func printToGraph(pkgList map[string][][]string) {
	fmt.Println("digraph mygraph {")
	fmt.Println("\tfontname=\"Helvetica,Arial,sans-serif\"")
	fmt.Println("\tnode [fontName=\"Helvetica,Arial,sans-serif\"]")
	fmt.Println("\tedge [fontName=\"Helvetica,Arial,sans-serif\"]")
	fmt.Println("\tnode [shape=box];")

	for name, elem := range pkgList {

		for _, fileDeps := range elem {
			fmt.Printf("\t\"%s\" -> \"%s\"\n", name, fileDeps[0])

			for i := 1; i < len(fileDeps); i++ {
				fmt.Printf("\t\"%s\" -> \"%s\"\n", fileDeps[0], fileDeps[i])
			}
		}
	}

	fmt.Println("}")
}
