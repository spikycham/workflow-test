import { defineConfig } from "vite";

export default defineConfig({
    base: "/cicd/",
    server: {
        open: true,
    },
    test: {
        globals: true,
    },
});
