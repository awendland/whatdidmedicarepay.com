// Snowpack Configuration File
// See all supported options: https://www.snowpack.dev/reference/configuration

/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    public: "/",
    src: "/dist",
  },
  plugins: [
    "@snowpack/plugin-svelte",
    [
      "@snowpack/plugin-build-script",
      { cmd: "postcss", input: [".css"], output: [".css"] },
    ],
    "@snowpack/plugin-typescript",
  ],
  routes: [
    /* Enable an SPA Fallback in development: */
    // {"match": "routes", "src": ".*", "dest": "/index.html"},
  ],
  packageOptions: {
    /* ... */
  },
  devOptions: {
    /* ... */
  },
  optimize: {
    bundle: true,
    minify: true,
    treeshake: true,
    target: "es2018",
  },
  buildOptions: {
    clean: true,
    out: "../server/static",
  },
};
