package main

import (
	"fmt"
	"os"

	"github.com/richargh/workdown/internal/cli"
	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiracli"
)

func main() {
	environment := env.NewOSEnv(readVersion(), ".", os.Stdout, os.Stderr)

	root := cli.New(environment, jiracli.New)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readVersion() kernel.WorkdownVersion {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return kernel.UnknownWorkdownVersion
	}

	return kernel.ParseWorkdownVersion(string(data))
}
