<template>
  <div class="auth-container">
    <!-- Success registration notification -->
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

    <!-- Authentication card -->
    <div class="auth-card" :class="{ 'animate-container': !isLoginMode }">
      <!-- Header -->
      <div class="auth-header">
        <transition name="fade-text" mode="out-in">
          <h2 :key="isLoginMode">
            {{ isLoginMode ? "Вход" : "Регистрация" }}
          </h2>
        </transition>
      </div>

      <form @submit.prevent="handleSubmit" novalidate>
        <!-- Username input (only for registration) -->
        <div
          class="extra-field username-field"
          :class="{ active: !isLoginMode }"
        >
          <div class="auth-input-group">
            <input
              v-model="username"
              type="text"
              placeholder=" "
              class="auth-input"
              @input="clearError('username')"
            />
            <label
              class="auth-input-label"
              :class="{ 'error-label': errors.username }"
            >
              Имя пользователя
            </label>
            <ErrorMessage :message="errors.username" />
          </div>
        </div>

        <!-- Email input group -->
        <div class="auth-input-group">
          <input
            v-model="email"
            type="email"
            placeholder=" "
            class="auth-input"
            @input="clearError('email')"
          />
          <label
            class="auth-input-label"
            :class="{ 'error-label': errors.email }"
          >
            Почта
          </label>
          <ErrorMessage :message="errors.email" />
        </div>

        <!-- Password input group -->
        <div class="auth-input-group">
          <input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            placeholder=" "
            class="auth-input"
            @input="clearError('password')"
          />
          <label
            class="auth-input-label"
            :class="{ 'error-label': errors.password }"
          >
            Пароль
          </label>

          <!-- Password toggle button -->
          <button type="button" class="password-toggle" @click="togglePassword">
            <svg
              class="eye-icon"
              :class="{ 'animate-closed': !showPassword }"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="butt"
              stroke-linejoin="round"
            >
              <path
                d="M1 12c1.67-4.33 5.477-7 11-7s9.33 2.67 11 7c-1.67 4.33-5.477 7-11 7s-9.33-2.67-11-7zM12 15a3 3 0 100-6 3 3 0 000 6z"
              />
              <line
                class="eye-cross"
                :class="{ drawn: !showPassword }"
                x1="4"
                y1="4"
                x2="20"
                y2="20"
              />
            </svg>
          </button>
          <ErrorMessage :message="errors.password" />
        </div>

        <!-- Confirm Password input (only for registration) -->
        <div
          class="extra-field confirm-field"
          :class="{ active: !isLoginMode }"
        >
          <div class="auth-input-group">
            <input
              v-model="confirmPassword"
              type="password"
              placeholder=" "
              class="auth-input"
              @input="clearError('confirmPassword')"
            />
            <label
              class="auth-input-label"
              :class="{ 'error-label': errors.confirmPassword }"
            >
              Подтвердите пароль
            </label>
            <ErrorMessage :message="errors.confirmPassword" />
          </div>
        </div>

        <!-- Submit button -->
        <button type="submit" class="auth-button">
          {{ isLoginMode ? "Войти" : "Зарегистрироваться" }}
        </button>
      </form>
    </div>

    <!-- Footer -->
    <div class="auth-footer">
      <button class="mode-toggle" @click="toggleMode">
        {{ isLoginMode ? "Нет аккаунта? Создать" : "Уже есть аккаунт? Войти" }}
      </button>
      <a href="#" class="forgot-password">Забыли пароль?</a>
    </div>

    <!-- Logo -->
    <div class="logo-corner">© &laquo;GoChat&raquo; 2025</div>
    <!-- Theme switch button -->
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
import ErrorMessage from "../ErrorMessage.vue";
import ErrorMessageMobile from "../ErrorMessageMobile.vue";
import { login, register } from "@/composables/useAuthAPI.js";

export default {
  name: "AuthPage",
  components: {
    ErrorMessage,
    ErrorMessageMobile,
  },
  data() {
    return {
      isLoginMode: true,
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      showPassword: false,
      showSuccessNotification: false,
      errorMessage: "",
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
    resetErrors() {
      this.errors = {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      };
    },
    clearError(field) {
      this.errors[field] = "";
    },
    showMobileError() {
      const firstError = Object.values(this.errors).find((error) => error);
      if (firstError) {
        this.errorMessage = firstError;
        if (this.errorTimeout) clearTimeout(this.errorTimeout);
        this.errorTimeout = setTimeout(() => {
          this.errorMessage = "";
        }, 3500);
      }
    },
    /**
     * Validates the form and returns true if the form is valid, false otherwise.
     * Will update the `errors` object with any errors it finds.
     * @returns {boolean} Whether the form is valid or not.
     */
    validateForm() {
      let hasErrors = false;
      if (!this.email.trim()) {
        this.errors.email = "Введите почту.";
        hasErrors = true;
      } else if (!/.+@.+\..+/.test(this.email)) {
        this.errors.email = "Неверный формат почты.";
        hasErrors = true;
      }
      if (!this.password.trim()) {
        this.errors.password = "Введите пароль.";
        hasErrors = true;
      }
      if (!this.isLoginMode) {
        if (!this.username.trim()) {
          this.errors.username = "Введите имя пользователя.";
          hasErrors = true;
        }
        if (this.password !== this.confirmPassword) {
          this.errors.confirmPassword = "Пароли не совпадают.";
          hasErrors = true;
        }
      }
      return !hasErrors;
    },
    /**
     * Sends authentication request to the server.
     *
     * If the authentication was successful, it logs in the user via Vuex store
     * and redirects to the main page.
     *
     * @throws {Error} - if the authentication request failed
     */
    async sendAuthRequest() {
      const APP_ID = 1;
      const authData = {
        email: this.email,
        password: this.password,
        app_id: APP_ID,
        ...(this.isLoginMode ? {} : { username: this.username }),
      };

      try {
        let data;
        if (this.isLoginMode) {
          // When logging in, expect a JWT token
          data = await login(authData);
          console.log("accepted token:", data.token);
          await this.$store.dispatch("auth/login", data);
          this.$router.push("/");
        } else {
          // When registering, expect a userId or success flag to go to the login page
          data = await register(authData);
          console.log("uid:", data.userId, "status:", data.status);
          this.$router.replace({
            path: "/auth",
            query: { registered: "true" },
          });
          this.isLoginMode = true;
        }
      } catch (error) {
        this.handleAuthError(error, this.isLoginMode);
      }
    },

    handleAuthError(error, isLoginMode) {
      console.error("Error sending authentication request:", error);
      let serverError = "Authentication error";

      // Trying to parse error in JSON format
      try {
        const errData = JSON.parse(error.message);
        serverError = errData.error || serverError;
      } catch (e) {
        serverError = error.message || serverError;
      }

      // Mapping server messages to user-friendly signatures
      const errorMap = {
        "user already exists": "Пользователь уже существует.",
        "invalid email or password": "Неверный логин или пароль.",
        // Add additional mappings if needed
      };

      let displayError = serverError;
      const lowerError = serverError.toLowerCase();
      Object.keys(errorMap).forEach((key) => {
        if (lowerError.includes(key)) {
          displayError = errorMap[key];
        }
      });

      // Distribute the error across fields or set a global message
      if (isLoginMode) {
        if (lowerError.includes("email")) {
          this.errors.email = displayError;
        } else if (lowerError.includes("password")) {
          this.errors.password = displayError;
        } else {
          this.errorMessage = displayError;
        }
      } else {
        if (lowerError.includes("user already exists")) {
          this.errors.username = displayError;
        } else if (lowerError.includes("email")) {
          this.errors.email = displayError;
        } else if (lowerError.includes("password")) {
          this.errors.password = displayError;
        } else {
          this.errorMessage = displayError;
        }
      }
      this.showMobileError();
    },
    handleSubmit() {
      this.resetErrors();
      if (!this.validateForm()) {
        this.showMobileError();
        return;
      }
      this.sendAuthRequest();
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
      this.username = "";
      this.email = "";
      this.password = "";
      this.confirmPassword = "";
      this.isLoginMode = !this.isLoginMode;
    },
    togglePassword() {
      this.showPassword = !this.showPassword;
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
