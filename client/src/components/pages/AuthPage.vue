<template>
  <div class="auth-container">
    <!-- Success register notification -->
    <transition name="slide-fade">
      <div
        v-if="showSuccessNotification"
        class="success-notification"
        key="notification"
      >
        <div class="notification-content">
          <span>Регистрация прошла успешно! Пожалуйста, выполните вход.</span>
          <svg class="checkmark" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" class="checkmark-circle" />
            <path d="M16 8L10.5 16L8 13.5" class="checkmark-path" />
          </svg>
        </div>
      </div>
    </transition>

    <!-- Mobile error message -->
    <ErrorMessageMobile :message="errorMessage" />

    <!-- Auth card -->
    <div class="auth-card" :class="{ 'animate-container': !isLoginMode }">
      <!-- Auth header -->
      <div class="auth-header">
        <transition name="fade-text" mode="out-in">
          <h2 :key="isLoginMode">
            {{ isLoginMode ? "Вход" : "Регистрация" }}
          </h2>
        </transition>
      </div>

      <!-- Auth form -->
      <AuthForm
        ref="authForm"
        :isLoginMode="isLoginMode"
        :serverError="serverError"
        @auth-submit="handleAuthFormSubmit"
        @update:errorMessage="handleErrorMessageUpdate"
        @update:serverError="updateServerError"
      />
    </div>

    <!-- Footer -->
    <div class="auth-footer">
      <button class="mode-toggle" @click="toggleMode">
        {{ isLoginMode ? "Нет аккаунта? Создать" : "Уже есть аккаунт? Войти" }}
      </button>
      <a href="#" class="forgot-password">Забыли пароль?</a>
    </div>

    <!-- Corner logo -->
    <div class="logo-corner">© &laquo;GoChat&raquo; 2025</div>

    <!-- Theme switcher -->
    <button
      class="theme-switch-corner"
      @click="toggleTheme"
      aria-label="Switch theme"
    >
      <svg
        class="theme-icon"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path :d="isDarkTheme ? moonIcon : sunIcon" />
      </svg>
    </button>
  </div>
</template>

<script>
import AuthForm from "../AuthForm.vue";
import ErrorMessageMobile from "../ErrorMessageMobile.vue";
import { login, register } from "@/composables/useAuthAPI.js";

export default {
  name: "AuthPage",
  components: {
    AuthForm,
    ErrorMessageMobile,
  },
  data() {
    return {
      isLoginMode: true,
      showSuccessNotification: false,
      errorMessage: "",
      serverError: "",
      errors: {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      },
      sunIcon:
        "M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8a4 4 0 110 8 4 4 0 010-8z",
      moonIcon: "M21 12.79A9 9 0 1111.21 3a7 7 0 009.79 9.79z",
    };
  },
  computed: {
    isDarkTheme() {
      return this.$store.getters["theme/isDarkTheme"];
    },
  },
  mounted() {
    document.documentElement.classList.toggle("dark-theme", this.isDarkTheme);
  },
  watch: {
    "$route.query.registered": {
      immediate: true,
      handler() {
        this.handleRegistration();
      },
    },
  },
  methods: {
    // resetErrors clears all error messages
    resetErrors() {
      this.errors = {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      };
    },
    // Shows an error message in the mobile error message component if
    // the window width is less than or equal to 568px
    showMobileError() {
      if (window.innerWidth <= 568) {
        let messageToShow =
          this.errorMessage ||
          Object.values(this.errors).find((error) => error) ||
          "";
        if (messageToShow) {
          this.errorMessage = messageToShow;
          if (this.errorTimeout) {
            clearTimeout(this.errorTimeout);
          }
          this.errorTimeout = setTimeout(() => {
            this.errorMessage = "";
          }, 3500);
        } else {
          this.errorMessage = "";
        }
      }
    },
    /**
     * Sends an authentication request to the server, either to login or
     * register.
     *
     * @param {Object} authData - The authentication data containing
     * user credentials.
     *
     * If the request is successful, dispatches an action to update the
     * authentication
     * state and redirects to the chat page if logging in, or to the same
     * page with a query parameter "registered=true" if registering.
     *
     * If the request is unsuccessful, catches the error, extracts the error
     * message from the error object, and calls the handleServerError method
     * of the AuthForm component to display the error message.
     */
    async sendAuthRequest(authData) {
      const APP_ID = 1;
      const requestData = {
        app_id: APP_ID,
        ...authData,
      };
      try {
        let data;
        if (this.isLoginMode) {
          data = await login(requestData);
          console.log("login:", data);
          // const userData = await getUser();
          // console.log("sendAuthRequest: getUser data", userData);
          this.$router.push("/");
        } else {
          data = await register(requestData);
          console.log("register:", data);
          this.$router.replace({
            path: "/auth",
            query: { registered: "true" },
          });
          this.isLoginMode = true;
        }
        this.serverError = "";
      } catch (error) {
        console.error(error);
        let serverError = "Authentication error";
        try {
          const errData = JSON.parse(error.message);
          serverError = errData.error || serverError;
        } catch (e) {
          serverError = error.message || serverError;
        }
        this.serverError = serverError;
        this.$refs.authForm.handleServerError(serverError);
      }
    },

    handleAuthFormSubmit(authData) {
      this.resetErrors();
      if (!authData) {
        this.showMobileError();
        return;
      }
      this.sendAuthRequest(authData);
    },
    handleErrorMessageUpdate(message) {
      this.errorMessage = message;
      this.showMobileError();
    },
    updateServerError(newError) {
      this.serverError = newError;
    },
    handleRegistration() {
      if (this.$route.query.registered === "true") {
        this.showSuccessNotification = true;
        setTimeout(() => {
          this.showSuccessNotification = false;
          this.$router.replace({ query: null });
        }, 5000);
      }
    },
    toggleMode() {
      this.resetErrors();
      this.isLoginMode = !this.isLoginMode;
      this.serverError = "";
    },
    toggleTheme() {
      this.$store.dispatch("theme/toggleTheme");
    },
  },
};
</script>

<style scoped>
@import "@/styles/authpage.css";
</style>
