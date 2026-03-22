import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss(), reactRouter()],
  server: {
    host: "::",
    port: 3000,
  },
  resolve: {
    tsconfigPaths: true,
  },
});
