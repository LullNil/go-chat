import { getUser, logout } from "@/composables/useAuthAPI";

export default {
  namespaced: true,
  state: {
    token: null,
    isAuthenticated: false,
  },
  mutations: {
    setAuth(state, payload) {
      state.isAuthenticated = payload.isAuthenticated;
      state.user = payload.user || null;
    },
  },
  actions: {
    // Initializing authentication state when the application starts
    async initAuth({ commit }) {
      try {
        const userData = await getUser();
        commit("setAuth", { isAuthenticated: true, user: userData });
      } catch (error) {
        // If the error is due to a missing cookie, the user is not authenticated
        if (error.message.includes("named cookie not present")) {
          console.warn("User is not authenticated.");
          commit("setAuth", { isAuthenticated: false, user: null });
        } else {
          console.error("Error initializing auth:", error);
          commit("setAuth", { isAuthenticated: false, user: null });
        }
      }
    },

    async logout({ commit }) {
      try {
        await logout();
        commit("setAuth", { isAuthenticated: false, user: null });
      } catch (error) {
        console.error("Error during logout:", error);
      }
    },
  },
  getters: {
    isAuthenticated: (state) => state.isAuthenticated,
    user: (state) => state.user,
  },
};
