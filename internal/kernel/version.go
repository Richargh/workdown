package kernel

import "strings"

var UnknownWorkdownVersion = WorkdownVersion{}

type WorkdownVersion struct {
	value string
}

func ParseWorkdownVersion(value string) WorkdownVersion {
	version := strings.TrimRight(value, "\r\n")
	if version == "" || version == "?.?.?" {
		return UnknownWorkdownVersion
	}
	return WorkdownVersion{value: version}
}

func (v WorkdownVersion) String() string {
	if v.value == "" {
		return "?.?.?"
	}
	return v.value
}
