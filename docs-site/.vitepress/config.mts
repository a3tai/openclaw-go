import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/openclaw-go/',
  title: 'openclaw-go',
  description: 'The Go client library for OpenClaw — connect agents, stream responses, discover gateways.',

  lastUpdated: true,
  cleanUrls: true,
  appearance: 'dark',

  head: [
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { href: 'https://fonts.googleapis.com/css2?family=Outfit:wght@400;500;600;700;800&family=Albert+Sans:ital,wght@0,400;0,500;0,600;1,400&family=JetBrains+Mono:wght@400;500&display=swap', rel: 'stylesheet' }],
    ['link', { rel: 'icon', href: '/openclaw-go/mascot-gpt.png', type: 'image/png' }],
    ['meta', { name: 'theme-color', content: '#f5a623' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'openclaw-go' }],
    ['meta', { property: 'og:description', content: 'The Go client library for OpenClaw' }],
  ],

  themeConfig: {
    logo: { src: '/mascot-gpt.png', alt: 'openclaw-go' },

    nav: [
      { text: 'Guide', link: '/guide/getting-started', activeMatch: '/guide/' },
      { text: 'Packages', link: '/packages/gateway', activeMatch: '/packages/' },
      { text: 'Changelog', link: '/changelog' },
      { text: 'GitHub', link: 'https://github.com/a3tai/openclaw-go' },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Authentication', link: '/guide/authentication' },
            { text: 'Streaming', link: '/guide/streaming' },
          ],
        },
      ],
      '/packages/': [
        {
          text: 'Packages',
          items: [
            { text: 'gateway', link: '/packages/gateway' },
            { text: 'Gateway Methods', link: '/packages/gateway-methods' },
            { text: 'protocol', link: '/packages/protocol' },
            { text: 'identity', link: '/packages/identity' },
            { text: 'chatcompletions', link: '/packages/chatcompletions' },
            { text: 'openresponses', link: '/packages/openresponses' },
            { text: 'toolsinvoke', link: '/packages/toolsinvoke' },
            { text: 'discovery', link: '/packages/discovery' },
            { text: 'acp', link: '/packages/acp' },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/a3tai/openclaw-go' },
    ],

    editLink: {
      pattern: 'https://github.com/a3tai/openclaw-go/edit/main/docs-site/:path',
      text: 'Edit this page on GitHub',
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: '© A3T · Maintained by <a href="https://a3t.ai">A3T</a>',
    },

    search: {
      provider: 'local',
    },
  },
})
