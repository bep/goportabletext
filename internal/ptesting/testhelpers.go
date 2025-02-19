// Copyright 2025 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package ptesting

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// IsCI reports whether the tests are running in a CI environment.
func IsCI() bool {
	return os.Getenv("CI") != ""
}

type BlockContentToMarkdownHelper interface {
	SkipFiles() map[string]bool
	GenerateGolden() error
	ForEachGoldenPair(func(sourceName, targetName, sourceContent, targetContent string) error) error
}

func NewUpstreamBlockContentToMarkdownHelper() (BlockContentToMarkdownHelper, error) {
	return NewBlockContentToMarkdownHelper("testdata/upstream", "testdata/golden_upstream/markdown")
}

func NewSamplesBlockContentToMarkdownHelper() (BlockContentToMarkdownHelper, error) {
	return NewBlockContentToMarkdownHelper("testdata/samples", "testdata/golden_samples/markdown")
}

func NewBlockContentToMarkdownHelper(sourceDir, targetDir string) (BlockContentToMarkdownHelper, error) {
	_, bb, _, _ := runtime.Caller(0)
	workingDir := filepath.Dir(bb)
	baseDir := filepath.Join(workingDir, "..", "..")

	sourceDir = filepath.Join(baseDir, sourceDir)
	targetDir = filepath.Join(baseDir, targetDir)

	b := &blockContentToMarkdown{
		workingDir: workingDir,
		SourceDir:  sourceDir,
		TargetDir:  targetDir,
	}
	return b, b.init()
}

// blockContentToMarkdown is a test helper that converts JSON block content files to Markdown.
type blockContentToMarkdown struct {
	workingDir string

	SourceDir string
	TargetDir string
}

func (b *blockContentToMarkdown) SkipFiles() map[string]bool {
	skipFilesFilename := filepath.Join(b.SourceDir, "skipfiles.txt")
	skipFilesContent, err := os.ReadFile(skipFilesFilename)
	if err != nil {
		panic(err)
	}
	skipFiles := make(map[string]bool)

	for _, file := range strings.Split(string(skipFilesContent), "\n") {
		skipFiles[strings.TrimSpace(file)] = true
	}

	return skipFiles
}

func (b *blockContentToMarkdown) ForEachGoldenPair(fn func(sourceName, targetName, sourceContent, targetContent string) error) error {
	skipFiles := b.SkipFiles()

	sourceDirFiles, err := os.ReadDir(b.SourceDir)
	if err != nil {
		return err
	}

	for _, file := range sourceDirFiles {
		if skipFiles[file.Name()] || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		sourceName := file.Name()
		targetName := sourceName + ".md"
		sourceFilename := filepath.Join(b.SourceDir, sourceName)
		targetFilename := filepath.Join(b.TargetDir, targetName)

		sourceContent, err := os.ReadFile(sourceFilename)
		if err != nil {
		}
		targetContent, err := os.ReadFile(targetFilename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Create the target file and fill it with some content that will fail the test.
				targetContent = []byte("FAIL")
				if err := os.WriteFile(targetFilename, targetContent, 0o644); err != nil {
					return err
				}
			}
			return err
		}

		if err := fn(sourceName, targetName, string(sourceContent), string(targetContent)); err != nil {
			return err
		}

	}

	return nil
}

func (b *blockContentToMarkdown) init() error {
	if b.SourceDir == "" || b.TargetDir == "" {
		return errors.New("sourceDir and targetDir must be set")
	}
	return nil
}

func (b *blockContentToMarkdown) GenerateGolden() error {
	// Read all JSON files in sourceDir.
	// For each JSON file, convert it to Markdown.
	// Write the Markdown to targetDir.
	sourceDirFiles, err := os.ReadDir(b.SourceDir)
	if err != nil {
		return err
	}

	// Clear out the targetDir.
	if err := os.RemoveAll(b.TargetDir); err != nil {
		return err
	}
	if err := os.MkdirAll(b.TargetDir, 0o755); err != nil {
		return err
	}

	skipFiles := b.SkipFiles()

	for _, file := range sourceDirFiles {
		if skipFiles[file.Name()] || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		filename := filepath.Join(b.SourceDir, file.Name())
		outFilename := filepath.Join(b.TargetDir, file.Name()+".md")
		if err := func() error {
			outFile, err := os.Create(outFilename)
			if err != nil {
				return err
			}
			defer outFile.Close()
			fmt.Println(filename)
			script := filepath.Join(b.workingDir, "block-to-md.js")
			cmd := exec.Command(script, filename)
			cmd.Dir = b.workingDir
			cmd.Stdout = outFile
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}(); err != nil {
			return err
		}
	}

	return nil
}
