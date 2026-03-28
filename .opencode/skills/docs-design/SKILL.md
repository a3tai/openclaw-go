---
name: docs-design
description: Improve the visual design and UX of the VitePress docs site — layout, colors, typography, components, and responsiveness
license: MIT
metadata:
  audience: designers
  workflow: docs
---

## What I do

Analyze and improve the visual design, user experience, and polish of the openclaw-go documentation site built with VitePress.

## When to use me

Use this skill when the docs site needs visual improvements, when adding new UI components, or when the site feels rough or inconsistent.

## Design system

### Brand colors
- **Primary (amber)**: `#f5a623`
- **Primary hover**: `#d4891e`
- **Primary deep**: `#b8730f`
- **Primary soft**: `rgba(245, 166, 35, 0.12)`

### Dark mode surfaces
- **Background**: `#0d0d0f`
- **Soft background**: `#1a1a1f`
- **Mute background**: `#141418`
- **Alt background**: `#101013`
- **Border/divider**: `#2a2a35`

### Mascot
- File: `docs-site/public/mascot-gpt.png` — a Go gopher with lobster claws
- Used as: navbar logo (64x64px) and favicon
- Style: playful but professional

## Key files to edit

| File | What it controls |
|------|-----------------|
| `docs-site/.vitepress/theme/custom.css` | All custom styles — colors, component overrides, dark mode |
| `docs-site/.vitepress/theme/index.ts` | Theme setup, layout component registration |
| `docs-site/.vitepress/config.mts` | Nav, sidebar, logo, footer, head meta |
| `docs-site/index.md` | Home page hero, feature cards, install block |

## Design principles

1. **Dark-first** — the site is primarily viewed in dark mode. Light mode should work but dark mode is the priority.
2. **Readable code** — code blocks need high contrast, clear syntax highlighting, and generous padding.
3. **Amber accents** — use the brand amber sparingly for links, active states, and interactive elements. Don't overwhelm with color.
4. **Whitespace** — generous spacing between sections. Don't crowd content.
5. **Mobile-friendly** — test at 375px width. Navigation, code blocks, and tables should be usable on mobile.

## Improvement areas to consider

### Homepage
- Hero section: gradient, tagline sizing, CTA button styling
- Feature cards: hover effects, icon sizing, card grid layout
- Install/quickstart section: code block styling, spacing

### Navigation
- Logo sizing and alignment in the navbar
- Active state styling for nav items
- Sidebar width and item spacing
- Mobile hamburger menu behavior

### Content pages
- Heading hierarchy and spacing
- Table styling (borders, padding, zebra striping)
- Inline code vs code block contrast
- Tip/warning/info callout boxes
- "On this page" ToC sidebar styling

### Typography
- Font stack for body, headings, and code
- Line height and paragraph spacing
- Heading sizes and weight progression

### Interactions
- Link hover states
- Button focus rings
- Smooth scroll behavior
- Page transition effects

## Workflow

1. **Assess current state** — take screenshots, identify issues
2. **Prioritize** — fix broken/ugly things first, polish second
3. **Edit CSS** — make changes in `custom.css`
4. **Build and preview**:
   ```sh
   cd docs-site && npm run build && npm run preview
   ```
5. **Compare** — screenshot before/after
6. **Test responsiveness** — check at 375px, 768px, and 1440px widths
7. **Test both themes** — verify changes work in dark and light mode

## Rules

- Stay within VitePress's theming system — override CSS variables and class selectors, don't fight the framework
- Keep custom CSS maintainable — comment sections, use the existing color variables
- Don't add JavaScript for purely visual effects — CSS-only where possible
- Test in both dark and light mode even though dark is primary
