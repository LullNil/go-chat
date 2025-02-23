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
    <div class="logo-corner">&laquo;GoChat©&raquo; 2025</div>

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
  name: "AuthForm",
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

/* Fade error transition */
.fade-error-enter-active,
.fade-error-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-error-enter-from,
.fade-error-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
.fade-error-enter-to,
.fade-error-leave-from {
  opacity: 1;
  transform: translateY(0);
}

/* Header styles */
.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}
.auth-header h2 {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-light);
  margin: 0;
}
.dark-theme .auth-header h2 {
  color: #fff;
}
.dark-theme .mode-toggle {
  color: var(--gradient-sent-end);
}

/* Footer styles */
.auth-footer {
  margin-top: 1rem;
  text-align: center;
}
.mode-toggle {
  background: none;
  border: none;
  color: var(--gradient-sent-end);
  cursor: pointer;
  font-weight: 500;
  transition: opacity 0.3s ease;
}
.mode-toggle:hover {
  opacity: 0.8;
}
.forgot-password {
  display: block;
  margin-top: 0.5rem;
  color: var(--text-light);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.5s ease;
}
.dark-theme .forgot-password {
  color: var(--text-light);
}

/* Smooth appearance for extra fields with a "jelly" effect */
.extra-field {
  z-index: 0;
  position: relative;
  max-height: 0;
  opacity: 0;
  transform: scaleY(0.8);
  transition: max-height 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275),
    opacity 0.5s ease, transform 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.extra-field.active {
  max-height: 200px;
  opacity: 1;
  transform: scaleY(1);
}
.extra-field .auth-input-group {
  position: relative;
}

/* Header text animation for login/registration */
.fade-text-enter-active {
  animation: textBounceIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
@keyframes textBounceIn {
  0% {
    opacity: 0;
    transform: scale(0.8);
  }
  50% {
    opacity: 1;
    transform: scale(1.1);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}
.fade-text-leave-active {
  animation: textFadeOut 0.3s ease;
}
@keyframes textFadeOut {
  0% {
    opacity: 1;
    transform: scale(1);
  }
  100% {
    opacity: 0;
    transform: scale(0.9);
  }
}

/* Ensure label stays above the input */
.auth-input-label {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  transition: all 0.3s ease;
  color: var(--text-light);
  background: none;
  padding: 0 0.5rem;
  z-index: 2;
}
.dark-theme .auth-input-label {
  color: var(--text-dark);
}
.auth-input:focus ~ .auth-input-label,
.auth-input:not(:placeholder-shown) ~ .auth-input-label {
  top: 0;
  font-size: 0.875rem;
  color: var(--gradient-sent-end);
}
.error-label {
  color: #ff4f4f;
}
.dark-theme .error-label {
  color: #ff4f4f;
}

/* Error messages and margin transitions */
.error-message {
  position: absolute;
  left: 1rem;
  bottom: -1.2rem;
  font-size: 0.875rem;
  color: #ff4f4f;
  transition: opacity 0.3s ease, transform 0.3s ease;
}

/* Password visibility toggle button */
.password-toggle {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
}
.eye-icon {
  width: 24px;
  height: 24px;
  transition: stroke 0.3s ease, transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.dark-theme .eye-icon {
  stroke: #fff;
}
.eye-cross {
  stroke: currentColor;
  stroke-width: 1.5;
  stroke-linecap: butt;
  stroke-dasharray: 22.63;
  stroke-dashoffset: 22.63;
  opacity: 0;
  transition: stroke-dashoffset 0.15s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.15s ease;
}
.dark-theme .eye-cross {
  stroke: #fff;
}
.eye-cross.drawn {
  stroke-dashoffset: 0;
  opacity: 1;
}
.eye-icon.animate-closed {
  animation: eyeElastic 0.7s ease-in-out forwards;
}
@keyframes eyeElastic {
  0% {
    transform: scale(1) rotate(0deg);
    filter: brightness(1);
  }
  20% {
    transform: scale(1.1) rotate(5deg) translateY(-2px);
    filter: brightness(1.2);
  }
  40% {
    transform: scale(0.95) rotate(-3deg) translateY(2px);
    filter: brightness(0.9);
  }
  60% {
    transform: scale(1.03) rotate(2deg) translateY(-1px);
    filter: brightness(1.05);
  }
  80% {
    transform: scale(0.98) rotate(-1deg) translateY(1px);
    filter: brightness(0.98);
  }
  100% {
    transform: scale(1) rotate(0deg) translateY(0);
    filter: brightness(1);
  }
}

/* Authorization container styles */
.auth-container {
  user-select: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: var(--primary-light);
  transition: background 0.5s ease;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  overflow-y: auto;
}

/* Logo and theme switch button positioning */
.logo-corner {
  position: absolute;
  bottom: 16px;
  left: 10px;
  font-size: 1.1rem;
  font-weight: 100;
  color: rgba(0, 0, 0, 0.5);
  z-index: 10;
}
.dark-theme .logo-corner {
  color: rgba(255, 255, 255, 0.3);
}
.theme-switch-corner {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  z-index: 10;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.theme-switch-corner .theme-icon {
  width: 32px;
  height: 32px;
  transition: transform 0.6s ease, stroke 0.3s ease;
  stroke: rgba(0, 0, 0, 0.5);
}
.dark-theme .theme-switch-corner .theme-icon {
  stroke: rgba(255, 255, 255, 0.3);
}
.animate-icon {
  animation: icon-pop 0.6s ease;
}
@keyframes icon-pop {
  0% {
    transform: scale(1) rotate(0deg);
  }
  50% {
    transform: scale(1.2) rotate(45deg);
  }
  100% {
    transform: scale(1) rotate(0deg);
  }
}

/* Authentication card with container animation */
.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2rem 2.5rem;
  background: var(--glass-effect-light, rgba(255, 255, 255, 0.9));
  backdrop-filter: blur(15px);
  border-radius: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.animate-container {
  animation: containerJelly 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
@keyframes containerJelly {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.03);
  }
  100% {
    transform: scale(1);
  }
}
.dark-theme .auth-container {
  background: var(--primary-dark);
}
.dark-theme .auth-card {
  background: var(--secondary-dark);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 0.05);
}

/* Input group and input styles */
.auth-input-group {
  position: relative;
  transition: margin-bottom 0.5s ease;
  margin-bottom: 1.5rem !important;
}
.auth-input {
  width: 100%;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 0.75rem;
  font-size: 1rem;
  transition: all 0.3s ease;
  outline: none;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}
.auth-input:focus {
  outline: none;
  box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}
.dark-theme .auth-input {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-dark);
}

/* Submit button styles */
.auth-button {
  position: relative;
  z-index: 1;
  width: 100%;
  padding: 1rem;
  background: var(--gradient-sent-end);
  color: white;
  border: none;
  border-radius: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 1rem;
}
.auth-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(88, 101, 242, 0.25);
}

.mobile-error-message {
  display: none;
}

/* Mobile adaptation */
@media (max-width: 768px) {
  .error-message {
    display: none;
  }

  .dark-theme .error-label {
    color: #dc2626;
  }

  /* Креативная анимация с эффектом пружины */
  .bounce-pop-enter-active {
    animation: bounce-in 0.5s cubic-bezier(0.18, 0.89, 0.32, 1.28);
  }

  .bounce-pop-leave-active {
    animation: fade-slide-up 0.3s cubic-bezier(0.55, 0.06, 0.68, 0.19);
  }

  @keyframes bounce-in {
    0% {
      transform: translateY(8px) scale(0.98);
      opacity: 0;
    }
    60% {
      transform: translateY(-4px) scale(1.05);
    }
    80% {
      transform: translateY(2px) scale(0.995);
    }
    100% {
      transform: translateY(0) scale(1);
      opacity: 1;
    }
  }

  @keyframes fade-slide-up {
    0% {
      transform: translateY(0) scale(1);
      opacity: 1;
    }
    100% {
      transform: translateY(-40px) scale(0.85);
      opacity: 0;
    }
  }

  .mobile-error-message {
    position: fixed;
    top: 8vh;
    /* left: 5%;
    right: 5%; */
    width: 75%;
    z-index: 1000;
    background: linear-gradient(145deg, #ff5f5f, #ff4040);
    color: white;
    padding: 15px 25px;
    border-radius: 15px;
    font-size: 0.95rem;
    box-shadow: 0 8px 25px rgba(255, 79, 79, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    backdrop-filter: blur(8px);
    border: 2px solid rgba(255, 255, 255, 0.15);
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }

  .error-icon {
    animation: icon-shake 1.5s ease infinite;
    font-size: 1.4rem;
    filter: drop-shadow(0 2px 2px rgba(0, 0, 0, 0.2));
  }

  @keyframes icon-shake {
    0%,
    60%,
    100% {
      transform: rotate(0deg) scale(1);
      animation-timing-function: cubic-bezier(0.4, 0, 0.6, 1);
    }
    10%,
    30%,
    50% {
      transform: rotate(12deg) scale(1.08);
    }
    20%,
    40% {
      transform: rotate(-12deg) scale(1.08);
    }
  }

  .element-with-icon {
    animation: icon-shake 1.4s infinite;
  }

  /* Адаптация для темной темы */
  .dark-theme .mobile-error-message {
    background: linear-gradient(145deg, #dc2626, #b91c1c);
    border-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 8px 25px rgba(220, 38, 38, 0.2);
  }

  .auth-container {
    align-items: center;
    padding: 1rem;
    height: 100vh;
  }
  .auth-card {
    width: 100%;
    max-width: 100%;
    padding: 1.5rem;
    margin-top: 1rem;
    background: transparent;
    box-shadow: none;
    border-radius: 0;
    transform: none;
  }
  .auth-header h2 {
    font-size: 1.5rem;
  }
  .auth-input {
    font-size: 0.9rem;
    padding: 0.8rem;
  }
  .auth-input-group {
    margin-bottom: 0.75rem !important; /* Уменьшаем отступы между полями */
  }

  .auth-button {
    padding: 0.8rem;
    font-size: 0.9rem;
    height: 45px;
  }
  .logo-corner,
  .theme-switch-corner {
    position: fixed;
  }
  .logo-corner {
    bottom: 10px;
    left: 10px;
    font-size: 0.9rem;
  }
  .theme-switch-corner {
    bottom: 5px;
    right: 10px;
    width: 40px;
    height: 40px;
  }
  .theme-switch-corner .theme-icon {
    width: 28px;
    height: 28px;
  }
  .dark-theme .auth-card {
    background: transparent;
    box-shadow: none;
    border-radius: 0;
  }
}
</style>
