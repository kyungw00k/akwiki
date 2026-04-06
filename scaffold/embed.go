package scaffold

import "embed"

//go:embed Home.md config.yml deploy.yml
var Files embed.FS
