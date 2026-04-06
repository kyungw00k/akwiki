package builder

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/kyungw00k/akwiki/internal/config"
	"github.com/kyungw00k/akwiki/internal/content"
	"github.com/kyungw00k/akwiki/internal/render"
	"github.com/kyungw00k/akwiki/internal/search"
	"github.com/kyungw00k/akwiki/internal/wiki"
	"github.com/kyungw00k/akwiki/theme"
)

// Build runs the full static-site build pipeline.
func Build(rootDir, outDir string) error {
	// 1. Load config
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. Load pages
	pagesDir := filepath.Join(rootDir, "pages")
	allPages, err := content.LoadPages(pagesDir)
	if err != nil {
		return fmt.Errorf("load pages: %w", err)
	}

	// 3. Filter out private pages, track private names
	privateNames := make(map[string]bool)
	var pages []content.Page
	for _, p := range allPages {
		if p.Private {
			privateNames[p.Name] = true
		} else {
			pages = append(pages, p)
		}
	}

	// 4. Build alias map (alias → page name)
	aliasMap := make(map[string]string)
	for _, p := range pages {
		for _, alias := range p.Aliases {
			aliasMap[alias] = p.Name
		}
	}

	// 5. Build link maps
	bodyMap := make(map[string][]byte, len(pages))
	for _, p := range pages {
		bodyMap[p.Name] = p.RawBody
	}
	linkMap, backlinkMap := wiki.BuildLinkMaps(bodyMap)

	// 6. Build TF-IDF index
	docTexts := make(map[string]string, len(pages))
	for _, p := range pages {
		docTexts[p.Name] = string(p.RawBody)
	}
	tfidfIndex := search.NewTFIDFIndex(docTexts)

	// 7. Create pageByName lookup
	pageByName := make(map[string]content.Page, len(pages))
	for _, p := range pages {
		pageByName[p.Name] = p
	}

	// 8. Prepare output directory
	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("remove old outDir: %w", err)
	}
	for _, d := range []string{outDir, filepath.Join(outDir, "pages"), filepath.Join(outDir, "assets")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}

	// 9. Copy embedded theme static files
	if err := copyEmbeddedStatic(theme.DefaultTheme, outDir); err != nil {
		return fmt.Errorf("copy theme static: %w", err)
	}

	// 10. Copy user's public/ directory to outDir/assets/
	publicDir := filepath.Join(rootDir, "public")
	if info, err := os.Stat(publicDir); err == nil && info.IsDir() {
		if err := copyDir(publicDir, filepath.Join(outDir, "assets")); err != nil {
			return fmt.Errorf("copy public dir: %w", err)
		}
	}

	// 11. Copy user's .akwiki/theme/static/ overrides to outDir/assets/
	themeStaticDir := filepath.Join(rootDir, ".akwiki", "theme", "static")
	if info, err := os.Stat(themeStaticDir); err == nil && info.IsDir() {
		if err := copyDir(themeStaticDir, filepath.Join(outDir, "assets")); err != nil {
			return fmt.Errorf("copy theme static overrides: %w", err)
		}
	}

	// 12. Create template engine
	overrideDir := filepath.Join(rootDir, ".akwiki", "theme", "templates")
	engine, err := render.NewTemplateEngine(theme.DefaultTheme, overrideDir)
	if err != nil {
		return fmt.Errorf("create template engine: %w", err)
	}

	// 13. Build link resolver for wikilinks
	linkResolver := func(target string) wiki.LinkStatus {
		// Check alias map first
		resolvedTarget := target
		if aliasTarget, ok := aliasMap[target]; ok {
			resolvedTarget = aliasTarget
		}

		if privateNames[resolvedTarget] {
			return wiki.LinkPrivate
		}
		if _, ok := pageByName[resolvedTarget]; ok {
			return wiki.LinkExists
		}
		return wiki.LinkMissing
	}

	// 14. Render each non-private page
	for i := range pages {
		p := &pages[i]

		// a. Render markdown to HTML with link resolution
		htmlBytes, err := render.RenderMarkdownWithResolver(p.RawBody, cfg.BasePath()+cfg.Build.PageRoute, linkResolver)
		if err != nil {
			return fmt.Errorf("render markdown for %s: %w", p.Name, err)
		}

		// b. Extract TOC
		toc := wiki.ExtractTOC(p.RawBody)

		// c. Build backlinks list
		var backlinks []render.PageRef
		for _, blName := range backlinkMap[p.Name] {
			if privateNames[blName] {
				continue
			}
			if blPage, ok := pageByName[blName]; ok {
				backlinks = append(backlinks, render.PageRef{
					Name:  blPage.Name,
					Title: blPage.Title,
					Brief: blPage.Brief,
				})
			}
		}

		// d. Build related content (grouped by Type)
		related := make(map[string][]render.PageRef)
		for _, sim := range tfidfIndex.MostSimilar(p.Name, 5) {
			if privateNames[sim.Name] {
				continue
			}
			if relPage, ok := pageByName[sim.Name]; ok {
				typ := relPage.Type
				if typ == "" {
					typ = "Article"
				}
				related[typ] = append(related[typ], render.PageRef{
					Name:  relPage.Name,
					Title: relPage.Title,
					Brief: relPage.Brief,
					Type:  typ,
					Score: sim.Score,
				})
			}
		}

		// e. Build links list
		var links []render.PageRef
		for _, lnkName := range linkMap[p.Name] {
			if privateNames[lnkName] {
				continue
			}
			if lnkPage, ok := pageByName[lnkName]; ok {
				links = append(links, render.PageRef{
					Name:  lnkPage.Name,
					Title: lnkPage.Title,
					Brief: lnkPage.Brief,
				})
			}
		}

		// f. Generate JSON-LD
		jsonld := render.GenerateJSONLD(*p)

		// g. Set page.RawURL
		p.RawURL = cfg.BasePath() + cfg.Build.PageRoute + "/" + url.PathEscape(p.Name) + ".txt"

		// h. Create TemplateContext and render
		ctx := &render.TemplateContext{
			Site:      &cfg,
			Page:      *p,
			Content:   template.HTML(htmlBytes),
			TOC:       toc,
			Links:     links,
			Backlinks: backlinks,
			Related:   related,
			JSONLD:    template.JS(jsonld),
		}

		pageHTML, err := engine.RenderPage(ctx)
		if err != nil {
			return fmt.Errorf("render page %s: %w", p.Name, err)
		}

		// i. Write HTML to outDir/pages/{pageName}/index.html
		pageDir := filepath.Join(outDir, "pages", p.Name)
		if err := os.MkdirAll(pageDir, 0o755); err != nil {
			return fmt.Errorf("create page dir %s: %w", p.Name, err)
		}
		if err := os.WriteFile(filepath.Join(pageDir, "index.html"), pageHTML, 0o644); err != nil {
			return fmt.Errorf("write page HTML %s: %w", p.Name, err)
		}

		// j. Write .txt (raw markdown source)
		if err := os.WriteFile(filepath.Join(outDir, "pages", p.Name+".txt"), p.RawSource, 0o644); err != nil {
			return fmt.Errorf("write page txt %s: %w", p.Name, err)
		}
	}

	// 14. Generate search index
	searchData := search.BuildSearchIndex(pages)
	if err := os.WriteFile(filepath.Join(outDir, "search-index.json"), searchData, 0o644); err != nil {
		return fmt.Errorf("write search index: %w", err)
	}

	// 15. Generate index.html redirect to Home
	homePath := cfg.BasePath() + cfg.Build.PageRoute + "/Home"
	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="refresh" content="0; url=%s">
  <title>Redirecting...</title>
</head>
<body>
  <a href="%s">Redirecting to Home...</a>
</body>
</html>
`, homePath, homePath)
	if err := os.WriteFile(filepath.Join(outDir, "index.html"), []byte(indexHTML), 0o644); err != nil {
		return fmt.Errorf("write index.html: %w", err)
	}

	return nil
}

// copyEmbeddedStatic walks "default/static" in the embedded FS and copies files to outDir/assets/.
func copyEmbeddedStatic(fsys embed.FS, outDir string) error {
	return fs.WalkDir(fsys, "default/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path from "default/static"
		rel, err := filepath.Rel("default/static", path)
		if err != nil {
			return err
		}

		dst := filepath.Join(outDir, "assets", rel)

		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}

		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
}

// copyDir copies the contents of src directory into dst directory recursively.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}
