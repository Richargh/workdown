package env

import (
	"io"
	"net/http"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/fs"
)

func NewOSEnv(version kernel.WorkdownVersion, root string, stdout io.Writer, stderr io.Writer) Env {
	return Env{
		Version:    version,
		Files:      fs.NewOSStore(root),
		HTTPClient: http.DefaultClient,
		Stdout:     stdout,
		Stderr:     stderr,
	}
}
