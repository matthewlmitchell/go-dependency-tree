package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	ErrNoGoModFound      = errors.New("no go.mod found, is this a project?")
	ErrNoModuleNameFound = errors.New("no module name found")
	ErrNoPkgNameFound    = errors.New("no package name found")
	ErrNoImportsFound    = errors.New("no imports found")
	ErrNoMatchesFound    = errors.New("no matches found")
	ErrNoSubgroupsFound  = errors.New("no subgroups found")
	regPkgImports        = regexp.MustCompile(`import \((?P<pkgNames>[^)]*)\)`)
	regPkgName           = regexp.MustCompile(`(?m)^package\s(?P<projPkgName>\S+)`)
)

func regexpStringSubmatchGroups(re *regexp.Regexp, input *string) (map[string]string, error) {
	regexMatches := make(map[string]string)

	match := re.FindStringSubmatch(*input)
	if match == nil {
		return nil, ErrNoMatchesFound
	}
	if re.SubexpNames() == nil {
		return nil, ErrNoSubgroupsFound
	}

	for i, name := range re.SubexpNames() {
		if i != 0 && strings.TrimSpace(name) != "" {
			regexMatches[name] = match[i]
		}
	}

	return regexMatches, nil
}

func parsePackageName(input *string) (string, error) {
	regexMatches, err := regexpStringSubmatchGroups(regPkgName, input)
	if err != nil {
		return "", err
	}

	projPkgName, ok := regexMatches["projPkgName"]
	if !ok {
		return "", ErrNoPkgNameFound
	}

	return projPkgName, nil
}

func parseGoDependencies(input *string) ([]string, error) {

	projPkgName, err := parsePackageName(input)
	if err != nil {
		return nil, err
	}

	regexMatches, err := regexpStringSubmatchGroups(regPkgImports, input)
	if err == ErrNoMatchesFound {
		return []string{projPkgName}, nil
	} else if err != nil {
		return nil, err
	}

	regexMatches["pkgNames"] = strings.ReplaceAll(regexMatches["pkgNames"], "\"", "")
	regexMatches["pkgNames"] = strings.ReplaceAll(regexMatches["pkgNames"], "\t", "")

	importNamesList := strings.Split(regexMatches["pkgNames"], "\n")

	importedPackages := []string{}
	importedPackages = append(importedPackages, projPkgName)
	for _, elem := range importNamesList {
		if strings.TrimSpace(elem) != "" {
			importedPackages = append(importedPackages, elem)
		}
	}

	return importedPackages, nil
}

func parseGoMod(input *string) (string, error) {
	var re = regexp.MustCompile(`module (?P<projectRoot>.*)`)

	regexMatches, err := regexpStringSubmatchGroups(re, input)
	if err != nil {
		return "", ErrNoModuleNameFound
	}

	return regexMatches["projectRoot"], nil
}

func readAndParseFileString(path string, linesToRead int32,
	parserFunc func(*string) (string, error)) (string, error) {

	fileContents, err := readPkgImportsToString(path, linesToRead)
	if err != nil {
		return "", err
	}

	value, err := parserFunc(fileContents)
	if err != nil {
		return "", err
	}

	return value, nil

}

func readAndParseFileArray(path string, linesToRead int32,
	parserFunc func(*string) ([]string, error)) ([]string, error) {

	fileContents, err := readPkgImportsToString(path, linesToRead)
	if err != nil {
		return nil, err
	}

	// fmt.Println(*fileContents)

	value, err := parserFunc(fileContents)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func readPkgImportsToString(path string, lineReadLimit int32) (*string, error) {
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
		fmt.Printf("WARNING: File %s is really large (size = %d), reading first %d lines", path, fileStats.Size(), lineReadLimit)
	}

	scanner := bufio.NewScanner(file)

	var output string
	var currentLine string
	ignoreLine := false
	foundImportBlock := false

	var i int32 = 0
	for scanner.Scan() {
		currentLine = scanner.Text()

		// Remove trailing and leading whitespace so we can properly detect
		// comments in indented code blocks
		currentLine = strings.TrimSpace(currentLine)

		// If we find a multi-line comment that doesnt start with "/*",
		// save the line then skip the next ones until we find "*/"
		if strings.Contains(currentLine, "/*") &&
			!strings.Contains(currentLine, "*/") {
			if !strings.HasPrefix(currentLine, "/*") {
				output += fmt.Sprintf("%s\n", scanner.Text())
				i++
				ignoreLine = true
				continue
			}

			// When the multi-line comment starts with "/*",
			// ignore the entire line
			ignoreLine = true
		} else if strings.HasPrefix(currentLine, "//") {
			// If the line is a single-line comment, skip the entire line
			continue
		} else if ignoreLine && strings.Contains(currentLine, "*/") {
			// If we find the end of the multi-line comment, skip the line
			// and dont ignore the next one
			ignoreLine = false
			continue
		} else if !ignoreLine && strings.HasPrefix(currentLine, "import (") {
			// If we find the beginning of a multi-line import block that is not
			// on an ignored line
			foundImportBlock = true
		} else if !ignoreLine && foundImportBlock &&
			strings.TrimSpace(currentLine) == ")" {
			// If we are in an import block and we find a line only containing
			// the ending delimiter ")", add the line and exit
			i = lineReadLimit + 1
		} else if !ignoreLine && strings.HasPrefix(currentLine, "import \"") {
			i = lineReadLimit + 1
		}

		if !ignoreLine {
			output += fmt.Sprintf("%s\n", scanner.Text())
			i++
		}

		if i > lineReadLimit {
			break
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &output, nil
}
