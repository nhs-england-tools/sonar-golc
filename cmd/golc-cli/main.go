package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SonarSource-Demos/sonar-golc/assets"
	"github.com/SonarSource-Demos/sonar-golc/pkg/goloc"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a directory\n", absDir)
		os.Exit(1)
	}

	params := goloc.Params{
		Path:              absDir,
		ByFile:            false,
		ExcludePaths:      []string{},
		ExcludeExtensions: []string{},
		IncludeExtensions: []string{},
		FolderKeywords:    []string{},
		FileNamePatterns:  []string{},
		OrderByCode:       true,
		Order:             "DESC",
		OutputName:        "Result_",
		OutputPath:        "",
		ReportFormats:     []string{"prompt"},
		Cloned:            true,
		Repopath:          absDir,
	}

	gc, err := goloc.NewGCloc(params, assets.Languages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing analysis: %v\n", err)
		os.Exit(1)
	}

	if err := gc.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during analysis: %v\n", err)
		os.Exit(1)
	}
}
