package testdata

import _ "embed"

var (
	//go:embed assignments.pas
	Assignments string
	//go:embed constants.pas
	Constants string
)
