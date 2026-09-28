package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	changelogPath = "CHANGELOG.md"
	versionPath   = "VERSION"
)

var versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`)

func main() {
	version := flag.String("version", "", "release version to prepare, for example 1.2.3")
	notesFor := flag.String("notes-for", "", "release version to extract notes for, for example 1.2.3")
	output := flag.String("output", "", "file to write release notes to; defaults to stdout")
	flag.Parse()

	fail(validOptions(*version, *notesFor, *output))

	if *version != "" {
		fail(prepareRelease(*version, time.Now().UTC()))
		return
	}

	notes, err := releaseNotes(*notesFor)
	fail(err)
	if *output == "" {
		fmt.Print(notes)
		return
	}
	fail(os.WriteFile(*output, []byte(notes), 0o600))
}

func validOptions(version string, notesFor string, output string) error {
	if (version == "") == (notesFor == "") {
		return errors.New("set exactly one of --version or --notes-for")
	}
	if version != "" && output != "" {
		return errors.New("--output can only be used with --notes-for")
	}
	return nil
}

func prepareRelease(version string, now time.Time) error {
	if err := validVersion(version); err != nil {
		return err
	}

	changelog, err := os.ReadFile(changelogPath)
	if err != nil {
		return err
	}

	updated, err := moveUnreleasedToVersion(string(changelog), version, now.Format("2006-01-02"))
	if err != nil {
		return err
	}

	// #nosec G703 -- changelogPath is a fixed repository file path.
	if err := os.WriteFile(changelogPath, []byte(updated), fileMode(changelogPath)); err != nil {
		return err
	}
	// #nosec G703 -- versionPath is a fixed repository file path.
	return os.WriteFile(versionPath, []byte(version+"\n"), fileMode(versionPath))
}

func validVersion(version string) error {
	if !versionPattern.MatchString(version) {
		return fmt.Errorf("version %q must use semantic version format, for example 1.2.3", version)
	}
	return nil
}

func moveUnreleasedToVersion(changelog string, version string, date string) (string, error) {
	if strings.Contains(changelog, "## ["+version+"]") {
		return "", fmt.Errorf("CHANGELOG.md already contains version %s", version)
	}

	unreleasedStart, nextStart, err := findSection(changelog, "## [Unreleased]")
	if err != nil {
		return "", err
	}

	sectionBody := strings.TrimSpace(changelog[unreleasedStart:nextStart])
	if sectionBody == "" {
		return "", errors.New("CHANGELOG.md has no unreleased notes to release")
	}

	var out bytes.Buffer
	out.WriteString(changelog[:unreleasedStart])
	out.WriteString("\n")
	out.WriteString("## [")
	out.WriteString(version)
	out.WriteString("] - ")
	out.WriteString(date)
	out.WriteString("\n\n")
	out.WriteString(sectionBody)
	out.WriteString("\n\n")
	out.WriteString(changelog[nextStart:])
	return out.String(), nil
}

func releaseNotes(version string) (string, error) {
	changelog, err := os.ReadFile(changelogPath)
	if err != nil {
		return "", err
	}
	return releaseNotesFromChangelog(string(changelog), version)
}

func releaseNotesFromChangelog(changelog string, version string) (string, error) {
	if err := validVersion(version); err != nil {
		return "", err
	}

	start, end, err := findSection(changelog, "## ["+version+"]")
	if err != nil {
		return "", err
	}

	notes := strings.TrimSpace(changelog[start:end])
	if notes == "" {
		return "", fmt.Errorf("CHANGELOG.md has no notes for version %s", version)
	}
	return notes + "\n", nil
}

func findSection(changelog string, heading string) (int, int, error) {
	headingStart, bodyStart := findHeading(changelog, heading)
	if headingStart == -1 {
		return 0, 0, fmt.Errorf("CHANGELOG.md does not contain %q", heading)
	}
	if bodyStart == -1 {
		return 0, 0, fmt.Errorf("CHANGELOG.md heading %q has no body", heading)
	}

	nextHeading := findNextVersionHeading(changelog, bodyStart)
	if nextHeading == -1 {
		return bodyStart, len(changelog), nil
	}
	return bodyStart, nextHeading, nil
}

func findHeading(changelog string, heading string) (int, int) {
	lineStart := 0
	for lineStart < len(changelog) {
		lineEnd := strings.IndexByte(changelog[lineStart:], '\n')
		if lineEnd == -1 {
			lineEnd = len(changelog)
		} else {
			lineEnd += lineStart
		}

		if isRequestedHeading(changelog[lineStart:lineEnd], heading) {
			return lineStart, nextLineStart(changelog, lineEnd)
		}
		lineStart = nextLineStart(changelog, lineEnd)
	}
	return -1, -1
}

func isRequestedHeading(line string, heading string) bool {
	return line == heading || strings.HasPrefix(line, heading+" ")
}

func findNextVersionHeading(changelog string, start int) int {
	lineStart := start
	for lineStart < len(changelog) {
		lineEnd := strings.IndexByte(changelog[lineStart:], '\n')
		if lineEnd == -1 {
			lineEnd = len(changelog)
		} else {
			lineEnd += lineStart
		}

		if strings.HasPrefix(changelog[lineStart:lineEnd], "## [") {
			return lineStart
		}
		lineStart = nextLineStart(changelog, lineEnd)
	}
	return -1
}

func nextLineStart(text string, lineEnd int) int {
	if lineEnd == len(text) {
		return len(text)
	}
	return lineEnd + 1
}

func fileMode(path string) os.FileMode {
	info, err := os.Stat(path)
	if err != nil {
		return 0o600
	}
	return info.Mode().Perm()
}

func fail(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
