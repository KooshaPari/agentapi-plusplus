import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'agentapi++',
  description: 'Agent API server docs',
  base: '/agentapi-plusplus/',
  themeConfig: {
    nav: [
      { text: 'Start Here', link: '/index' },
      { text: 'Tutorials', link: '/tutorials/' },
      { text: 'How-to', link: '/how-to/' },
      { text: 'Explanation', link: '/explanation/' },
      { text: 'Operations', link: '/operations/' },
      { text: 'API', link: '/api/' }
    ]
  }
})
