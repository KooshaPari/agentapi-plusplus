import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

import { defineConfig } from 'vitepress'
import { resolveDocsBase } from '../../../docs-hub/.vitepress/base.config'

const docsBase = resolveDocsBase()

export default defineConfig({
  title: 'agentapi++',
  description: 'Agent API server docs',
  base: docsBase,

  vite: {
    resolve: {
      alias: {
        '@phenodocs-theme': phenodocsTheme,
      },
    },
    server: {
      fs: {
        allow: [phenodocsRoot],
      },
    },
  },
  themeConfig: {
    nav: [
      { text: 'Wiki', link: '/wiki/' },
      { text: 'Development Guide', link: '/development-guide/' },
      { text: 'Document Index', link: '/document-index/' },
      { text: 'API', link: '/api/' },
      { text: 'Roadmap', link: '/roadmap/' }
    ],
    sidebar: [{ text: 'Categories', items: [
      { text: 'Wiki', link: '/wiki/' },
      { text: 'Development Guide', link: '/development-guide/' },
      { text: 'Document Index', link: '/document-index/' },
      { text: 'API', link: '/api/' },
      { text: 'Roadmap', link: '/roadmap/' }
    ] }]
  }
})
