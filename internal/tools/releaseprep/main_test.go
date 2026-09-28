package main

import (
	"strings"
	"testing"
)

func TestReleasePreparationOptions(t *testing.T) {
	t.Run("rejects output when preparing a release", func(t *testing.T) {
		// given
		testee := validOptions

		// when
		err := testee("1.2.3", "", "notes.md")

		// then
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestReleasePreparationMovesManualChangelogNotes(t *testing.T) {
	t.Run("moves unreleased notes into a dated version section", func(t *testing.T) {
		// given
		testee := moveUnreleasedToVersion
		changelog := `# Changelog

## [Unreleased]

### Added

- Added one thing.

## [0.0.0] - 2026-09-23

### Added

- Initial version.
`

		// when
		got, err := testee(changelog, "1.2.3", "2026-09-28")

		// then
		if err != nil {
			t.Fatalf("move unreleased notes: %v", err)
		}
		want := `# Changelog

## [Unreleased]

## [1.2.3] - 2026-09-28

### Added

- Added one thing.

## [0.0.0] - 2026-09-23

### Added

- Initial version.
`
		if got != want {
			t.Fatalf("changelog mismatch\nwant:\n%s\ngot:\n%s", want, got)
		}
	})
}

func TestReleasePreparationRejectsInvalidChangelogState(t *testing.T) {
	t.Run("rejects a release when unreleased notes are empty", func(t *testing.T) {
		// given
		testee := moveUnreleasedToVersion
		changelog := `# Changelog

## [Unreleased]

## [0.0.0] - 2026-09-23
`

		// when
		_, err := testee(changelog, "1.2.3", "2026-09-28")

		// then
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestChangelogSectionParsing(t *testing.T) {
	t.Run("ignores matching text outside heading lines", func(t *testing.T) {
		// given
		testee := releaseNotesFromChangelog
		changelog := `# Changelog

## [Unreleased]

- The text ## [1.2.3] is not a release heading.

## [1.2.3] - 2026-09-28

- Release item.
`

		// when
		got, err := testee(changelog, "1.2.3")

		// then
		if err != nil {
			t.Fatalf("release notes: %v", err)
		}
		if strings.Contains(got, "not a release heading") {
			t.Fatalf("release notes contain text from wrong section: %s", got)
		}
		if !strings.Contains(got, "Release item") {
			t.Fatalf("release notes do not contain release item: %s", got)
		}
	})
}

func TestGitHubReleaseNotesComeFromRequestedChangelogVersion(t *testing.T) {
	t.Run("extracts only the requested version section", func(t *testing.T) {
		// given
		testee := releaseNotesFromChangelog
		changelog := `# Changelog

## [Unreleased]

## [1.2.3] - 2026-09-28

### Added

- Added one thing.

## [0.0.0] - 2026-09-23
`

		// when
		got, err := testee(changelog, "1.2.3")

		// then
		if err != nil {
			t.Fatalf("release notes: %v", err)
		}
		if strings.Contains(got, "0.0.0") {
			t.Fatalf("release notes contain next release section: %s", got)
		}
		if !strings.Contains(got, "Added one thing") {
			t.Fatalf("release notes do not contain release item: %s", got)
		}
	})
}
