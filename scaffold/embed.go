package scaffold

import "embed"

//go:embed Home.md config.yml deploy.yml AGENTS.md
var Files embed.FS
