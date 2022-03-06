package testdata

import _ "embed"

var (
	//go:embed 1.pas
	File1 string
	//go:embed 2.pas
	File2 string
)
