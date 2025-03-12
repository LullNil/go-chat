<template>
  <div class="auth-container">
    <!-- Error message for mobile devices -->
    <transition name="bounce-pop">
      <div v-if="errorMessage" class="mobile-error-message">
        <span class="error-icon">⚠</span>
        {{ errorMessage }}
      </div>
    </transition>

    <!-- Logo in the bottom-left corner -->
    <div class="logo-corner">© &laquo;GoChat&raquo; 2025</div>

    <!-- Theme switch button in the bottom-right corner -->
    <button
      class="theme-switch-corner"
      @click="toggleTheme"
      aria-label="Switch theme"
    >
      <svg
        class="theme-icon"
        :class="{ 'animate-icon': animateIcon }"
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

    <!-- Authentication card with header and form -->
    <div class="auth-card" :class="{ 'animate-container': !isLoginMode }">
      <div class="auth-header">
        <transition name="fade-text" mode="out-in">
          <h2 :key="isLoginMode">
            {{ isLoginMode ? "Вход" : "Регистрация" }}
          </h2>
        </transition>
      </div>
      <form @submit.prevent="handleSubmit" novalidate>
        <!-- Username input: appears above email in registration mode -->
        <div
          class="extra-field username-field"
          :class="{ active: !isLoginMode }"
        >
          <div
            class="auth-input-group"
            :style="{ marginBottom: errors.username ? '2.5rem' : '1.5rem' }"
          >
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
            <transition name="fade-error">
              <div v-if="errors.username" class="error-message">
                {{ errors.username }}
              </div>
            </transition>
          </div>
        </div>

        <!-- Email input group (always displayed) -->
        <div
          class="auth-input-group"
          :style="{ marginBottom: errors.email ? '2.5rem' : '1.5rem' }"
        >
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
          <transition name="fade-error">
            <div v-if="errors.email" class="error-message">
              {{ errors.email }}
            </div>
          </transition>
        </div>

        <!-- Password input group (always displayed) -->
        <div
          class="auth-input-group"
          :style="{ marginBottom: errors.password ? '2.5rem' : '1.5rem' }"
        >
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
          <transition name="fade-error">
            <div v-if="errors.password" class="error-message">
              {{ errors.password }}
            </div>
          </transition>
        </div>

        <!-- Confirm Password input: appears below password in registration mode -->
        <div
          class="extra-field confirm-field"
          :class="{ active: !isLoginMode }"
        >
          <div
            class="auth-input-group"
            :style="{
              marginBottom: errors.confirmPassword ? '2.5rem' : '1.5rem',
            }"
          >
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
            <transition name="fade-error">
              <div v-if="errors.confirmPassword" class="error-message">
                {{ errors.confirmPassword }}
              </div>
            </transition>
          </div>
        </div>

        <!-- Submit button -->
        <button type="submit" class="auth-button">
          {{ isLoginMode ? "Войти" : "Зарегистрироваться" }}
        </button>
      </form>
    </div>

    <!-- Footer: mode toggle and "Forgot password?" link -->
    <div class="auth-footer">
      <button class="mode-toggle" @click="toggleMode">
        {{ isLoginMode ? "Нет аккаунта? Создать" : "Уже есть аккаунт? Войти" }}
      </button>
      <a href="#" class="forgot-password">Забыли пароль?</a>
    </div>
  </div>
</template>

<script>
export default {
  name: "AuthPage",
  props: {
    isDarkTheme: { type: Boolean, default: false },
  },
  data() {
    return {
      // Mode and form data
      isLoginMode: true,
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      // UI state flags
      showPassword: false,
      animateIcon: false,
      errorMessage: "",
      // Error messages for individual fields
      errors: {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      },
      // Icons for theme switch
      sunIcon:
        "M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8a4 4 0 110 8 4 4 0 010-8z",
      moonIcon: "M21 12.79A9 9 0 1111.21 3a7 7 0 009.79 9.79z",
      // Unused error message array and id (reserved for potential future use)
      errorMessages: [],
      messageId: 0,
    };
  },
  watch: {
    // Animate the theme icon when the theme prop changes
    isDarkTheme() {
      this.animateIcon = true;
      setTimeout(() => {
        this.animateIcon = false;
      }, 600);
    },
  },
  methods: {
    /**
     * Resets all field error messages.
     */
    resetErrors() {
      this.errors = {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      };
    },

    /**
     * Clears the error message for a specific field.
     * @param {string} field - The field name whose error should be cleared.
     */
    clearError(field) {
      this.errors[field] = "";
    },

    /**
     * Checks if there are no validation errors.
     * @returns {boolean} True if there are no errors, false otherwise.
     */
    noErrors() {
      return Object.values(this.errors).every((error) => error === "");
    },

    /**
     * Displays the first error message in a mobile-friendly format.
     * The error message will automatically disappear after a short delay.
     */
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
     * Removes an error message from the errorMessages array by its ID.
     * (Reserved for potential future use)
     * @param {number} id - The ID of the message to remove.
     */
    removeMessage(id) {
      this.errorMessages = this.errorMessages.filter((msg) => msg.id !== id);
    },

    /**
     * Form validation.
     */
    validateForm() {
      let hasErrors = false;

      // Validate email field
      if (!this.email.trim()) {
        this.errors.email = "Пожалуйста, введите почту.";
        hasErrors = true;
      } else if (!/.+@.+\..+/.test(this.email)) {
        this.errors.email = "Неверный формат почты.";
        hasErrors = true;
      }

      // Validate password field
      if (!this.password.trim()) {
        this.errors.password = "Пожалуйста, введите пароль.";
        hasErrors = true;
      }

      // Additional validations for registration mode
      if (!this.isLoginMode) {
        if (!this.username.trim()) {
          this.errors.username = "Пожалуйста, введите имя пользователя.";
          hasErrors = true;
        }
        if (this.password !== this.confirmPassword) {
          this.errors.confirmPassword = "Пароли не совпадают.";
          hasErrors = true;
        }
      }

      return !hasErrors;
    },

    sendAuthRequest() {
      const url = this.isLoginMode
        ? "http://localhost:8081/login"
        : "http://localhost:8081/register";
      const authData = this.isLoginMode
        ? { email: this.email, password: this.password }
        : {
            username: this.username,
            email: this.email,
            password: this.password,
          };

      // TODO: delete when backend is ready
      this.$emit(this.isLoginMode ? "login" : "register", authData);

      fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(authData),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.$emit(this.isLoginMode ? "login" : "register", data);
        })
        .catch((error) => {
          console.error("Error sending authentication request:", error);
        })
        .catch((error) => {
          console.error("Error sending authentication request:", error);
          this.errorMessage = error.message;
          this.showMobileError();
        });
    },

    /**
     * Handles the form submission by validating inputs and emitting either
     * a 'login' or 'register' event with the authentication data if validation passes.
     */
    handleSubmit() {
      this.resetErrors();

      // Display mobile error if any validation errors exist
      if (!this.validateForm()) {
        this.showMobileError();
        return;
      }

      // sdfsdfg@c.com
      this.sendAuthRequest();
    },

    /**
     * Toggles between login and registration modes.
     * Resets form fields and error messages.
     */
    toggleMode() {
      this.resetErrors();
      this.username = "";
      this.email = "";
      this.password = "";
      this.confirmPassword = "";
      this.isLoginMode = !this.isLoginMode;
    },

    /**
     * Toggles the visibility of the password field.
     */
    togglePassword() {
      this.showPassword = !this.showPassword;
    },

    /**
     * Emits an event to toggle the application theme.
     */
    toggleTheme() {
      this.$emit("toggle-theme");
    },
  },
};
</script>

<style scoped>
@import "@/styles/variables.css";
@import "@/styles/chat.css";
@import "@/styles/authpage.css";
</style>
