package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site      SiteConfig      `yaml:"site"`
	Build     BuildConfig     `yaml:"build"`
	Analytics AnalyticsConfig `yaml:"analytics"`
	Theme     ThemeConfig     `yaml:"theme"`
}

type SiteConfig struct {
	Title    string `yaml:"title"`
	Author   string `yaml:"author"`
	URL      string `yaml:"url"`
	Language string `yaml:"language"`
}

type BuildConfig struct {
	OutDir    string `yaml:"outDir"`
	PageRoute string `yaml:"pageRoute"`
}

type AnalyticsConfig struct {
	GA string `yaml:"ga"`
}

type ThemeConfig struct {
	Colors ThemeColors `yaml:"colors"`
	Fonts  ThemeFonts  `yaml:"fonts"`
	Layout ThemeLayout `yaml:"layout"`
	Footer ThemeFooter `yaml:"footer"`
	Edit   ThemeEdit   `yaml:"edit"`
}

type ThemeColors struct {
	Background  string `yaml:"background"`
	Text        string `yaml:"text"`
	Link        string `yaml:"link"`
	LinkPrivate string `yaml:"link-private"`
	Accent      string `yaml:"accent"`
}

type ThemeFonts struct {
	Heading string `yaml:"heading"`
	Body    string `yaml:"body"`
	Code    string `yaml:"code"`
}

type ThemeLayout struct {
	MaxWidth  string `yaml:"max-width"`
	TOC       bool   `yaml:"toc"`
	Backlinks bool   `yaml:"backlinks"`
	Related   bool   `yaml:"related"`
	Search    bool   `yaml:"search"`
}

type ThemeFooter struct {
	Copyright string       `yaml:"copyright"`
	Links     []FooterLink `yaml:"links"`
}

type FooterLink struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
}

type ThemeEdit struct {
	URL string `yaml:"url"`
}

// rawThemeLayout uses pointers to distinguish between "not set" and explicit false.
type rawThemeLayout struct {
	MaxWidth  string `yaml:"max-width"`
	TOC       *bool  `yaml:"toc"`
	Backlinks *bool  `yaml:"backlinks"`
	Related   *bool  `yaml:"related"`
	Search    *bool  `yaml:"search"`
}

type rawThemeConfig struct {
	Colors ThemeColors    `yaml:"colors"`
	Fonts  ThemeFonts     `yaml:"fonts"`
	Layout rawThemeLayout `yaml:"layout"`
	Footer ThemeFooter    `yaml:"footer"`
	Edit   ThemeEdit      `yaml:"edit"`
}

type rawConfig struct {
	Site      SiteConfig      `yaml:"site"`
	Build     BuildConfig     `yaml:"build"`
	Analytics AnalyticsConfig `yaml:"analytics"`
	Theme     rawThemeConfig  `yaml:"theme"`
}

func defaults() Config {
	return Config{
		Site: SiteConfig{
			Language: "ko",
		},
		Build: BuildConfig{
			OutDir:    "dist",
			PageRoute: "/pages",
		},
		Theme: ThemeConfig{
			Layout: ThemeLayout{
				TOC:       true,
				Backlinks: true,
				Related:   true,
				Search:    true,
			},
		},
	}
}

// Load reads .akwiki/config.yml from rootDir and merges it with defaults.
// If no config file is found, returns defaults.
func Load(rootDir string) (Config, error) {
	cfg := defaults()

	configPath := filepath.Join(rootDir, ".akwiki", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return cfg, err
	}

	// Merge non-zero string/struct fields
	if raw.Site.Title != "" {
		cfg.Site.Title = raw.Site.Title
	}
	if raw.Site.Author != "" {
		cfg.Site.Author = raw.Site.Author
	}
	if raw.Site.URL != "" {
		cfg.Site.URL = raw.Site.URL
	}
	if raw.Site.Language != "" {
		cfg.Site.Language = raw.Site.Language
	}
	if raw.Build.OutDir != "" {
		cfg.Build.OutDir = raw.Build.OutDir
	}
	if raw.Build.PageRoute != "" {
		cfg.Build.PageRoute = raw.Build.PageRoute
	}
	if raw.Analytics.GA != "" {
		cfg.Analytics.GA = raw.Analytics.GA
	}

	// Theme Colors
	if raw.Theme.Colors.Background != "" {
		cfg.Theme.Colors.Background = raw.Theme.Colors.Background
	}
	if raw.Theme.Colors.Text != "" {
		cfg.Theme.Colors.Text = raw.Theme.Colors.Text
	}
	if raw.Theme.Colors.Link != "" {
		cfg.Theme.Colors.Link = raw.Theme.Colors.Link
	}
	if raw.Theme.Colors.LinkPrivate != "" {
		cfg.Theme.Colors.LinkPrivate = raw.Theme.Colors.LinkPrivate
	}
	if raw.Theme.Colors.Accent != "" {
		cfg.Theme.Colors.Accent = raw.Theme.Colors.Accent
	}

	// Theme Fonts
	if raw.Theme.Fonts.Heading != "" {
		cfg.Theme.Fonts.Heading = raw.Theme.Fonts.Heading
	}
	if raw.Theme.Fonts.Body != "" {
		cfg.Theme.Fonts.Body = raw.Theme.Fonts.Body
	}
	if raw.Theme.Fonts.Code != "" {
		cfg.Theme.Fonts.Code = raw.Theme.Fonts.Code
	}

	// Theme Layout (pointer-based booleans)
	if raw.Theme.Layout.MaxWidth != "" {
		cfg.Theme.Layout.MaxWidth = raw.Theme.Layout.MaxWidth
	}
	if raw.Theme.Layout.TOC != nil {
		cfg.Theme.Layout.TOC = *raw.Theme.Layout.TOC
	}
	if raw.Theme.Layout.Backlinks != nil {
		cfg.Theme.Layout.Backlinks = *raw.Theme.Layout.Backlinks
	}
	if raw.Theme.Layout.Related != nil {
		cfg.Theme.Layout.Related = *raw.Theme.Layout.Related
	}
	if raw.Theme.Layout.Search != nil {
		cfg.Theme.Layout.Search = *raw.Theme.Layout.Search
	}

	// Theme Footer
	if raw.Theme.Footer.Copyright != "" {
		cfg.Theme.Footer.Copyright = raw.Theme.Footer.Copyright
	}
	if len(raw.Theme.Footer.Links) > 0 {
		cfg.Theme.Footer.Links = raw.Theme.Footer.Links
	}

	// Theme Edit
	if raw.Theme.Edit.URL != "" {
		cfg.Theme.Edit.URL = raw.Theme.Edit.URL
	}

	return cfg, nil
}
