---
name: docs-site
description: Build, preview, and update the VitePress documentation site in docs-site/
license: MIT
metadata:
  audience: contributors
  workflow: docs
---

## What I do

Manage the openclaw-go VitePress documentation site — build it, preview it, add or update pages, and fix issues.

## When to use me

Use this skill when working on the docs site, adding new package documentation, updating guides, or debugging build/deploy issues.

## Site structure

```
docs-site/
  .vitepress/
    config.mts          — VitePress configuration (nav, sidebar, theme)
    theme/
      index.ts          — custom theme extending DefaultTheme
      custom.css         — amber brand colors, dark mode, component styles
  public/
    mascot-gpt.png      — site logo and favicon
  index.md              — home page (hero + feature cards)
  changelog.md          — includes ../CHANGELOG.md via @include directive
  guide/
    getting-started.md
    authentication.md
    streaming.md
  packages/
    gateway.md
    protocol.md
    chatcompletions.md
    openresponses.md
    toolsinvoke.md
    discovery.md
    acp.md
```

## Commands

All commands run from `docs-site/`:

| Command | Purpose |
|---------|---------|
| `npm ci` | Install dependencies (clean) |
| `npm run dev` | Dev server with hot reload |
| `npm run build` | Build static site to `.vitepress/dist/` |
| `npm run preview` | Serve built site locally at `http://localhost:4173/openclaw-go/` |

## Key configuration

- **Base path**: `/openclaw-go/` (GitHub Pages project site)
- **Clean URLs**: enabled (no `.html` extensions)
- **Search**: local provider
- **Logo/favicon**: `/mascot-gpt.png`
- **Edit links**: point to `https://github.com/a3tai/openclaw-go/edit/main/docs-site/:path`

## Adding a new package page

1. Create `docs-site/packages/<name>.md`
2. Add to sidebar in `.vitepress/config.mts` under `/packages/` items
3. Add a feature card to `index.md` if appropriate
4. Build and preview to verify

## Adding a new guide page

1. Create `docs-site/guide/<name>.md`
2. Add to sidebar in `.vitepress/config.mts` under `/guide/` items
3. Add "Next Steps" links from related pages
4. Build and preview to verify

## Deployment

The site deploys to GitHub Pages via `.github/workflows/docs-deploy.yml`:
- Triggered on version tags (`v*`) or manual `workflow_dispatch`
- Uses Node 22, `npm ci`, `npm run build`
- Deploys `.vitepress/dist/` as a Pages artifact

## Rules

- Always build and preview before committing docs changes
- The changelog page auto-includes `../CHANGELOG.md` — do not duplicate content
- Use the amber brand color `#f5a623` for any new styled elements
- Keep code examples consistent with actual library API
