package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/richargh/workdown/internal/cli"
	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiracli"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var check bool
	var outputDir string
	flag.BoolVar(&check, "check", false, "check that generated CLI docs are current")
	flag.StringVar(&outputDir, "dir", filepath.Join("docs", "cli"), "CLI docs output directory")
	flag.Parse()

	tempDir, err := os.MkdirTemp("", "workdown-cli-docs-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	if err := generate(tempDir); err != nil {
		return err
	}

	if check {
		if err := compareDirectories(outputDir, tempDir); err != nil {
			return fmt.Errorf("generated CLI docs are not current: %w\nrun go run ./internal/tools/clidocs", err)
		}
		return nil
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := copyDirectory(tempDir, outputDir); err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "generated CLI docs in %s\n", outputDir)
	return err
}

func generate(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return err
	}
	environment := env.NewOSEnv(kernel.UnknownWorkdownVersion, ".", io.Discard, io.Discard)
	root := cli.New(environment, jiracli.New)
	disableAutoGenTag(root)
	return doc.GenMarkdownTree(root, outputDir)
}

func disableAutoGenTag(cmd *cobra.Command) {
	cmd.DisableAutoGenTag = true
	for _, child := range cmd.Commands() {
		disableAutoGenTag(child)
	}
}

func compareDirectories(actualDir string, expectedDir string) error {
	actualFiles, err := markdownFiles(actualDir)
	if err != nil {
		return err
	}
	expectedFiles, err := markdownFiles(expectedDir)
	if err != nil {
		return err
	}
	if !equalStrings(actualFiles, expectedFiles) {
		return fmt.Errorf("file list differs: have %v, want %v", actualFiles, expectedFiles)
	}
	for _, file := range expectedFiles {
		actual, err := os.ReadFile(filepath.Join(actualDir, file)) // #nosec G304 -- file is from the generated docs file list.
		if err != nil {
			return err
		}
		expected, err := os.ReadFile(filepath.Join(expectedDir, file)) // #nosec G304 -- file is from the generated docs file list.
		if err != nil {
			return err
		}
		if !bytes.Equal(actual, expected) {
			return fmt.Errorf("%s differs", filepath.Join(actualDir, file))
		}
	}
	return nil
}

func markdownFiles(root string) ([]string, error) {
	var files []string
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func copyDirectory(sourceDir string, targetDir string) error {
	return filepath.WalkDir(sourceDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(targetDir, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o750)
		}
		data, err := os.ReadFile(path) // #nosec G304,G122 -- path is from the generated docs directory walk.
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o600) // #nosec G703 -- target is below the configured docs output directory.
	})
}

func equalStrings(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
