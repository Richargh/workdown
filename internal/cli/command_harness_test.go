package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func runCommand(t *testing.T, cmd *cobra.Command, args ...string) (stdout string, err error) {
	t.Helper()

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return out.String(), err
}
