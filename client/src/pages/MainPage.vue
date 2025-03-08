<template>
  <!-- Auth Form Overlay -->
  <AuthPage
    v-if="!isAuthenticated"
    :isDarkTheme="isDarkTheme"
    @toggle-theme="toggleTheme"
    @login="handleLogin"
  />

  <!-- Main Chat Interface -->
  <div
    v-show="isAuthenticated"
    :class="['chat-container', { 'dark-theme': isDarkTheme }]"
  >
    <div class="main-content">
      <!-- Main chat interface -->
      <ChatInterface :isDarkTheme="isDarkTheme" @toggle-theme="toggleTheme" />
    </div>
  </div>
</template>

<script>
import ChatInterface from "../components/ChatInterface.vue";
import AuthPage from "./AuthPage.vue";

export default {
  name: "MainPage",
  components: {
    ChatInterface,
    AuthPage,
  },
  data() {
    return {
      isAuthenticated:
        Boolean(localStorage.getItem("isAuthenticated")) || false,
      isDarkTheme: localStorage.getItem("isDarkTheme")
        ? JSON.parse(localStorage.getItem("isDarkTheme"))
        : false,
    };
  },

  mounted() {
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },

  methods: {
    handleLogin() {
      this.isAuthenticated = true;
      localStorage.setItem("isAuthenticated", "true");
    },
    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme;
      localStorage.setItem("isDarkTheme", JSON.stringify(this.isDarkTheme));
      document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
    },
  },
};
</script>
