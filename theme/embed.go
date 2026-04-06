package theme

import "embed"

//go:embed default/templates/*.html default/templates/partials/*.html default/static/*
var DefaultTheme embed.FS
