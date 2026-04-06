# akwiki Project Guide

## Project Structure

```
pages/          — Markdown wiki pages
public/         — Static assets (images, etc.)
.akwiki/        — Site configuration
dist/           — Built static site (generated)
```

## Creating Pages

Add `.md` files to `pages/`. The file name (without extension) becomes the page URL.

### Frontmatter

```yaml
---
title: Page Title
titleKo: 페이지 제목       # optional Korean title
aliases: [alt name, 다른 이름]  # optional search aliases
tags: [tag1, tag2]           # optional tags
---
```

## Wiki Syntax

- `[[Page Name]]` — link to another wiki page
- `[[Display Text|Page Name]]` — link with custom text
- Standard Markdown for everything else (headings, lists, code, tables, etc.)

## Commands

```bash
akwiki dev      # Local preview with auto-rebuild on file changes
akwiki build    # Generate static site to dist/
akwiki serve    # Serve built site on localhost:3000
```

## Conventions

- File names become page URLs (`pages/My Page.md` → `/pages/My Page`)
- Use wikilinks (`[[...]]`) to connect pages — backlinks are generated automatically
- Keep pages focused on single topics
- Related content is suggested automatically via TF-IDF similarity
