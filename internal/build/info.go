package build

import _ "embed"

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
	Name      = "KODKAFA"
	Tagline   = "A Persistent CLI with Memory"
	Url       = "https://kodkafa.com"
	Repo      = "https://github.com/kodkafa/kod"
)

//go:embed logo.txt
var LogoAscii string
