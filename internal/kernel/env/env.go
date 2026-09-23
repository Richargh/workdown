package env

import (
	"io"
	"net/http"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/config"
	"github.com/richargh/workdown/internal/kernel/credentials"
	"github.com/richargh/workdown/internal/kernel/fs"
)

type Env struct {
	Version     kernel.WorkdownVersion
	Config      config.Store
	Files       fs.Store
	Credentials credentials.Source
	HTTPClient  *http.Client
	Stdout      io.Writer
	Stderr      io.Writer
}
