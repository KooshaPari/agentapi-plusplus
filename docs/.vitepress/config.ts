import { existsSync } from "node:fs";
import { resolve } from "node:path";
import { pathToFileURL } from "node:url";
import { defineConfig } from "vitepress";

let createPhenotypeConfig = (config: Parameters<typeof defineConfig>[0]) =>
  defineConfig(config);

const vendoredConfigPath = resolve(
  process.cwd(),
  "../vendor/phenodocs/packages/docs/config.js",
);

if (existsSync(vendoredConfigPath)) {
  try {
    ({ createPhenotypeConfig } = await import(
      pathToFileURL(vendoredConfigPath).href
    ));
  } catch {
    // CI and standalone doc builds may not include the vendored shared docs package.
  }
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
