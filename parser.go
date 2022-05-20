package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
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
	regPkgName           = regexp.MustCompile(`package (?P<projPkgName>.*)`)
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
		fmt.Printf("%v", err)
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
