import { createApp } from "vue";
import App from "./App.vue";

import "bootstrap/dist/css/bootstrap.css"; // to use bootstrap CSS
import "bootstrap/dist/js/bootstrap.bundle.js"; // to use bootstrap JS
import "bootstrap-icons/font/bootstrap-icons.css"; // to use bootstrap icons

// Import CSS style files
import "@/styles/variables.css";
import "@/styles/layout.css";
import "@/styles/animations.css";
import "@/styles/responsive.css";
import "@/styles/theme.css";
import "@/styles/ui.css";

createApp(App).mount("#app");
