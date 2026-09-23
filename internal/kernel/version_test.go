package kernel_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel"
)

func TestParseWorkdownVersion(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "semantic version", value: "0.0.0", want: "0.0.0"},
		{name: "trims unix newline", value: "0.0.0\n", want: "0.0.0"},
		{name: "trims windows newline", value: "0.0.0\r\n", want: "0.0.0"},
		{name: "empty is unknown", value: "", want: "?.?.?"},
		{name: "explicit unknown", value: "?.?.?", want: "?.?.?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			version := kernel.ParseWorkdownVersion(tt.value)
			// then
			require.Equal(t, tt.want, version.String())
		})
	}
}

func TestWorkdownVersionStringUnknownValues(t *testing.T) {
	tests := []struct {
		name    string
		version kernel.WorkdownVersion
	}{
		{name: "unknown sentinel", version: kernel.UnknownWorkdownVersion},
		{name: "zero value", version: kernel.WorkdownVersion{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// then
			require.Equal(t, "?.?.?", tt.version.String())
		})
	}
}
