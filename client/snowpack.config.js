// Snowpack Configuration File
// See all supported options: https://www.snowpack.dev/reference/configuration
// svelte + typescript + tailwind support based on https://github.com/GarrettCannon/snowpack-svelte-ts-tw

const httpProxy = require("http-proxy");

const apiProxy = httpProxy.createServer({
  target: process.env.SNOWPACK_API_HOST || "http://localhost:3000",
});

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
    {
      src: "/api/.*",
      dest: (req, res) => apiProxy.web(req, res),
    },
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
