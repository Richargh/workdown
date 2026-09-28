package main

import "github.com/richargh/workdown/internal/kernel"

var version = "?.?.?"

func buildVersion() kernel.WorkdownVersion {
	return kernel.ParseWorkdownVersion(version)
}
