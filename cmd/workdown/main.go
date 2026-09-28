package main

import (
	"fmt"
	"os"

	"github.com/richargh/workdown/internal/cli"
	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiracli"
)

func main() {
	environment := env.NewOSEnv(buildVersion(), ".", os.Stdout, os.Stderr)

	root := cli.New(environment, jiracli.New)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
