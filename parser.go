package main

import (
	"fmt"
	"regexp"
	"strings"
)

func regexpStringSubmatchGroups(re *regexp.Regexp, input *string) (map[string]string, error) {
	regexMatches := make(map[string]string)

	match := re.FindStringSubmatch(*input)
	if match == nil {
		return nil, fmt.Errorf("no matches found")
	}
	if re.SubexpNames() == nil {
		return nil, fmt.Errorf("no subgroups found")
	}

	for i, name := range re.SubexpNames() {
		if i != 0 && strings.TrimSpace(name) != "" {
			regexMatches[name] = match[i]
		}
	}

	return regexMatches, nil
}

func parseImports(input *string) ([]string, error) {
	var re = regexp.MustCompile(`import \((?P<pkgNames>[^)]*)\)`)

	regexMatches, err := regexpStringSubmatchGroups(re, input)
	if err != nil {
		return nil, ErrNoImportsFound
	}

	regexMatches["pkgNames"] = strings.ReplaceAll(regexMatches["pkgNames"], "\"", "")
	regexMatches["pkgNames"] = strings.ReplaceAll(regexMatches["pkgNames"], "\t", "")

	importNamesList := strings.Split(regexMatches["pkgNames"], "\n")

	importedPackages := []string{}
	for _, elem := range importNamesList {
		if strings.TrimSpace(elem) != "" {
			importedPackages = append(importedPackages, elem)
		}
	}

	return importedPackages, nil
}

func parsePackageRoot(input *string) (string, error) {
	var re = regexp.MustCompile(`module (?P<projectRoot>.*)`)

	regexMatches, err := regexpStringSubmatchGroups(re, input)
	if err != nil {
		return "", ErrNoModuleNameFound
	}

	return regexMatches["projectRoot"], nil
}

func readAndParseFileString(path string, linesToRead int32,
	parserFunc func(*string) (string, error)) (string, error) {

	fileContents, err := readFileToString(path, linesToRead)
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

	fileContents, err := readFileToString(path, linesToRead)
	if err != nil {
		return nil, err
	}

	value, err := parserFunc(fileContents)
	if err != nil {
		return nil, err
	}

	return value, nil
}
