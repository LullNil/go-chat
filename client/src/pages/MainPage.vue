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
      <!-- Combined Chat Component -->
      <ChatComponent :isDarkTheme="isDarkTheme" @toggle-theme="toggleTheme" />
    </div>
  </div>
</template>

<script>
import ChatComponent from "../components/ChatComponent.vue";
import AuthPage from "./AuthPage.vue";

export default {
  components: {
    ChatComponent,
    AuthPage,
  },

  data() {
    return {
      // User authentication state
      isAuthenticated:
        Boolean(localStorage.getItem("isAuthenticated")) || false,
      // Theme state (dark/light)
      isDarkTheme: localStorage.getItem("isDarkTheme")
        ? JSON.parse(localStorage.getItem("isDarkTheme"))
        : false,
    };
  },

  mounted() {
    // Apply theme on initial load
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },

  methods: {
    // Set authentication state to true and hide auth form
    handleLogin() {
      this.isAuthenticated = true;
      localStorage.setItem("isAuthenticated", "true");
    },

    // Toggle theme between dark and light modes
    toggleTheme() {
      this.isDarkTheme = !this.isDarkTheme;
      localStorage.setItem("isDarkTheme", JSON.stringify(this.isDarkTheme));
      document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
    },
  },
};
</script>
