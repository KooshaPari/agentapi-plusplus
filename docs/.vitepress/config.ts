import { defineConfig } from "vitepress";

let createPhenotypeConfig = (config: Parameters<typeof defineConfig>[0]) =>
  defineConfig(config);

try {
  ({ createPhenotypeConfig } = await import("@phenotype/docs/config"));
} catch {
  // CI and standalone doc builds may not include the vendored shared docs package.
}

export default createPhenotypeConfig({
  title: "agentapi++",
  description: "Agent API server docs",
  base: process.env.GITHUB_ACTIONS ? "/agentapi-plusplus/" : "/",
  srcDir: ".",
  ignoreDeadLinks: true,
  githubOrg: "KooshaPari",
  githubRepo: "agentapi-plusplus",
  nav: [
    { text: "Wiki", link: "/wiki/" },
    { text: "Development Guide", link: "/development-guide/" },
    { text: "Document Index", link: "/document-index/" },
    { text: "API", link: "/api/" },
    { text: "Roadmap", link: "/roadmap/" },
  ],
  sidebar: [
    {
      text: "Categories",
      items: [
        { text: "Wiki", link: "/wiki/" },
        { text: "Development Guide", link: "/development-guide/" },
        { text: "Document Index", link: "/document-index/" },
        { text: "API", link: "/api/" },
        { text: "Roadmap", link: "/roadmap/" },
      ],
    },
  ],
});
