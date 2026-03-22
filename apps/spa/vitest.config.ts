import { defineConfig } from "vitest/config";
import path from "path";

export default defineConfig({
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./app/test/setup.ts"],
    include: ["app/**/*.{test,spec}.{ts,tsx}"],
  },
  resolve: {
    tsconfigPaths: true,
    alias: { "~": path.resolve(__dirname, "./app") },
  },
});
