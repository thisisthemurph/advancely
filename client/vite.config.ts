/// <reference types="vitest" />
/// <reference types="vite/client" />

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

export default defineConfig(({ mode }) => {
  const envDir = mode == "production" ? "/etc/secrets" : process.cwd();

  console.log({
    mode,
    cwd: process.cwd(),
  });

  return {
    plugins: [react()],
    envDir,
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./src/test/setup.ts",
      css: false,
    },
  }
});
