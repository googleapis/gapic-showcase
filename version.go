package version

import (
	_ "embed"
)

// Version contains the current version of this project, as declared in the
// version.txt file, for its tools to reference.
//go:embed version.txt
var Version string
