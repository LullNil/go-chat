export default {
  namespaced: true,
  state: {
    token: null,
    isAuthenticated: false,
  },
  mutations: {
    setAuth(state, payload) {
      state.isAuthenticated = payload.isAuthenticated;
      state.token = payload.token || null;
    },
  },
  actions: {
    // Initializing authentication state when the application starts
    async initAuth({ commit }) {
      try {
        const token = await window.electronAPI.getToken();
        console.log("Loaded token from secure store:", token);
        if (token) {
          commit("setAuth", { isAuthenticated: true, token });
        } else {
          commit("setAuth", { isAuthenticated: false, token: null });
        }
      } catch (error) {
        console.error("Error initializing auth:", error);
        commit("setAuth", { isAuthenticated: false, token: null });
      }
    },
    async login({ commit }, data) {
      try {
        // data.token received from server after successful login
        await window.electronAPI.saveToken(data.token);
        commit("setAuth", { isAuthenticated: true, token: data.token });
      } catch (error) {
        console.error("Error during login:", error);
      }
    },
    async logout({ commit }) {
      try {
        await window.electronAPI.deleteToken();
        commit("setAuth", { isAuthenticated: false, token: null });
      } catch (error) {
        console.error("Error during logout:", error);
      }
    },
  },
  getters: {
    isAuthenticated: (state) => state.isAuthenticated,
    token: (state) => state.token,
  },
};
