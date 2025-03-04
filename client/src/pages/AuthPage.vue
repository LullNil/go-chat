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
      isLoginMode: true,
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      showPassword: false,
      errors: {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      },
      sunIcon: `M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8a4 4 0 110 8 4 4 0 010-8z`,
      moonIcon: `M21 12.79A9 9 0 1111.21 3a7 7 0 009.79 9.79z`,
      animateIcon: false,
      errorMessage: "",
      errorMessages: [],
      messageId: 0,
    };
  },
  watch: {
    // Animate the icon when the theme changes
    isDarkTheme() {
      this.animateIcon = true;
      setTimeout(() => {
        this.animateIcon = false;
      }, 600);
    },
  },
  methods: {
    // Handle form submission and validation
    handleSubmit() {
      this.resetErrors();
      let hasErrors = false;

      // Валидация
      if (!this.email.trim()) {
        this.errors.email = "Пожалуйста, введите почту.";
        hasErrors = true;
      } else if (!/.+@.+\..+/.test(this.email)) {
        this.errors.email = "Неверный формат почты.";
        hasErrors = true;
      }

      if (!this.password.trim()) {
        this.errors.password = "Пожалуйста, введите пароль.";
        hasErrors = true;
      }

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

      // Показываем сообщение об ошибке для мобильных
      if (hasErrors) {
        this.showMobileError();
      }

      // If there are no errors, emit the corresponding event with authentication data
      if (this.noErrors()) {
        const authData = this.isLoginMode
          ? { email: this.email, password: this.password }
          : {
              username: this.username,
              email: this.email,
              password: this.password,
            };

        this.$emit(this.isLoginMode ? "login" : "register", authData);
      }
    },

    // Reset error messages
    resetErrors() {
      this.errors = {
        username: "",
        email: "",
        password: "",
        confirmPassword: "",
      };
    },

    // Check if there are no errors
    noErrors() {
      return Object.values(this.errors).every((error) => error === "");
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

    removeMessage(id) {
      this.errorMessages = this.errorMessages.filter((msg) => msg.id !== id);
    },

    // Toggle between login and registration mode
    toggleMode() {
      this.resetErrors();
      this.username = "";
      this.email = "";
      this.password = "";
      this.confirmPassword = "";
      this.isLoginMode = !this.isLoginMode;
    },

    // Toggle password visibility
    togglePassword() {
      this.showPassword = !this.showPassword;
    },

    // Emit event to toggle theme
    toggleTheme() {
      this.$emit("toggle-theme");
    },

    // Clear error for a specific field
    clearError(field) {
      this.errors[field] = "";
    },
  },
};
</script>

<style scoped>
@import "@/styles/variables.css";
@import "@/styles/chat.css";
@import "@/styles/authpage.css";
</style>
