/// <reference types="vitest" />
/// <reference types="vite/client" />

import react from "@vitejs/plugin-react-swc";
import {defineConfig, loadEnv} from "vite";
import { resolve } from "path"

export default defineConfig(({ mode }) => {
  const envDir = mode == "production"
    ? resolve(__dirname, "..")
    : process.cwd();

  const env = loadEnv(mode, envDir, "");

  return {
    plugins: [react()],
    define: {
      __APP_ENV__: JSON.stringify(env.APP_ENV),
    },
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./src/test/setup.ts",
      css: false,
    },
  }
});


// https://vitejs.dev/config/
// export default defineConfig({
//   plugins: [react()],
//   test: {
//     globals: true,
//     environment: "jsdom",
//     setupFiles: "./src/test/setup.ts",
//     css: false,
//   },
// });
