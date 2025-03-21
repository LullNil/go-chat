export default {
  namespaced: true,
  state: {
    isDarkTheme: localStorage.getItem("isDarkTheme")
      ? JSON.parse(localStorage.getItem("isDarkTheme"))
      : false,
  },
  mutations: {
    setDarkTheme(state, payload) {
      state.isDarkTheme = payload;
      localStorage.setItem("isDarkTheme", JSON.stringify(payload));
    },
  },
  actions: {
    toggleTheme({ commit, state }) {
      const newTheme = !state.isDarkTheme;
      commit("setDarkTheme", newTheme);
      // Обновляем CSS-класс на корневом элементе для применения темы
      document.documentElement.classList.toggle("dark-theme", newTheme);
    },
  },
  getters: {
    isDarkTheme: (state) => state.isDarkTheme,
  },
};
