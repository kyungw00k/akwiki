package theme

import "embed"

//go:embed default/templates/*.html default/templates/partials/*.html
var DefaultTheme embed.FS
