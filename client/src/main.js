import { createApp } from "vue";
import App from "./App.vue";

import "bootstrap/dist/css/bootstrap.css"; // to use bootstrap CSS
import "bootstrap/dist/js/bootstrap.bundle.js"; // to use bootstrap JS
import "bootstrap-icons/font/bootstrap-icons.css"; // to use bootstrap icons

import "@/styles/global.css";
import "@/styles/theme.css";
import "@/styles/fonts.css";
import "@/styles/ui.css";

createApp(App).mount("#app");
