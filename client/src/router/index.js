import { createRouter, createWebHashHistory } from "vue-router";
import AuthPage from "../components/pages/AuthPage.vue";
import ChatPage from "../components/pages/ChatPage.vue";
import store from "../store";

const routes = [
  {
    path: "/",
    name: "ChatPage",
    component: ChatPage,
    meta: { requiresAuth: true },
  },
  {
    path: "/auth",
    name: "AuthPage",
    component: AuthPage,
  },
  {
    path: "/:catchAll(.*)",
    redirect: "/auth",
  },
];

const router = createRouter({
  history: createWebHashHistory(), // Usee hash mode for Electron
  routes,
});

router.beforeEach(async (to, from, next) => {
  if (!store.getters["auth/token"]) {
    await store.dispatch("auth/initAuth");
  }
  const isAuthenticated = store.getters["auth/isAuthenticated"];

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: "AuthPage" });
  } else {
    next();
  }
});

export default router;
