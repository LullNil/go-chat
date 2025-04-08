import { getUserProfile, logout } from "@/composables/useAuthAPI";

/**
 * Normalize a user object from the API into a shape that's consistent with
 * what the rest of the app expects. This includes adding default values for
 * first_name, last_name, and avatar_url if they are not present in the API
 * response.
 *
 * @param {Object} apiUser - The user object from the API.
 * @returns {Object|null} - The normalized user object, or null if the input
 * is invalid.
 */
function normalizeUserData(apiUser) {
  if (!apiUser) return null;
  return {
    userId: apiUser.user_id,
    email: apiUser.email,
    username: apiUser.username,
    firstName: apiUser.first_name ?? "",
    lastName: apiUser.last_name ?? "",
    avatarUrl: apiUser.avatar_url ?? null,
    createdAt: apiUser.created_at,
    updatedAt: apiUser.updated_at,
  };
}

export default {
  namespaced: true,
  state: {
    isAuthenticated: false,
    user: null,
  },
  mutations: {
    /**
     * Updates the authentication state in the Vuex store with data from the
     * API. This mutation is used both when the user logs in and when the user
     * data is updated.
     *
     * @param {Object} state - The Vuex state object
     * @param {Object} payload - The payload with data from the API
     * @param {boolean} payload.isAuthenticated - Whether the user is authenticated
     * @param {Object} payload.user - The user data from the API
     */
    setAuth(state, payload) {
      state.isAuthenticated = payload.isAuthenticated;
      state.user = normalizeUserData(payload.user);
    },

    /**
     * Updates the user data in the store with data from the API, after
     * normalizing it. This is used to update the user data in the store
     * when the user logs in or when the user data is updated.
     *
     * @param {Object} state - The Vuex state object
     * @param {Object} userDataFromApi - The user data from the API
     */
    SET_USER(state, userDataFromApi) {
      const normalizedUserUpdate = normalizeUserData(userDataFromApi);
      if (normalizedUserUpdate) {
        state.user = normalizedUserUpdate;
      } else {
        console.warn("Received invalid data in SET_USER mutation");
      }
    },
  },
  actions: {
    /**
     * Initializes the authentication state by calling the getUserProfile
     * function from useAuthAPI and committing the result to the Vuex store.
     * If the call to getUserProfile fails, it sets isAuthenticated to false
     * and user to null.
     *
     * @param {Object} context - The Vuex context object
     * @param {function} context.commit - The Vuex mutation commit function
     * @returns {Promise<void>}
     */
    async initAuth({ commit }) {
      try {
        const apiUserData = await getUserProfile();
        commit("setAuth", { isAuthenticated: true, user: apiUserData });
      } catch (error) {
        console.warn(
          "initAuth: User not authenticated:",
          error.message || error
        );
        commit("setAuth", { isAuthenticated: false, user: null });
      }
    },

    /**
     * Logs out the user and updates the Vuex store with the new
     * authentication state.
     *
     * @param {Object} context - The Vuex context object
     * @param {function} context.commit - The Vuex mutation commit function
     * @returns {Promise<void>}
     */
    async logout({ commit }) {
      try {
        await logout();
        commit("setAuth", { isAuthenticated: false, user: null });
      } catch (error) {
        console.error("Error during logout API call:", error);
        commit("setAuth", { isAuthenticated: false, user: null });
      }
    },

    /**
     * Updates the user data in the Vuex store with the data from the API
     * by committing the SET_USER mutation.
     *
     * @param {Object} context - The Vuex context object
     * @param {function} context.commit - The Vuex mutation commit function
     * @param {Object} userDataFromApi - The user data received from the API
     * @returns {undefined}
     */
    updateUser({ commit }, userDataFromApi) {
      commit("SET_USER", userDataFromApi);
    },
  },
  getters: {
    isAuthenticated: (state) => state.isAuthenticated,
    user: (state) => state.user,
  },
};
