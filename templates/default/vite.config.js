import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
  plugins: [
    VitePWA({
      registerType: "autoUpdate",
      injectRegister: "inline", 
      workbox: {
        globPatterns: ["**/*.{js,css,html,ico,png,svg}"],
      },
      manifest: {
        name: "{{PROJECT_NAME}}",
        short_name: "{{PROJECT_NAME}}",
        description: "A local-first application built with LIAB.",
        theme_color: "#ffffff",
        background_color: "#ffffff",
        display: "standalone",
        start_url: "/",
        icons: [
          {
            src: "https://js.org", // Stand-in asset icon
            sizes: "512x512",
            type: "image/png"
          }
        ]
      }
    })
  ]
});
