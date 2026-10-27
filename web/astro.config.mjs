import { defineConfig } from "astro/config";
import { fileURLToPath } from "node:url";

export default defineConfig({
  output: "static",
  outDir: fileURLToPath(new URL("./dist/", import.meta.url)),
  publicDir: fileURLToPath(new URL("./public/", import.meta.url)),
  vite: {
    server: {
      proxy: {
        "/api": "http://localhost:8080",
      },
    },
  },
});